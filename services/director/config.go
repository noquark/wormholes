package director

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

// Top Level Config.
type Config struct {
	Port      int `env:"PORT" envDefault:"5000"`
	BatchSize int `env:"BATCH_SIZE" envDefault:"10000"`
	Streams   int `env:"STREAMS" envDefault:"8"`
}

func DefaultConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Panic().Err(err)
	}

	return &cfg
}
