package generator

import (
	"fmt"
	"net"
	"wormholes/internal/header"
	"wormholes/protos"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Run() {
	header.Show("Generator")

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
