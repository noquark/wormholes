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

	dbConf := db.Load()
	postgres := dbConf.Postgres.Connect()
	timescale := dbConf.Timescale.Connect()
	redis := dbConf.Redis.Connect()

	db.InitPg(postgres)
	db.InitTS(timescale)

	switch dbConf.Mode {
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
