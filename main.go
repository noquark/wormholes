package main

import (
	"os"
	"wormholes/internal/db"
	"wormholes/internal/modes"
	"wormholes/internal/unified"
	"wormholes/services/creator"
	"wormholes/services/director"
	"wormholes/services/generator"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf := db.Load()
	postgres := conf.Postgres.Connect()
	timescale := conf.Timescale.Connect()
	redis := conf.Redis.Connect()

	db.InitPg(postgres)
	db.InitTS(timescale)

	switch conf.Mode {
	case modes.Generator:
		generator.Run(postgres)
	case modes.Creator:
		creator.Run(postgres, redis)
	case modes.Director:
		director.Run(postgres, timescale, redis)
	case modes.Unified:
		unified.Run(postgres, timescale, redis)
	}
}
