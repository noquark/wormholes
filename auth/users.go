package auth

import (
	"github.com/google/uuid"
)

// user model and constructor

type User struct {
	Id      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	IsAdmin bool      `json:"is_admin"`
}

func New(email string, isAdmin bool) *User {
	return &User{
		Email:   email,
		IsAdmin: isAdmin,
	}
}
