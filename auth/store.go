package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// postgres implementation of user db store

type Store struct {
	db *pgxpool.Pool

	sqlInsert     string
	sqlInsertSafe string
	sqlGet        string
	sqlDelete     string
	sqlValidate   string
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{
		db:            pool,
		sqlInsert:     `INSERT INTO wh_users (email, hashed_password, is_admin) VALUES ($1, $2, $3);`,
		sqlInsertSafe: `INSERT INTO wh_users (email, hashed_password, is_admin) VALUES ($1, $2, $3) on conflict do nothing;`,
		sqlGet:        `SELECT id, email, is_admin FROM wh_users where email=$1`,
		sqlDelete:     `DELETE FROM wh_users WHERE email=$1;`,
		sqlValidate:   `SELECT hashed_password, is_admin FROM wh_users WHERE email=$1`,
	}
}

func (p *Store) InsertSafe(user *User, hash string) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlInsertSafe, user.Email, hash, user.IsAdmin,
	)
	if err != nil {
		return fmt.Errorf("failed to create user safely: %w", err)
	}

	return nil
}

func (p *Store) Insert(user *User, hash string) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlInsert, user.Email, hash, user.IsAdmin,
	)
	if err != nil {
		log.Printf("Error creating user : %v", err)

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (p *Store) Get(email string) (*User, error) {
	var user User

	rows := p.db.QueryRow(context.Background(),
		p.sqlGet,
		email,
	)
	err := rows.Scan(&user.ID, &user.Email, &user.IsAdmin)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("User not found : %v", err)

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (p *Store) Delete(email string) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlDelete,
		email,
	)
	if err != nil {
		log.Printf("Error deleting user: %v", err)

		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (p *Store) ValidateAuth(email, password string) (bool, bool) {
	var hashedPassword string

	var isAdmin bool

	err := p.db.QueryRow(
		context.Background(),
		p.sqlValidate,
		email,
	).Scan(&hashedPassword, &isAdmin)
	if err != nil {
		log.Printf("Error getting hash %v", err)

		return false, false
	}

	err = CompareHashAndPassword(hashedPassword, []byte(password))

	return err == nil, isAdmin
}
