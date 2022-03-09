package db

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

//go:embed sql/postgres.sql
var tsSchema string

// Config for TimescaleDB.
type Timescale struct {
	URI string `env:"TS_URI" envDefault:"postgres://postgres:postgres@localhost:5433/postgres"`
}

func (db *Timescale) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(db.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("timescale: failed to connect")
	}

	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("timescale: failed to connect")
	}

	return dbpool
}

func InitTS(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tsSchema)
	if err != nil {
		log.Fatal().Err(err).Msg("timescale: failed to create required tables")
	}
}
