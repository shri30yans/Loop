package models

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token has expired")
    ErrMalformedToken = errors.New("malformed token")
)

type AuthenticatedUser struct {
    UserID string
    HashedPassword string
}

// Claims represents the JWT claims structure
type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// Session represents an authenticated user session
type Session struct {
    UserID    string     `json:"user_id"`
    Token     string    `json:"token"`
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt time.Time `json:"expires_at"`
}

// NewSession creates a new session instance
func NewSession(userID string, token string, duration time.Duration) *Session {
    now := time.Now()
    return &Session{
        UserID:    userID,
        Token:     token,
        CreatedAt: now,
        ExpiresAt: now.Add(duration),
    }
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
    return time.Now().After(s.ExpiresAt)
}