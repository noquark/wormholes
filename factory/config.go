package factory

import (
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

type Conf struct {
	MaxLimit   uint
	ErrorRate  float64
	MaxTry     int
	IdSize     int
	BackupPath string
}

func Default() Conf {
	bp, err := bloomPath()
	if err != nil {
		log.Panicln(err.Error())
	}
	return Conf{
		MaxLimit:   1e7,
		ErrorRate:  1e-3,
		MaxTry:     10,
		IdSize:     7,
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
	wormDir := path.Join(home, ".wormholes")
	_, err = os.Stat(wormDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(wormDir, 0775); err != nil {
			return "", err
		}
	}
	return path.Join(wormDir, "bloom.db"), nil
}
