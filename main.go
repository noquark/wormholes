package main

import (
	"flag"
	"fmt"
	"os"
	"wormholes/internal/db"
	"wormholes/internal/header"
	"wormholes/services/creator"
	"wormholes/services/creator/reserve"
	"wormholes/services/director"
	"wormholes/services/generator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	Generator = "generator"
	Creator   = "creator"
	Director  = "director"
	Unified   = "unified"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	as := flag.String("as", Unified, "Run wormhole as")
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
	case Unified:
		genConf := generator.DefaultConfig()
		creatorConf := creator.DefaultConfig()
		directorConf := director.DefaultConfig()

		factory := generator.NewFactory(genConf, pg).Prepare().Run()
		ingestor := creator.NewIngestor(pg, creatorConf.BatchSize).Start()
		reserve := reserve.WithLocal(factory)
		store := creator.NewPgStore(pg)
		pipe := director.NewPipe(directorConf, ts).Start().Wait()

		cHandler := creator.NewHandler(store, ingestor, redis, reserve)
		dHandler := director.NewHandler(pipe, pg, redis)

		app := fiber.New(fiber.Config{
			DisableStartupMessage:   true,
			EnableTrustedProxyCheck: true,
			ServerHeader:            "wormholes",
		})

		apiV1 := app.Group("v1")
		linksAPI := apiV1.Group("links")
		cHandler.Setup(linksAPI)

		redirectAPI := app.Group("l")
		redirectAPI.Get("/:id", dHandler.Redirect)

		app.Use(etag.New())
		app.Use(recover.New())

		header.Show("Unified")
		log.Info().Msgf("Running on port %d", directorConf.Port)

		if err := app.Listen(fmt.Sprintf(":%d", directorConf.Port)); err != nil {
			log.Error().Err(err).Msg("failed to start server")
		}
	}
}
