package generator

import (
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

const (
	BucketEmpty = iota
	BucketBusy  = iota
	BucketFull  = iota
)

// An in-memory store with multiple buckets and their status.
type MemStore struct {
	mutex    sync.RWMutex
	status   map[int]int
	buckets  [][]string
	capacity int
}

// Create a new memory store for given bucket size and capacity.
func NewMemStore(size, capacity int) *MemStore {
	memStore := &MemStore{
		mutex:    sync.RWMutex{},
		status:   make(map[int]int, size),
		buckets:  make([][]string, size),
		capacity: capacity,
	}
	for i := range memStore.buckets {
		memStore.buckets[i] = make([]string, capacity)
		memStore.status[i] = BucketEmpty
	}

	log.Info().Msgf("memstore: bucket capacity %s", humanize.Comma(int64(capacity)))
	log.Info().Msgf("memstore: number of buckets %d", size)

	return memStore
}

// Get index of all buckets that are empty.
func (s *MemStore) GetEmpty() []int {
	var emptyBucketIDs []int

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, status := range s.status {
		if status == BucketEmpty {
			emptyBucketIDs = append(emptyBucketIDs, id)
		}
	}

	return emptyBucketIDs
}

// Pop first bucket that is full.
func (s *MemStore) Pop() []string {
	var data []string

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, status := range s.status {
		if status == BucketFull {
			data = make([]string, s.capacity)
			copy(data, s.buckets[id])
			s.status[id] = BucketEmpty
			s.buckets[id] = make([]string, s.capacity)

			break
		}
	}

	return data
}
