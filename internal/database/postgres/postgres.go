package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

const INITIAL_SCHEMA = `
CREATE TABLE IF NOT EXISTS links (
  id varchar(255),
  tag varchar(255),
  target varchar(255),
  created_at timestamptz NULL,
  updated_at timestamptz NULL,
  CONSTRAINT links_pk PRIMARY KEY (id)
);
`

// Postgres connnection
type Postgres struct {
	Username string
	Password string
	Host     string
	Database string
	Port     int
}

func Default() Postgres {
	return Postgres{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     5432,
		Database: "postgres",
	}
}

func (p *Postgres) Connect() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.Username, p.Password, p.Host, p.Port, p.Database),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	rows, err := dbpool.Query(context.Background(), INITIAL_SCHEMA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create required tables: %v\n", err)
		os.Exit(1)
	}
	rows.Close()

	return dbpool
}
