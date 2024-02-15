package main

import (
	"creator/ingestor"
	"creator/reserve"
	"creator/store"
	"fmt"
	"lib/cache"
	"lib/db"
	"lib/header"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dbconf := db.Load()

	postgres := dbconf.Postgres.Connect()
	cache := cache.New(dbconf.REDIS_URI)

	conf := DefaultConfig()
	db := store.WithPg(postgres)
	pipe := ingestor.New(postgres, conf.BatchSize).Start()
	reserve := reserve.WithGrpc(conf.GenAddr)
	handler := NewHandler(db, pipe, cache, reserve)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		Prefork:                 true,
		ServerHeader:            "wormholes",
	})

	if !fiber.IsChild() {
		header.Show("Creator")

		log.Info().Msgf("Running on port %d", conf.Port)
	}

	app.Use(etag.New())
	app.Use(recover.New())

	handler.Setup(app)

	if err := app.Listen(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
