package director

import (
	"fmt"
	"wormholes/internal/header"

	"github.com/go-redis/redis/v8"
	json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func Run(pg *pgxpool.Pool, ts *pgxpool.Pool, redis *redis.Client) {
	conf := DefaultConfig()

	pipe := NewPipe(conf, ts).Start().Wait()
	handler := NewHandler(pipe, pg, redis)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		Prefork:                 true,
		ServerHeader:            "wormholes",
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		GETOnly:                 true,
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
