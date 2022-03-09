package db

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Postgres
	Timescale
	Redis
}

func Load() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("db_config: failed to parse")
	}

	return &cfg
}
