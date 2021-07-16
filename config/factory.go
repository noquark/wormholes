package config

import (
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/mohitsinghs/wormholes/constants"
)

type FactoryConfig struct {
	MaxLimit   uint
	ErrorRate  float64
	MaxTry     int
	IdSize     int
	BackupPath string
}

func DefaultFactory() FactoryConfig {
	bp, err := bloomPath()
	if err != nil {
		log.Panicln(err.Error())
	}
	return FactoryConfig{
		MaxLimit:   constants.MAX_LIMIT,
		ErrorRate:  constants.ERROR_RATE,
		MaxTry:     constants.MAX_TRY,
		IdSize:     constants.ID_SIZE,
		BackupPath: bp,
	}
}

// Get backup path for bloom filter creating directories
// if they don't exist
func bloomPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	wormDir := path.Join(home, constants.DOT_DIR)
	_, err = os.Stat(wormDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(wormDir, constants.DIR_PERM); err != nil {
			return "", err
		}
	}
	return path.Join(wormDir, constants.BLOOM_DB), nil
}
