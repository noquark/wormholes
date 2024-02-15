package main

import (
	"fmt"
	"lib/cache"
	"lib/db"
	"lib/header"
	"os"
	"redirector/pipe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dbconf := db.Load()

	postgres := dbconf.Postgres.Connect()
	timescale := dbconf.Timescale.Connect()
	cache := cache.New(dbconf.REDIS_URI)

	conf := DefaultConfig()

	pipe := pipe.New(conf.BatchSize, conf.Streams, timescale).Start().Wait()
	handler := NewHandler(pipe, postgres, cache)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		Prefork:                 true,
		ServerHeader:            "wormholes",
	})

	if !fiber.IsChild() {
		header.Show("Redirector")

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

func Run(pg *pgxpool.Pool, ts *pgxpool.Pool, cache *cache.Cache) {

}
