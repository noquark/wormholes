package unified

import (
	"fmt"
	"wormholes/internal/header"
	"wormholes/services/creator"
	"wormholes/services/creator/ingestor"
	"wormholes/services/creator/reserve"
	"wormholes/services/creator/store"
	"wormholes/services/generator"
	"wormholes/services/redirector"
	"wormholes/services/redirector/pipe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/zerolog/log"
)

func Run(
	postgres *pgxpool.Pool,
	timescale *pgxpool.Pool,
	cache radix.Client,
) {
	gConf := generator.DefaultConfig()
	cConf := creator.DefaultConfig()
	dConf := redirector.DefaultConfig()

	factory := generator.NewFactory(gConf, postgres).Prepare().Run()
	ingest := ingestor.New(postgres, cConf.BatchSize).Start()
	reserve := reserve.WithLocal(factory)
	cStore := store.WithPg(postgres)
	pipe := pipe.New(dConf.BatchSize, dConf.Streams, timescale).Start().Wait()

	cHandler := creator.NewHandler(cStore, ingest, cache, reserve)
	dHandler := redirector.NewHandler(pipe, postgres, cache)

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
	log.Info().Msgf("Running on port %d", dConf.Port)

	if err := app.Listen(fmt.Sprintf(":%d", dConf.Port)); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
