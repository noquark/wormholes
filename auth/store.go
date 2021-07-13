package auth

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// postgres implementation of user db store

type Store struct {
	db *pgxpool.Pool

	sqlInsert   string
	sqlGet      string
	sqlDelete   string
	sqlValidate string
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{
		db:          pool,
		sqlInsert:   `INSERT INTO wh_users (email, hashed_password) VALUES ($1, $2);`,
		sqlGet:      `SELECT id, email FROM wh_users where email=$1`,
		sqlDelete:   `DELETE FROM wh_users WHERE email=$1;`,
		sqlValidate: `SELECT hashed_password FROM wh_users WHERE email=$1`,
	}
}

func (p *Store) Insert(user *User, hash string) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlInsert, user.Email, hash,
	)
	if err != nil {
		log.Printf("Error creating user : %v", err)
	}
	return err
}

func (p *Store) Get(email string) (*User, error) {
	var user User
	rows := p.db.QueryRow(context.Background(),
		p.sqlGet,
		email,
	)
	err := rows.Scan(&user.Id, &user.Email)
	switch err {
	case pgx.ErrNoRows:
		log.Printf("User not found : %v", err)
		return nil, errors.New("User not found")
	case nil:
		return &user, nil
	default:
		log.Printf("Error getting user : %v", err)
		return nil, err
	}
}

func (p *Store) Delete(email string) error {
	_, err := p.db.Exec(context.Background(),
		p.sqlDelete,
		email,
	)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
	}
	return err
}

func (p *Store) ValidateAuth(email, password string) bool {
	var hashedPassword string
	err := p.db.QueryRow(
		context.Background(),
		p.sqlValidate,
		email,
	).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error getting hash %v", err)
		return false
	}
	err = CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}
