package creator

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed sql/update_link.sql
var sqlUpdate string

//go:embed sql/get_link.sql
var sqlGet string

//go:embed sql/delete_link.sql
var sqlDelete string

// postgres implementation of link db store.
type PgStore struct {
	db *pgxpool.Pool
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{
		db: pool,
	}
}

func (p *PgStore) Get(id string) (*Link, error) {
	var link Link

	rows := p.db.QueryRow(context.Background(),
		sqlGet,
		id,
	)

	err := rows.Scan(&link.ID, &link.Target, &link.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve link: %w", err)
	}

	return &link, nil
}

func (p *PgStore) Update(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		sqlUpdate,
		link.Target, link.Tag, link.ID,
	)
	if err != nil {
		log.Printf("Error updating link : %v", err)

		return fmt.Errorf("failed to update link: %w", err)
	}

	return nil
}

func (p *PgStore) Delete(id string) error {
	_, err := p.db.Exec(context.Background(),
		sqlDelete,
		id,
	)
	if err != nil {
		log.Printf("Error deleting link %v", err)

		return fmt.Errorf("failed to delete link: %w", err)
	}

	return nil
}
