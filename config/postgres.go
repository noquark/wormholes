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

const (
	TypePostgres         = "postgres"
	TypeTimescale        = "timescale"
	DefaultPostgresPort  = 5432
	DefaultTimescalePort = 5432
	DefaultMaxConn       = 5000
)

// Postgres connection.
type Postgres struct {
	DBType   string
	Username string
	Password string
	Host     string
	Database string
	Port     int
}

func DefaultPostgres() Postgres {
	return Postgres{
		DBType:   TypePostgres,
		Username: "postgres",
		Password: "postgres",
		Host:     "postgres",
		Port:     DefaultPostgresPort,
		Database: "postgres",
	}
}

func DefaultTimescale() Postgres {
	return Postgres{
		DBType:   TypeTimescale,
		Username: "postgres",
		Password: "postgres",
		Host:     "timescale",
		Port:     DefaultTimescalePort,
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

	if p.DBType == TypePostgres {
		log.Info().Msg("setting max connections to 5000")

		config.MaxConns = DefaultMaxConn
	}

	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if p.DBType == TypePostgres {
		p.CreateTables(dbpool, POSTGRES)
	} else {
		p.CreateTables(dbpool, TIMESCALE)
	}

	return dbpool
}
