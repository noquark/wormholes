package main

import (
	"fmt"
	"os"
	"wormholes/internal/header"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	conf := DefaultConfig()

	var store Store

	pool := conf.Postgres.Connect()
	store = NewPgStore(pool)

	ingestor := NewIngestor(pool, conf.BatchSize).Start()
	cache := conf.Redis.Connect()
	reserve := NewReserve(conf.GenAddr)

	handler := NewHandler(store, ingestor, cache, reserve)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		Prefork:                 true,
		ServerHeader:            "wormholes-creator",
	})

	if !fiber.IsChild() {
		header.Show("Creator")

		log.Info().Msgf("Running on port %d", conf.Port)
	}

	app.Use(etag.New())
	app.Use(recover.New())

	handler.Setup(app)

	if err := app.Listen(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Error().Err(err).Msg("error starting server")
	}
}
