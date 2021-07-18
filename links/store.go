package links

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// postgres implementation of link db store

type Store struct {
	db *pgxpool.Pool

	sqlInsert string
	sqlUpdate string
	sqlGet    string
	sqlDelete string
	sqlIds    string
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{
		db:        pool,
		sqlInsert: `INSERT INTO wh_links (id, target, tag) VALUES ($1, $2, $3);`,
		sqlUpdate: `UPDATE wh_links SET target=$1, tag=$2 where id=$3`,
		sqlGet:    `SELECT id, target, tag FROM wh_links where id=$1`,
		sqlDelete: `DELETE FROM wh_links WHERE id=$1`,
		sqlIds:    `SELECT id from links`,
	}
}

func (p *Store) Insert(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlInsert,
		link.Id, link.Target, link.Tag,
	)
	if err != nil {
		log.Printf("Error creating link : %v", err)
	}
	return err
}

func (p *Store) Update(link *Link) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlUpdate,
		link.Target, link.Tag, link.Id,
	)
	if err != nil {
		log.Printf("Error updating link : %v", err)
	}
	return err
}

func (p *Store) Get(id string) (*Link, error) {
	var link Link
	rows := p.db.QueryRow(context.Background(),
		p.sqlGet,
		id,
	)
	err := rows.Scan(&link.Id, &link.Target, &link.Tag)
	switch err {
	case pgx.ErrNoRows:
		log.Printf("link not found : %v", err)
		return nil, errors.New("link not found")
	case nil:
		return &link, nil
	default:
		log.Printf("Error getting link : %v", err)
		return nil, err
	}
}

func (p *Store) Delete(id string) error {
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

func (p *Store) Ids() ([]string, error) {
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
