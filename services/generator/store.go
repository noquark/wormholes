package main

const (
	BUCKET_EMPTY = iota
	BUCKET_BUSY  = iota
	BUCKET_FULL  = iota
)

// An in-memory store with multiple buckets and their status
type MemStore struct {
	status   map[int]int
	buckets  [][]string
	capacity int
}

// Create a new memory store for given bucket size and capacity
func NewMemStore(size, capacity int) *MemStore {
	memStore := &MemStore{
		status:   make(map[int]int, size),
		buckets:  make([][]string, size),
		capacity: capacity,
	}
	for i := range memStore.buckets {
		memStore.buckets[i] = make([]string, capacity)
		memStore.status[i] = BUCKET_EMPTY
	}
	return memStore
}

// Get index of all buckets that are empty
func (s *MemStore) GetEmpty() []int {
	var emptyBucketIDs []int
	for id, status := range s.status {
		if status == BUCKET_EMPTY {
			emptyBucketIDs = append(emptyBucketIDs, id)
		}
	}
	return emptyBucketIDs
}

// Pop first bucket that is full
func (s *MemStore) Pop() []string {
	var data []string
	for id, status := range s.status {
		if status == BUCKET_FULL {
			data = make([]string, s.capacity)
			copy(data, s.buckets[id])
			s.status[id] = BUCKET_EMPTY
			s.buckets[id] = make([]string, s.capacity)
			break
		}
	}
	return data
}
