package auth

import (
	"github.com/google/uuid"
)

// user model and constructor

type User struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func New(email string) *User {
	return &User{
		Email: email,
	}
}
