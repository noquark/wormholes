package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// Config for Redis.
type Redis struct {
	URI string `env:"REDIS_URI" envDefault:"redis://:redis@localhost:6379/0"`
}

func (r *Redis) Connect() *redis.Client {
	opts, err := redis.ParseURL(r.URI)
	if err != nil {
		log.Error().Err(err).Msg("redis: failed to parse url")
	}

	return redis.NewClient(opts)
}
