package main

import (
	"fmt"
	"os"
	"wormholes/internal/header"
	"wormholes/protos"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	conf := DefaultConfig()

	conn, err := grpc.Dial(conf.GenAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("grpc: failed to connect")
	}
	defer conn.Close()
	client := protos.NewBucketServiceClient(conn)

	var store Store

	pool := conf.Postgres.Connect()
	store = NewPgStore(pool)

	ingestor := NewIngestor(pool, conf.BatchSize)
	cache := conf.Redis.Connect()

	handler := NewHandler(store, ingestor, cache, client)

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
