package database

import "time"

// User represents a user in the database.
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	Location     string    `json:"location"`
	Bio          string    `json:"bio"`
	PasswordHash string    `json:"password_hash"`
}