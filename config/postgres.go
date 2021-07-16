package config

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed database.sql
var SCHEMA string

// Postgres connection
type Postgres struct {
	Username string
	Password string
	Host     string
	Database string
	Port     int
}

func DefaultPostgres() Postgres {
	return Postgres{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     5432,
		Database: "postgres",
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
	config.MaxConns = 200
	dbpool, err := pgxpool.ConnectConfig(
		context.Background(), config,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	_, err = dbpool.Exec(context.Background(), SCHEMA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create required tables: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
