package db

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

//go:embed sql/postgres.sql
var pgSchema string

// Config for PostgreSQL.
type Postgres struct {
	URI      string `env:"PG_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres"`
	MaxConns int32  `env:"PG_MAX_CONN" envDefault:"5000"`
}

func (db *Postgres) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(db.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres: failed to parse config")
	}

	config.MaxConns = db.MaxConns

	dbpool, err := pgxpool.NewWithConfig(
		context.Background(), config,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres: failed to connect")
	}

	return dbpool
}

func InitPg(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), pgSchema)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres: failed to create required tables")
	}
}
