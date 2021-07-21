package factory

import (
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mohitsinghs/wormholes/config"
)

type Factory struct {
	ID, Cookie, Token *Bloom
	conf              *config.FactoryConfig
}

func New(config *config.FactoryConfig) *Factory {
	return &Factory{
		ID:     NewBloom("id", config.BackupPath, config.MaxLimit, config.ErrorRate),
		Cookie: NewBloom("cookie", config.BackupPath, config.MaxLimit, config.ErrorRate),
		Token:  NewBloom("token", config.BackupPath, config.MaxLimit, config.ErrorRate),
		conf:   config,
	}
}

func (f *Factory) NewId() (string, error) {
	id, err := gonanoid.New(f.conf.IdSize)
	if err != nil || f.ID.Exists([]byte(id)) {
		id = f.failSafe(f.conf.IdSize, f.ID)
	}
	if id == "" {
		return "", errors.New("unable to generate valid id")
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

func (f *Factory) NewToken() string {
	token, err := gonanoid.New(f.conf.TokenSize)
	if err != nil || f.Token.Exists([]byte(token)) {
		token = f.failSafe(f.conf.TokenSize, f.Token)
	}
	f.Token.Add([]byte(token))
	return token
}

// Backup all bloom-filters
func (f *Factory) Backup() {
	f.ID.Backup()
	f.Cookie.Backup()
	f.Token.Backup()
}

// Restore all bloom-filters
func (f *Factory) Restore(restoreFunc func() ([]string, error)) {
	restored := f.ID.Restore()
	if !restored {
		f.ID.TryRestore(restoreFunc)
	}
	f.Cookie.Restore()
	f.Token.Restore()
}

// Try genrating new id at least specified times
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
