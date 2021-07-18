package stats

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool

	sqlLinks  string
	sqlUsers  string
	sqlDBSize string
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{
		db:        pool,
		sqlLinks:  `select count(id) as links, count(distinct tag) as tags from wh_links`,
		sqlUsers:  `select count(*) as clicks, count(distinct cookie) as users from wh_clicks`,
		sqlDBSize: `select pg_size_pretty(pg_relation_size('wh_links')) links, pg_size_pretty(hypertable_size('wh_clicks')) as clicks`,
	}
}

func (s *Store) Overview() (*Overview, error) {
	var overview Overview
	rows := s.db.QueryRow(context.Background(),
		s.sqlLinks,
	)
	err := rows.Scan(&overview.Links, &overview.Tags)
	if err != nil {
		return nil, err
	}
	rows = s.db.QueryRow(context.Background(),
		s.sqlUsers,
	)
	err = rows.Scan(&overview.Clicks, &overview.Users)
	if err != nil {
		return nil, err
	}
	return &overview, nil
}

func (s *Store) DBSize() (*DBSize, error) {
	var dbSize DBSize
	rows := s.db.QueryRow(context.Background(),
		s.sqlDBSize,
	)
	err := rows.Scan(&dbSize.Links, &dbSize.Clicks)
	if err != nil {
		return nil, err
	}
	return &dbSize, nil
}
