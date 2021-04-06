package links

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// postgres implementaion of link db store

type PgStore struct {
	db *pgxpool.Pool
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{
		db: pool,
	}
}

func (p *PgStore) Insert(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		`INSERT INTO links.links (id, target, tag, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);`,
		link.Id, link.Target, link.Tag, time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("Error creating link : %v", err)
	}
	return err
}

func (p *PgStore) Update(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		`UPDATE links.links SET target=$1, tag=$2, updated_at=$3 where id=$4`,
		link.Target, link.Tag, time.Now(), link.Id,
	)
	if err != nil {
		log.Printf("Error updating link : %v", err)
	}
	return err
}

func (p *PgStore) Get(id string) (*Link, error) {
	var link Link
	rows := p.db.QueryRow(context.Background(),
		`SELECT id, target, tag FROM links.links where id=$1`,
		id,
	)
	err := rows.Scan(&link.Id, &link.Target, &link.Tag)
	switch err {
	case pgx.ErrNoRows:
		log.Printf("Link not found : %v", err)
		return nil, errors.New("Link not found")
	case nil:
		return &link, nil
	default:
		log.Printf("Error getting link : %v", err)
		return nil, err
	}
}

func (p *PgStore) Delete(id string) error {
	_, err := p.db.Exec(context.Background(),
		`DELETE FROM links.links WHERE id=$1;`,
		id,
	)
	if err != nil {
		log.Printf("Error deleting link %v", err)
		return err
	}
	return nil
}

func (s *PgStore) Ids() ([]string, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT id from links.links;`,
	)
	if err != nil {
		log.Printf("Error during ids query : %v", err)
		return nil, errors.New("Error getting ids")
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		switch err {
		case pgx.ErrNoRows:
			log.Printf("Ids not found : %v", err)
			return nil, errors.New("Ids not found")
		case nil:
			ids = append(ids, id)
		default:
			log.Printf("Error getting id : %v", err)
			continue
		}
	}
	return ids, nil
}
