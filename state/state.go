package state

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mohitsinghs/wormholes/config"
	"github.com/mohitsinghs/wormholes/factory"
	"github.com/mohitsinghs/wormholes/pipe"
)

type State struct {
	Config  *config.Config
	Factory *factory.Factory
	Pipe    *pipe.Pipe
	DB      *pgxpool.Pool
}

func New(config *config.Config, factory *factory.Factory, pipe *pipe.Pipe) *State {
	return &State{
		Config:  config,
		Factory: factory,
		Pipe:    pipe,
		DB:      config.Postgres.Connect(),
	}
}
