package links

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// postgres implementation of link db store

type PgStore struct {
	db *pgxpool.Pool

	sqlInsert string
	sqlUpdate string
	sqlGet    string
	sqlDelete string
	sqlIds    string
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{
		db:        pool,
		sqlInsert: `INSERT INTO wh_links (id, target, tag, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);`,
		sqlUpdate: `UPDATE wh_links SET target=$1, tag=$2, updated_at=$3 where id=$4`,
		sqlGet:    `SELECT id, target, tag FROM wh_links where id=$1`,
		sqlDelete: `DELETE FROM wh_links WHERE id=$1`,
		sqlIds:    `SELECT id from links`,
	}
}

func (p *PgStore) Insert(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlInsert,
		link.Id, link.Target, link.Tag, time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("Error creating link : %v", err)
	}
	return err
}

func (p *PgStore) Update(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlUpdate,
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
		p.sqlGet,
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
		p.sqlDelete,
		id,
	)
	if err != nil {
		log.Printf("Error deleting link %v", err)
		return err
	}
	return nil
}

func (p *PgStore) Ids() ([]string, error) {
	rows, err := p.db.Query(context.Background(),
		p.sqlIds,
	)
	if err != nil {
		log.Printf("Error during ids query : %v", err)
		return nil, errors.New("error getting ids")
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
