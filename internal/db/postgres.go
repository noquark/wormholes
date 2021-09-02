package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

// Config for PostgreSQL.
type Postgres struct {
	URI      string `env:"PG_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres"`
	MaxConns int32  `env:"PG_MAX_CONN" envDefault:"5000"`
}

func (db *Postgres) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(db.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to postgres")
	}

	config.MaxConns = db.MaxConns

	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres: failed to connect")
	}

	return dbpool
}
