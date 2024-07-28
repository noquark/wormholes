package memstore

import (
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

type Bucket struct {
	sync.RWMutex
	Capacity int
	Data     []string
}

func (b *Bucket) Pop() []string {
	data := make([]string, b.Capacity)
	copy(data, b.Data)
	b.Data = nil
	return data
}

// An in memory store with buckets
type MemStore struct {
	Buckets []*Bucket
	Empty   chan int
}

// Create a new memory store for given bucket size and capacity.
func New(size, capacity int) *MemStore {
	memStore := &MemStore{
		Buckets: make([]*Bucket, size),
		Empty:   make(chan int, size),
	}
	for i := range memStore.Buckets {
		memStore.Buckets[i] = &Bucket{Capacity: capacity}
	}

	log.Info().Msgf("bucket capacity %s", humanize.Comma(int64(capacity)))
	log.Info().Msgf("number of buckets %d", size)

	return memStore
}

// Pop first bucket that is full.
func (s *MemStore) Pop() []string {
	for id, bucket := range s.Buckets {
		if isAvailable := bucket.TryLock(); isAvailable {
			if bucket.Data != nil {
				data := bucket.Pop()
				bucket.Unlock()
				s.Empty <- id
				log.Info().Msgf("popped bucket %d", id)
				return data
			} else {
				bucket.Unlock()
				continue
			}
		}
	}
	return nil
}
