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

	db := conf.Postgres.Connect()
	tsdb := conf.Timescale.Connect()
	cache := conf.Redis.Connect()
	pipe := NewPipe(conf, tsdb).Start().Wait()
	handler := NewHandler(pipe, db, cache)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		Prefork:                 true,
		ServerHeader:            "wormholes-director",
	})

	if !fiber.IsChild() {
		header.Show("Director")

		log.Info().Msgf("Running on port %d", conf.Port)
	}

	app.Use(etag.New())
	app.Use(recover.New())

	app.Get("/:id", handler.Redirect)

	err := app.Listen(fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Error().Err(err).Msg("error starting server")
	}
}
