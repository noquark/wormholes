package config

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

//go:embed sql/postgres.sql
var POSTGRES string

//go:embed sql/timescale.sql
var TIMESCALE string

const TYPE_POSTGRES = "postgres"
const TYPE_TIMESCALE = "timescale"

// Postgres connection
type Postgres struct {
	DbType   string
	Username string
	Password string
	Host     string
	Database string
	Port     int
}

func DefaultPostgres() Postgres {
	return Postgres{
		DbType:   TYPE_POSTGRES,
		Username: "postgres",
		Password: "postgres",
		Host:     "postgres",
		Port:     5432,
		Database: "postgres",
	}
}

func DefaultTimescale() Postgres {
	return Postgres{
		DbType:   TYPE_TIMESCALE,
		Username: "postgres",
		Password: "postgres",
		Host:     "timescale",
		Port:     5432,
		Database: "postgres",
	}
}

func (p *Postgres) CreateTables(pool *pgxpool.Pool, schema string) {
	_, err := pool.Exec(context.Background(), schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create required tables: %v\n", err)
		os.Exit(1)
	}

}

func (p *Postgres) Connect() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.Username, p.Password, p.Host, p.Port, p.Database),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	if p.DbType == TYPE_POSTGRES {
		log.Info().Msg("setting max connections to 5000")
		config.MaxConns = 5000
	}
	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if p.DbType == TYPE_POSTGRES {
		p.CreateTables(dbpool, POSTGRES)
	} else {
		p.CreateTables(dbpool, TIMESCALE)
	}

	return dbpool
}
