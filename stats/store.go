package stats

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db, tsdb *pgxpool.Pool

	sqlLinks string
	sqlUsers string
}

func NewStore(db, tsdb *pgxpool.Pool) Store {
	return Store{
		db:       db,
		tsdb:     tsdb,
		sqlLinks: `select count(id) as links, count(distinct tag) as tags from wh_links`,
		sqlUsers: `select count(*) as clicks, count(distinct cookie) as users from wh_clicks`,
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
	rows = s.tsdb.QueryRow(context.Background(),
		s.sqlUsers,
	)
	err = rows.Scan(&overview.Clicks, &overview.Users)
	if err != nil {
		return nil, err
	}
	return &overview, nil
}
