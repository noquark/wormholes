package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

// Config for TimescaleDB.
type Timescale struct {
	URI string `env:"TS_URI" envDefault:"postgres://postgres:postgres@localhost:5433/postgres"`
}

func (db *Timescale) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(db.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to timescale")
	}

	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to timescale: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
