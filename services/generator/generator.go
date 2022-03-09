package generator

import (
	"fmt"
	"net"
	"wormholes/internal/header"
	"wormholes/protos"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Run(pg *pgxpool.Pool) {
	header.Show("Generator")

	conf := DefaultConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start generator")
	}

	f := NewFactory(conf, pg).Prepare().Run()
	grpcServer := grpc.NewServer()
	protos.RegisterBucketServiceServer(grpcServer, f)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to start generator")
	}
}
