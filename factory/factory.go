package factory

import (
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mohitsinghs/wormholes/config"
)

var ErrGenerateID = errors.New("failed to generate valid id")

type Factory struct {
	ID, Cookie, Token *Bloom
	conf              *config.FactoryConfig
}

func New(config *config.FactoryConfig) *Factory {
	return &Factory{
		ID:     NewBloom("id", config.BackupPath, config.MaxLimit, config.ErrorRate),
		Cookie: NewBloom("cookie", config.BackupPath, config.MaxLimit, config.ErrorRate),
		conf:   config,
	}
}

func (f *Factory) NewID() (string, error) {
	id, err := gonanoid.New(f.conf.IDSize)
	if err != nil || f.ID.Exists([]byte(id)) {
		id = f.failSafe(f.conf.IDSize, f.ID)
	}

	if id == "" {
		return "", ErrGenerateID
	}

	f.ID.Add([]byte(id))

	return id, nil
}

func (f *Factory) NewCookie() string {
	cookie, err := gonanoid.New(f.conf.CookieSize)
	if err != nil || f.Cookie.Exists([]byte(cookie)) {
		cookie = f.failSafe(f.conf.CookieSize, f.Cookie)
	}

	f.Cookie.Add([]byte(cookie))

	return cookie
}

// Backup all bloom-filters.
func (f *Factory) Backup() {
	f.ID.Backup()
	f.Cookie.Backup()
}

// Restore all bloom-filters.
func (f *Factory) Restore(restoreFunc func() ([]string, error)) {
	restored := f.ID.Restore()
	if !restored {
		f.ID.TryRestore(restoreFunc)
	}

	f.Cookie.Restore()
}

// Try genrating new id at least specified times.
func (f *Factory) failSafe(size int, bloom *Bloom) string {
	id := ""
	for i := 0; i < f.conf.MaxTry; i++ {
		id, err := gonanoid.New(size)
		if err != nil || bloom.Exists([]byte(id)) {
			continue
		}

		bloom.Add([]byte(id))

		break
	}

	return id
}
