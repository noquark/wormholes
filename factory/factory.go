package factory

import (
	"errors"
	"log"
	"os"
	"path"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
	homedir "github.com/mitchellh/go-homedir"
	boom "github.com/tylertreat/BoomFilters"
)

const (
	MAX_TRIES = 10
	ID_SIZE   = 7
)

type Factory struct {
	Generated *boom.ScalableBloomFilter
	mutex     sync.RWMutex
}

func New() *Factory {
	return &Factory{
		Generated: boom.NewDefaultScalableBloomFilter(0.001),
		mutex:     sync.RWMutex{},
	}
}

func (f *Factory) NewId() (string, error) {
	return f.newWithSize(ID_SIZE)
}

// Generate a new nanoid of given size and fallback to
// failSafeGenId on error or collsion in simple id generation
func (f *Factory) newWithSize(size int) (string, error) {
	id, err := gonanoid.New(size)
	if err != nil || f.Exists(id) {
		id = f.failSafeGenId(size)
	}
	if id == "" {
		return "", errors.New("unable to generate valid id in max iterations")
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
func CanRestore() (bool, error) {
	bp, err := bloomPath()
	if err != nil {
		return false, err
	}
	info, err := os.Stat(bp)
	if err != nil {
		if os.IsNotExist(err) || info.Size() == 0 {
			return false, nil
		}
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// Backup generated bloom filter to file in $HOME/.wh/gen.db
func (f *Factory) Backup() error {
	bp, err := bloomPath()
	if err != nil {
		return err
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
func (f *Factory) TryRestore(idsFunc func() ([]string, error)) {
	canRestore, err := CanRestore()
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
}

// Restore generated bloom filter from $HOME/.wh/gen.db
func (f *Factory) RestoreFromFile() error {
	bp, err := bloomPath()
	if err != nil {
		log.Panicln(err.Error())
	}
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
	for i := 0; i < MAX_TRIES; i++ {
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

// Get backup path for bloom filter creating directories
// if they don't exist
func bloomPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	wormDir := path.Join(home, ".wh")
	_, err = os.Stat(wormDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(wormDir, 0775); err != nil {
			return "", err
		}
	}
	return path.Join(wormDir, "gen.db"), nil
}
