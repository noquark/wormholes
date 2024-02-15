package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Postgres struct {
	Username string `env:"PG_USERNAME" envDefault:"postgres"`
	Password string `env:"PG_PASSWORD" envDefault:"postgres"`
	Host     string `env:"PG_HOST" envDefault:"postgres"`
	Database string `env:"PG_DB" envDefault:"postgres"`
	Port     int    `env:"PG_PORT" envDefault:"5432"`
}

func (db *Postgres) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.Username, db.Password, db.Host, db.Port, db.Database),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to postgres")
	}

	dbpool, err := pgxpool.NewWithConfig(
		context.Background(), config,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to postgres")
	}

	return dbpool
}
