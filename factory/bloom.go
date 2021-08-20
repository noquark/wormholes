package factory

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/mohitsinghs/wormholes/constants"
	"github.com/rs/zerolog/log"
)

// A thread safe wrapper around bloom filters

type Bloom struct {
	name     string
	location string
	bloom    *bloom.BloomFilter
	mutex    sync.RWMutex
}

func NewBloom(name, location string, maxLimit uint, errorRate float64) *Bloom {
	return &Bloom{
		name:     name,
		location: location,
		bloom:    bloom.NewWithEstimates(maxLimit, errorRate),
		mutex:    sync.RWMutex{},
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

func (b *Bloom) Backup() {
	_, err := os.Stat(b.location)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(b.location, constants.DirPerm)
		} else {
			log.Panic().Err(err).Msg("can't access bloom path")
		}
	}

	info, err := os.Stat(b.location)
	if info.IsDir() {
		b.mutex.Lock()
		defer b.mutex.Unlock()

		data, err := b.bloom.GobEncode()
		if err != nil {
			log.Error().Err(err).Str("bloom", b.name).Msg("failed to backup")
		}

		err = os.WriteFile(path.Join(b.location, fmt.Sprintf("%s.bloom", b.name)), data, constants.FilePerm)
		if err == nil {
			log.Info().Str("bloom", b.name).Msg("backup created")
		} else {
			log.Error().Err(err).Str("bloom", b.name).Msg("failed to backup")
		}
	} else {
		log.Panic().Err(err).Msg("bloom path is not a directory")
	}
}

func (b *Bloom) Restore() bool {
	bp := path.Join(b.location, fmt.Sprintf("%s.bloom", b.name))

	info, err := os.Stat(bp)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info().Str("bloom", b.name).Msg("backup not found, skipping restore")
		} else {
			log.Error().Err(err).Str("bloom", b.name).Msg("failed to restore")
		}

		return false
	}

	if info.Size() == 0 {
		log.Info().Str("bloom", b.name).Msg("backup is empty, skipping restore")

		return false
	}

	if !info.Mode().IsRegular() {
		log.Info().Str("bloom", b.name).Msg("backup is not a file, skipping restore")

		return false
	}

	data, err := os.ReadFile(bp)
	if err != nil {
		log.Error().Err(err).Str("bloom", b.name).Msg("failed to restore")

		return false
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	err = b.bloom.GobDecode(data)
	if err != nil {
		log.Error().Err(err).Str("bloom", b.name).Msg("failed to restore")

		return false
	}

	log.Info().Str("bloom", b.name).Msg("restored")

	return true
}

func (b *Bloom) TryRestore(restoreFunc func() ([]string, error)) {
	if restoreFunc != nil {
		items, err := restoreFunc()
		if err == nil {
			b.mutex.Lock()
			defer b.mutex.Unlock()

			for _, item := range items {
				b.bloom.Add([]byte(item))
			}

			log.Info().Str("bloom", b.name).Msg("restored from database")
		} else {
			log.Error().Err(err).Str("bloom", b.name).Msg("failed to restore from database")
		}
	}
}
