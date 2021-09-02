package main

import (
	"fmt"
	"net"
	"os"
	"wormholes/internal/header"
	"wormholes/protos"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	header.Show("Generator")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	conf := DefaultConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start generator")
	}

	f := NewFactory(conf).Prepare().Run()
	grpcServer := grpc.NewServer()
	protos.RegisterBucketServiceServer(grpcServer, f)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to start generator")
	}
}
