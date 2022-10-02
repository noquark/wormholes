package unified

import (
	"fmt"
	"wormholes/internal/header"
	"wormholes/services/creator"
	"wormholes/services/creator/reserve"
	"wormholes/services/director"
	"wormholes/services/generator"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func Run(
	postgres *pgxpool.Pool,
	timescale *pgxpool.Pool,
	redis *redis.Client,
) {
	genConf := generator.DefaultConfig()
	creatorConf := creator.DefaultConfig()
	directorConf := director.DefaultConfig()

	factory := generator.NewFactory(genConf, postgres).Prepare().Run()
	ingestor := creator.NewIngestor(postgres, creatorConf.BatchSize).Start()
	reserve := reserve.WithLocal(factory)
	cStore := creator.NewPgStore(postgres)
	pipe := director.NewPipe(directorConf, timescale).Start().Wait()

	cHandler := creator.NewHandler(cStore, ingestor, redis, reserve)
	dHandler := director.NewHandler(pipe, postgres, redis)

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
