package main

import (
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
)

// A thread safe wrapper around bloom filter with backup and restore
type Bloom struct {
	bloom *bloom.BloomFilter
	mutex sync.RWMutex
}

func NewBloom(maxLimit uint, errorRate float64) *Bloom {
	return &Bloom{
		bloom: bloom.NewWithEstimates(maxLimit, errorRate),
		mutex: sync.RWMutex{},
	}
}

func (b *Bloom) Add(id []byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.bloom.Add(id)
}

func (b *Bloom) Exists(id []byte) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.bloom.Test(id)
}
