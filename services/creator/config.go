package main

import (
	"wormholes/internal/db"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port      int    `env:"PORT" envDefault:"5002"`
	GenAddr   string `env:"GEN_ADDR" envDefault:"localhost:5001"`
	BatchSize int    `env:"BATCH_SIZE" envDefault:"10000"`
	db.Postgres
	db.Redis
}

func DefaultConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err)
	}

	return &cfg
}
