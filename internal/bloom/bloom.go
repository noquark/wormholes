package bloom

import (
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

const (
	ByteSize = 8
)

// A thread safe wrapper around bloom filter with backup and restore.
type Bloom struct {
	bloom *bloom.BloomFilter
	mutex sync.RWMutex
}

func New(maxLimit uint, errorRate float64) *Bloom {
	b := &Bloom{
		bloom: bloom.NewWithEstimates(maxLimit, errorRate),
		mutex: sync.RWMutex{},
	}

	log.Info().Msgf("bloom-filter: size %s", humanize.Bytes(uint64(b.bloom.Cap()/ByteSize)))
	log.Info().Msgf("bloom-filter: limit %s", humanize.Comma(int64(maxLimit)))
	log.Info().Msgf("bloom-filter: errorRate %f", errorRate)

	return b
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
