package db

import (
	"context"

	"github.com/mediocregopher/radix/v4"
	"github.com/rs/zerolog/log"
)

// Config for Redis.
type Redis struct {
	URI string `env:"REDIS_URI" envDefault:"redis://localhost:6379/0"`
}

func (r *Redis) Connect() radix.Client {
	client, err := (radix.PoolConfig{}).New(context.Background(), "tcp", r.URI)
	if err != nil {
		log.Error().Err(err).Msg("redis: failed to connect")
	}
	return client
}
