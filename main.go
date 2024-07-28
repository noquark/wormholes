package main

import (
	"fmt"
	"net"
	"os"
	"wormholes/ingestor"
	"wormholes/internal/cache"
	"wormholes/internal/config"
	"wormholes/internal/db"
	"wormholes/internal/header"
	"wormholes/ipc"
	"wormholes/protos"
	"wormholes/store"

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
	conf := config.DefaultConfig()
	if !fiber.IsChild() {
		header.Show()

		log.Info().Msgf("Running on port %d", conf.Port)
	}
	dbconf := db.Load()

	postgres := dbconf.Postgres.Connect()
	cache := cache.New(dbconf.REDIS_URI)
	db.InitPg(postgres)

	backend := store.WithPg(postgres)
	pipe := ingestor.New(postgres, conf.BatchSize).Start()

	if !fiber.IsChild() {
		go func() {
			factory := ipc.NewFactory(conf, postgres).Prepare().Run(conf)
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GenPort))
			if err != nil {
				log.Fatal().Err(err).Msg("factory: failed to start")
			}

			grpcServer := grpc.NewServer()
			protos.RegisterBucketServiceServer(grpcServer, factory)

			if err := grpcServer.Serve(lis); err != nil {
				log.Fatal().Err(err).Msg("factory: failed to start")
			}
		}()
	}

	ipcStore := ipc.NewStore(fmt.Sprintf(":%d", conf.GenPort))
	handler := NewHandler(backend, pipe, cache, ipcStore)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		Prefork:                 true,
		ServerHeader:            "wormholes",
	})

	app.Use(etag.New())
	app.Use(recover.New())

	handler.Setup(app)

	if err := app.Listen(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
