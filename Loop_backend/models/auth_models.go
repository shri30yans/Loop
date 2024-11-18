package models

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)
// ----------------------------------------------------------------------------
// Authentication Structures
// ----------------------------------------------------------------------------

// Session represents a user session in the database.
type Session struct {
	SessionID    int       `json:"session_id"`
	UserID       int       `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// Login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Auth response
type AuthResponse struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type RegisterResponse struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

// Claims structure
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
