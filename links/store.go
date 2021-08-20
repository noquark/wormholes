package links

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// postgres implementation of link db store

type Store struct {
	db *pgxpool.Pool

	sqlUpdate string
	sqlGet    string
	sqlDelete string
	sqlIds    string
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{
		db:        pool,
		sqlUpdate: `UPDATE wh_links SET target=$1, tag=$2 where id=$3`,
		sqlGet:    `SELECT id, target, tag FROM wh_links where id=$1`,
		sqlDelete: `DELETE FROM wh_links WHERE id=$1`,
		sqlIds:    `SELECT id from wh_links`,
	}
}

func (p *Store) Update(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlUpdate,
		link.Target, link.Tag, link.ID,
	)
	if err != nil {
		log.Printf("Error updating link : %v", err)

		return fmt.Errorf("failed to update link: %w", err)
	}

	return nil
}

func (p *Store) Get(id string) (*Link, error) {
	var link Link

	rows := p.db.QueryRow(context.Background(),
		p.sqlGet,
		id,
	)

	err := rows.Scan(&link.ID, &link.Target, &link.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve link: %w", err)
	}

	return &link, nil
}

func (p *Store) Delete(id string) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlDelete,
		id,
	)
	if err != nil {
		log.Printf("Error deleting link %v", err)

		return fmt.Errorf("failed to delete link: %w", err)
	}

	return nil
}

func (p *Store) Ids() ([]string, error) {
	rows, err := p.db.Query(context.Background(),
		p.sqlIds,
	)
	if err != nil {
		log.Printf("Error during ids query : %v", err)

		return nil, fmt.Errorf("failed to retrieve ids: %w", err)
	}
	defer rows.Close()

	ids := []string{}

	for rows.Next() {
		var id string

		err := rows.Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Ids not found : %v", err)

			return nil, fmt.Errorf("failed to find ids: %w", err)
		}

		if err != nil {
			log.Printf("Error getting id : %v", err)

			continue
		}

		ids = append(ids, id)
	}

	return ids, nil
}
