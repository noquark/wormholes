package creator

import (
	"fmt"
	"wormholes/internal/header"
	"wormholes/services/creator/reserve"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func Run(pg *pgxpool.Pool, redis *redis.Client) {
	conf := DefaultConfig()
	store := NewPgStore(pg)

	ingestor := NewIngestor(pg, conf.BatchSize).Start()
	reserve := reserve.WithGrpc(conf.GenAddr)

	handler := NewHandler(store, ingestor, redis, reserve)

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
