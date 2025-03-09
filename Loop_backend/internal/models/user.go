package models

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

// User represents a user in the system
type User struct {
    ID        string      
    Username  string    
    Email     string    
    Bio       string
    Location string
    CreatedAt time.Time 
    UpdatedAt time.Time 
}

var (
    ErrInvalidID       = errors.New("invalid user ID")
    ErrInvalidEmail    = errors.New("invalid email address")
    ErrInvalidPassword = errors.New("invalid password")
)

func NewUser(email string, username string , bio string, location string) (*User, error) {
    if email == "" {
        return nil, ErrInvalidEmail
    }

    now := time.Now()
    return &User{
        ID: uuid.New().String(),
        Email:     email,
        Username:  username,
        Bio : bio,
        Location: location,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}