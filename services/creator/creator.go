package creator

import (
	"fmt"
	"wormholes/internal/header"
	"wormholes/services/creator/ingestor"
	"wormholes/services/creator/reserve"
	"wormholes/services/creator/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/zerolog/log"
)

func Run(pg *pgxpool.Pool, cache radix.Client) {
	conf := DefaultConfig()
	db := store.WithPg(pg)
	pipe := ingestor.New(pg, conf.BatchSize).Start()
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
