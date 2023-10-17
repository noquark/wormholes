package store

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"wormholes/internal/links"
	"wormholes/services/creator/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

// postgres implementation of link db store.
type PgStore struct {
	db *pgxpool.Pool
}

func WithPg(pool *pgxpool.Pool) *PgStore {
	return &PgStore{
		db: pool,
	}
}

func (p *PgStore) Get(id string) (links.Link, error) {
	var link links.Link

	err := p.db.QueryRow(context.Background(),
		sql.Get,
		id,
	).Scan(&link.ID, &link.Target, &link.Tag)
	if err != nil {
		return links.Link{}, fmt.Errorf("failed to retrieve link: %w", err)
	}

	return link, nil
}

func (p *PgStore) Update(link *links.Link) error {
	_, err := p.db.Exec(context.Background(),
		sql.Update,
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
		sql.Delete,
		id,
	)
	if err != nil {
		log.Printf("Error deleting link %v", err)

		return fmt.Errorf("failed to delete link: %w", err)
	}

	return nil
}
