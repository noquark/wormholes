package auth

import (
	"github.com/google/uuid"
)

// user model and constructor

type User struct {
	ID      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	IsAdmin bool      `json:"isAdmin"`
}

func New(email string, isAdmin bool) *User {
	return &User{
		Email:   email,
		IsAdmin: isAdmin,
	}
}
