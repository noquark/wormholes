package main

import (
	"flag"
	"os"
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

	switch *as {
	case Generator:
		generator.Run()
	case Creator:
		creator.Run()
	case Director:
		director.Run()
	default:
		log.Error().Msg("Available options are generator, creator and director")
	}
}
