package main

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"
	"wormholes/protos"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mohitsinghs/nanoid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	queryIDs string = `SELECT id from wh_links`
)

var ErrFactoryEmpty = status.New(codes.ResourceExhausted, "factory is empty").Err()

// An ID generation factory made of -
//  - A db connection for loading IDs on startup
//  - A bloom-filter based on all existing ids to avoid collisions
//  - size of ID (default is 7)
type Factory struct {
	protos.UnimplementedBucketServiceServer
	db    *pgxpool.Pool
	bloom *Bloom
	store *MemStore
	size  int
	tick  *time.Ticker
}

func NewFactory(config *Config) *Factory {
	return &Factory{
		db:    config.Postgres.Connect(),
		bloom: NewBloom(config.MaxLimit, config.ErrorRate),
		store: NewMemStore(config.Size, config.Capacity),
		size:  config.IDSize,
		tick:  time.NewTicker(time.Second),
	}
}

func (f *Factory) Prepare() *Factory {
	rows, err := f.db.Query(context.Background(), queryIDs)
	if err != nil {
		log.Warn().Err(err).Msg("failed to retrieve IDs")
		return f
	}
	defer rows.Close()
	var idCount uint64
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse ID")

			continue
		}

		atomic.AddUint64(&idCount, 1)
		f.bloom.Add(fasterByte(id))

	}
	if idCount > 0 {
		log.Info().Msgf("cached %d IDs for lookup", idCount)
	}

	return f
}

func (f *Factory) Run() *Factory {
	go func() {
		for {
			select {
			case <-f.tick.C:
				if emptyBuckets := f.store.GetEmpty(); len(emptyBuckets) > 0 {
					for _, idx := range emptyBuckets {
						f.store.mutex.Lock()
						go f.populateBucket(idx)
						f.store.status[idx] = BUCKET_BUSY
						f.store.mutex.Unlock()
					}
				}
			}
		}
	}()
	return f
}

func (f *Factory) Shutdown() {
	f.tick.Stop()
}

// populate bucket at given index until full
func (f *Factory) populateBucket(idx int) {
	t := time.Now()
	log.Info().Msgf("filling bucket %d", idx)

	fillCount := 0
	for fillCount < f.store.capacity {
		id, err := nanoid.New(f.size)
		if err == nil && id != "" {
			if !f.bloom.Exists(fasterByte(id)) {
				f.store.buckets[idx][fillCount] = id
				f.bloom.Add(fasterByte(id))
				fillCount++
			}
		}
	}
	f.store.mutex.Lock()
	f.store.status[idx] = BUCKET_FULL
	f.store.mutex.Unlock()
	log.Info().Msgf("filled bucket %d in %s", idx, time.Since(t).String())
}

func (f *Factory) GetBucket(context context.Context, empty *protos.Empty) (*protos.Bucket, error) {
	ids := f.store.Pop()
	if len(ids) == 0 {
		return nil, ErrFactoryEmpty
	}
	return &protos.Bucket{
		Ids: ids,
	}, nil
}

func fasterByte(s string) []byte {
	return unsafe.Slice(
		(*byte)(unsafe.Pointer(
			(*reflect.StringHeader)(unsafe.Pointer(&s)).Data)),
		len(s),
	)
}
