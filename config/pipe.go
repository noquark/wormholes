package config

import (
	"github.com/mohitsinghs/wormholes/constants"
)

type PipeConfig struct {
	Streams   int
	BatchSize int
}

func DefaultPipe() PipeConfig {
	return PipeConfig{
		Streams:   constants.Streams,
		BatchSize: constants.BatchSize,
	}
}
