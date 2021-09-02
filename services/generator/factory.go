package main

import (
	"context"
	"reflect"
	"time"
	"unsafe"
	"wormholes/protos"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mohitsinghs/nanoid"
	"github.com/rs/zerolog/log"
)

const (
	queryIDs string = `SELECT id from wh_links`
)

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
		bloom: NewBloom("id", uint(config.MaxLimit), config.ErrorRate),
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
	if len(rows.RawValues()) == 0 {
		log.Warn().Msg("there are no IDs yet")
	}
	for rows.Next() {
		log.Info().Msg("caching IDs")

		var id string
		err := rows.Scan(&id)
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse ID")

			continue
		}

		f.bloom.Add(unsafeFromString(id))

	}
	return f
}

func (f *Factory) Run() *Factory {
	go func() {
		for {
			select {
			case <-f.tick.C:
				if emptyBuckets := f.store.GetEmpty(); len(emptyBuckets) > 0 {
					for idx := range emptyBuckets {
						go f.populateBucket(idx)
						f.store.status[idx] = BUCKET_BUSY
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
	fillCount := 0
	for fillCount < f.store.capacity {
		id, err := nanoid.New(f.size)
		if err == nil && id != "" && !f.bloom.Exists(unsafeFromString(id)) {
			f.store.buckets[idx][fillCount] = id
			f.bloom.Add(unsafeFromString(id))
			fillCount++
		}
	}
	f.store.status[idx] = BUCKET_FULL
}

func (f *Factory) GetBucket(context context.Context, empty *protos.Empty) (*protos.Bucket, error) {
	return &protos.Bucket{
		Ids: f.store.Pop(),
	}, nil
}

func unsafeFromString(s string) []byte {
	return unsafe.Slice(
		(*byte)(unsafe.Pointer(
			(*reflect.StringHeader)(unsafe.Pointer(&s)).Data)),
		len(s),
	)
}
