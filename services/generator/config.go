package main

import (
	"wormholes/internal/db"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type BucketConfig struct {
	Size     int `env:"BUCKET_SIZE" envDefault:"8"`
	Capacity int `env:"BUCKET_CAPACITY" envDefault:"100000"`
}

type BloomConfig struct {
	MaxLimit  int     `env:"MAX_LIMIT" envDefault:"1000000"`
	ErrorRate float64 `env:"ERROR_RATE" envDefault:"0.0000001"`
}

type Config struct {
	Port   int `env:"PORT" envDefault:"5001"`
	IDSize int `env:"ID_SIZE" envDefault:"7"`
	BloomConfig
	BucketConfig
	db.Postgres
}

func DefaultConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("config: failed to parse")
	}

	return &cfg
}
