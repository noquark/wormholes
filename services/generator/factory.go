package generator

import (
	"context"
	"errors"
	"reflect"
	"time"
	"unsafe"
	"wormholes/protos"

	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noquark/nanoid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	queryIDs      string = `SELECT id from links`
	queryIDsCount string = `SELECT count(id) from links`
	maxBarWidth          = 64
)

// An ID generation factory made of -
//   - A db connection for loading IDs on startup
//   - A bloom-filter based on all existing ids to avoid collisions
//   - An in-memory store to hold generated Ids
type Factory struct {
	protos.UnimplementedBucketServiceServer
	db     *pgxpool.Pool
	bloom  *Bloom
	store  *MemStore
	config *Config
}

func NewFactory(config *Config, db *pgxpool.Pool) *Factory {
	return &Factory{
		db:     db,
		bloom:  NewBloom(config.BloomMaxLimit, config.BloomErrorRate),
		store:  NewMemStore(config.BucketSize, config.BucketCapacity),
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

		bar := pb.Full.Start(int(idCount)).SetMaxWidth(maxBarWidth)

		for rows.Next() {
			var id string

			err := rows.Scan(&id)
			if err != nil {
				log.Warn().Err(err).Msg("factory: failed to parse ID")

				continue
			}

			bar.Increment()
			f.bloom.Add(fasterByte(id))
		}
		bar.Finish()
		log.Info().Msgf("factory: cached %s IDs", humanize.Comma(int64(idCount)))
	}

	return f
}

func (f *Factory) Run() *Factory {
	for i := range f.store.buckets {
		go f.populateBucket(i)
	}
	go func() {
		for idx := range f.store.empty {
			go f.populateBucket(idx)
		}
	}()

	return f
}

func (f *Factory) Shutdown() {
	close(f.store.empty)
}

// populate bucket at given index until full.
func (f *Factory) populateBucket(idx int) {
	t := time.Now()
	fillCount := 0
	bucket := f.store.buckets[idx]
	idSize := f.config.IDSize
	if isAvailable := bucket.TryLock(); isAvailable {
		log.Info().Msgf("filling bucket %d", idx)
		bucket.data = make([]string, bucket.capacity)
		for fillCount < bucket.capacity {
			id, err := nanoid.New(idSize)
			if err == nil && id != "" {
				if !f.bloom.Exists(fasterByte(id)) {
					bucket.data[fillCount] = id
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

func (f *Factory) GetLocalBucket() ([]string, error) {
	ids := f.store.Pop()
	if ids == nil {
		timer := time.NewTimer(f.config.Timeout)
		for range timer.C {
			if ids = f.store.Pop(); ids != nil {
				return ids, nil
			}
		}
		log.Warn().Caller().Msgf("timed out, none of the buckets are filled")
		return nil, errors.New("factory: it's empty here")
	}

	return ids, nil
}

func fasterByte(s string) []byte {
	return unsafe.Slice(
		(*byte)(unsafe.Pointer(
			(*reflect.StringHeader)(unsafe.Pointer(&s)).Data)),
		len(s),
	)
}
