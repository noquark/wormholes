package main

import (
	"io/fs"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
)

const (
	DirPerm  fs.FileMode = 0x755
	FilePerm fs.FileMode = 0x600
)

// A thread safe wrapper around bloom filter with backup and restore
type Bloom struct {
	name  string
	bloom *bloom.BloomFilter
	mutex sync.RWMutex
}

func NewBloom(name string, maxLimit uint, errorRate float64) *Bloom {
	return &Bloom{
		name:  name,
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
