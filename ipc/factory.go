package ipc

import (
	"context"
	"time"
	"unsafe"
	"wormholes/internal/bloom"
	"wormholes/internal/config"
	"wormholes/internal/memstore"
	"wormholes/protos"

	"github.com/dustin/go-humanize"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noquark/nanoid"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	queryIDs      string = `SELECT id from links`
	queryIDsCount string = `SELECT count(id) from links`
	maxBarWidth          = 64
)

type Factory struct {
	protos.UnimplementedBucketServiceServer
	db     *pgxpool.Pool
	bloom  *bloom.Bloom
	store  *memstore.MemStore
	config *config.Config
}

func NewFactory(config *config.Config, db *pgxpool.Pool) *Factory {
	return &Factory{
		db:     db,
		bloom:  bloom.New(config.BloomMaxLimit, config.BloomErrorRate),
		store:  memstore.New(config.BucketSize, config.BucketCapacity),
		config: config,
	}
}

func (f *Factory) Prepare() *Factory {
	var idCount uint64

	err := f.db.QueryRow(context.Background(), queryIDsCount).Scan(&idCount)
	if err != nil {
		log.Warn().Err(err).Msg("factory: failed to get IDs count")
	}

	if idCount > 0 {
		rows, err := f.db.Query(context.Background(), queryIDs)
		if err != nil {
			log.Warn().Err(err).Msg("factory: failed to get IDs")

			return f
		}
		defer rows.Close()

		bar := progressbar.NewOptions(
			int(idCount),
			progressbar.OptionSetWidth(maxBarWidth),
			progressbar.OptionClearOnFinish(),
		)
		for rows.Next() {
			var id string

			err := rows.Scan(&id)
			if err != nil {
				log.Warn().Err(err).Msg("factory: failed to parse ID")

				continue
			}

			bar.Add(1)
			f.bloom.Add(fasterByte(id))
		}
		bar.Finish()
		log.Info().Msgf("factory: cached %s IDs", humanize.Comma(int64(idCount)))
	}

	return f
}

func (f *Factory) Run(conf *config.Config) *Factory {
	for i := range f.store.Buckets {
		go f.populateBucket(i)
	}
	go func() {
		for idx := range f.store.Empty {
			go f.populateBucket(idx)
		}
	}()

	return f
}

// populate bucket at given index until full.
func (f *Factory) populateBucket(idx int) {
	t := time.Now()
	fillCount := 0
	bucket := f.store.Buckets[idx]
	idSize := f.config.IDSize
	if isAvailable := bucket.TryLock(); isAvailable {
		log.Info().Msgf("filling bucket %d", idx)
		bucket.Data = make([]string, bucket.Capacity)
		for fillCount < bucket.Capacity {
			id, err := nanoid.New(idSize)
			if err == nil && id != "" {
				if !f.bloom.Exists(fasterByte(id)) {
					bucket.Data[fillCount] = id
					f.bloom.Add(fasterByte(id))
					fillCount++
				}
			}
		}
		bucket.Unlock()
		log.Info().Msgf("filled bucket %d in %s", idx, time.Since(t).String())
	} else {
		log.Warn().Msgf("bucket not available %d", idx)
	}
}

func (f *Factory) Shutdown() {
	close(f.store.Empty)
}
func (f *Factory) GetBucket(context context.Context, empty *protos.Empty) (*protos.Bucket, error) {
	t := time.Now()
	ids := f.store.Pop()
	if ids != nil {
		log.Info().Msgf("get bucket in %s", time.Since(t).String())
		p := &protos.Bucket{
			Ids: ids,
		}
		log.Info().Msgf("serialized bucket in %s", time.Since(t).String())
		return p, nil
	} else {
		timer := time.NewTimer(f.config.Timeout)
		for range timer.C {
			if ids = f.store.Pop(); ids != nil {
				return &protos.Bucket{
					Ids: ids,
				}, nil
			}
		}
		log.Warn().Caller().Msgf("timed out, none of the buckets are filled")
		return nil, status.New(codes.ResourceExhausted, "factory: it's empty here").Err()
	}
}

func fasterByte(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
