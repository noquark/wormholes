package redirector

import (
	"fmt"
	"wormholes/internal/cache"
	"wormholes/internal/header"
	"wormholes/services/redirector/pipe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func Run(pg *pgxpool.Pool, ts *pgxpool.Pool, cache *cache.Cache) {
	conf := DefaultConfig()

	pipe := pipe.New(conf.BatchSize, conf.Streams, ts).Start().Wait()
	handler := NewHandler(pipe, pg, cache)

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
