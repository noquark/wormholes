package main

import (
	"flag"
	"os"
	"wormholes/internal/db"
	"wormholes/services/creator"
	"wormholes/services/director"
	"wormholes/services/generator"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	Generator = "generator"
	Creator   = "creator"
	Director  = "director"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	as := flag.String("as", Generator, "Run wormhole as")
	flag.Parse()

	dbConf := db.Load()
	pg := dbConf.Postgres.Connect()
	ts := dbConf.Timescale.Connect()
	redis := dbConf.Redis.Connect()

	db.InitPg(pg)
	db.InitTS(ts)

	switch *as {
	case Generator:
		generator.Run(pg)
	case Creator:
		creator.Run(pg, redis)
	case Director:
		director.Run(pg, ts, redis)
	default:
		log.Error().Msg("Available options are generator, creator and director")
	}
}
