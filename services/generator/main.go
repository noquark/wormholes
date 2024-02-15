package main

import (
	"fmt"
	"lib/db"
	"lib/header"
	"net"
	"os"
	"protos"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dbconf := db.Load()
	postgres := dbconf.Postgres.Connect()

	header.Show("Generator")

	conf := DefaultConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start generator")
	}

	f := NewFactory(conf, postgres).Prepare().Run()
	grpcServer := grpc.NewServer()
	protos.RegisterBucketServiceServer(grpcServer, f)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to start generator")
	}
}
