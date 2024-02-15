package cache

import (
	"context"
	"lib/links"

	"github.com/mediocregopher/radix/v4"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	radix.Client
}

func New(uri string) *Cache {
	client, err := (radix.PoolConfig{}).New(context.Background(), "tcp", uri)
	if err != nil {
		log.Error().Err(err).Msg("cache: failed to connect")
	}
	return &Cache{
		client,
	}
}

func (c *Cache) GetLink(link *links.Link, shortID string) (err error) {
	err = c.Do(context.Background(), radix.Cmd(&link, "HGETALL", shortID))
	return err
}

func (c *Cache) SetLink(link links.Link, shortID string) (err error) {
	err = c.Do(context.Background(), radix.Cmd(nil, "HSET", shortID, "id", link.ID, "target", link.Target, "tag", link.Tag))
	return err
}
