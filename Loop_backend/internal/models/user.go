package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        string  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrInvalidID       = errors.New("invalid user ID")
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidPassword = errors.New("invalid password")
)

func NewUser(email string, username string, bio string, location string) (*User, error) {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Bio:       bio,
		Location:  location,
		CreatedAt: now,
	}, nil
}
