package generator

import (
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

type Bucket struct {
	sync.RWMutex
	capacity int
	data     []string
}

func (b *Bucket) Pop() []string {
	data := make([]string, b.capacity)
	copy(data, b.data)
	b.data = nil
	return data
}

// An in memory store with buckets
type MemStore struct {
	buckets []*Bucket
	empty   chan int
}

// Create a new memory store for given bucket size and capacity.
func NewMemStore(size, capacity int) *MemStore {
	memStore := &MemStore{
		buckets: make([]*Bucket, size),
		empty:   make(chan int, size),
	}
	for i := range memStore.buckets {
		memStore.buckets[i] = &Bucket{capacity: capacity}
	}

	log.Info().Msgf("bucket capacity %s", humanize.Comma(int64(capacity)))
	log.Info().Msgf("number of buckets %d", size)

	return memStore
}

// Pop first bucket that is full.
func (s *MemStore) Pop() []string {
	for id, bucket := range s.buckets {
		if isAvailable := bucket.TryLock(); isAvailable {
			if bucket.data != nil {
				data := bucket.Pop()
				bucket.Unlock()
				s.empty <- id
				return data
			} else {
				bucket.Unlock()
				continue
			}
		} else {
			log.Warn().Msgf("bucket %d is busy", id)
		}
	}
	return nil
}
