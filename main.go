package main

import (
	"os"
	"wormholes/internal/cache"
	"wormholes/internal/db"
	"wormholes/internal/modes"
	"wormholes/internal/unified"
	"wormholes/services/creator"
	"wormholes/services/generator"
	"wormholes/services/redirector"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf := db.Load()
	postgres := conf.Postgres.Connect()
	timescale := conf.Timescale.Connect()
	cache := cache.New(conf.REDIS_URI)

	db.InitPg(postgres)
	db.InitTS(timescale)

	switch conf.Mode {
	case modes.Generator:
		generator.Run(postgres)
	case modes.Creator:
		creator.Run(postgres, cache)
	case modes.Redirector:
		redirector.Run(postgres, timescale, cache)
	case modes.Unified:
		unified.Run(postgres, timescale, cache)
	}
}
