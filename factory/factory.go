package factory

import (
	"errors"
	"log"
	"os"
	"path"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mohitsinghs/wormholes/config"
)

type Factory struct {
	Generated *bloom.BloomFilter
	mutex     sync.RWMutex
	conf      *config.FactoryConfig
}

func New(config *config.FactoryConfig) *Factory {
	return &Factory{
		Generated: bloom.NewWithEstimates(config.MaxLimit, config.ErrorRate),
		mutex:     sync.RWMutex{},
		conf:      config,
	}
}

func (f *Factory) NewId() (string, error) {
	id, err := gonanoid.New(f.conf.IdSize)
	if err != nil || f.Exists(id) {
		id = f.failSafeGenId(f.conf.IdSize)
	}
	if id == "" {
		return "", errors.New("unable to generate valid id")
	}
	f.Add(id)
	return id, nil
}

func (f *Factory) NewCookie() (string, error) {
	id, err := gonanoid.New(21)
	if err != nil || f.Exists(id) {
		id = f.failSafeGenId(21)
	}
	if id == "" {
		return "", errors.New("unable to generate valid cookie")
	}
	return id, nil
}

// Add unique id to genrated
func (f *Factory) Add(id string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.Generated.Add([]byte(id))
}

// Check if id exits in generated
func (f *Factory) Exists(id string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.Generated.Test([]byte(id))
}

// Check if backup of generated bloom filter can be restored from file
func (f *Factory) CanRestore() (bool, error) {
	bp := f.conf.BackupPath
	info, err := os.Stat(bp)
	if err != nil {
		if os.IsNotExist(err) || info.Size() == 0 {
			return false, nil
		}
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// Backup generated bloom filter to file
func (f *Factory) Backup() error {
	bp := f.conf.BackupPath
	bpDir := path.Dir(bp)
	_, err := os.Stat(bpDir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(bpDir, 0775)
	}
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	data, err := f.Generated.GobEncode()
	if err != nil {
		return err
	}
	os.WriteFile(bp, data, 0664)
	return nil
}

// Try restoring bloom filter either from file or from database
func (f *Factory) TryRestore(idsFunc func() ([]string, error)) *Factory {
	canRestore, err := f.CanRestore()
	if err != nil {
		log.Printf("Error getting backup : %v", err)
	}
	if canRestore {
		if err := f.RestoreFromFile(); err != nil {
			log.Printf("Error restoring from file : %v", err)
		}
	} else {
		ids, err := idsFunc()
		if err != nil {
			log.Printf("Error getting ids : %v", err)
		}
		f.RestoreFromIds(ids)
	}
	return f
}

// Restore generated bloom filter from backup
func (f *Factory) RestoreFromFile() error {
	bp := f.conf.BackupPath
	data, err := os.ReadFile(bp)
	if err != nil {
		return err
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	err = f.Generated.GobDecode(data)
	if err != nil {
		return err
	}
	return nil
}

// Restore generated bloom filter from all ids
// This is useful when bloom filter backup is missing or corrupt
func (f *Factory) RestoreFromIds(ids []string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	log.Println("Len : ", len(ids))
	for _, id := range ids {
		f.Generated.Add([]byte(id))
	}
}

// Try genrating new id at least 10 times before failing
// return empty string on failure
func (f *Factory) failSafeGenId(size int) string {
	id := ""
	for i := 0; i < f.conf.MaxTry; i++ {
		id, err := gonanoid.New(size)
		if err != nil || f.Exists(id) {
			continue
		}
		f.mutex.Lock()
		f.Generated.Add([]byte(id))
		f.mutex.Unlock()
		break
	}
	return id
}
