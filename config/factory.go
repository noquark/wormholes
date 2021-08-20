package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/mohitsinghs/wormholes/constants"
)

type FactoryConfig struct {
	MaxLimit                              uint
	ErrorRate                             float64
	MaxTry, IDSize, CookieSize, TokenSize int
	BackupPath                            string
}

func DefaultFactory() FactoryConfig {
	bp, err := bloomPath()
	if err != nil {
		log.Panicln(err.Error())
	}

	return FactoryConfig{
		MaxLimit:   constants.MaxLimit,
		ErrorRate:  constants.ErrorRate,
		MaxTry:     constants.MaxTry,
		IDSize:     constants.IDSize,
		CookieSize: constants.CookieSize,
		TokenSize:  constants.TokenSize,
		BackupPath: bp,
	}
}

// Get backup path for bloom filter creating directories
// if they don't exist.
func bloomPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %w", err)
	}

	wormDir := path.Join(home, constants.DotDir)

	_, err = os.Stat(wormDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(wormDir, constants.DirPerm); err != nil {
			return "", fmt.Errorf("failed to create required directories  : %w", err)
		}
	}

	return wormDir, nil
}
