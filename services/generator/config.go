package main

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port           int           `env:"PORT" envDefault:"5001"`
	IDSize         int           `env:"ID_SIZE" envDefault:"7"`
	BucketSize     int           `env:"BUCKET_SIZE" envDefault:"16"`
	BucketCapacity int           `env:"BUCKET_CAP" envDefault:"100000"`
	BloomMaxLimit  uint          `env:"BLOOM_MAX" envDefault:"100000000"`
	BloomErrorRate float64       `env:"BLOOM_ERROR" envDefault:"0.0000001"`
	Timeout        time.Duration `env:"TIMEOUT" envDefault:"100ms"`
}

func DefaultConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("config: failed to parse")
	}

	return &cfg
}
