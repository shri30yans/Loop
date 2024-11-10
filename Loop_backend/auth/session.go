package auth

import (
	db "Loop/database"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"fmt"
)

const (
	SessionDuration = 7 * 24 * time.Hour * 1000
	TokenDuration   = 15 * time.Minute
)

func GenerateJWT(userID int) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Set expiration time
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func CreateSession(userID int) (*Session, error) {
	refreshToken, err := GenerateJWT(userID)
	if err != nil {
		return nil, err
	}

	session := &Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(SessionDuration),
	}

	err = db.DB.QueryRow(
		context.Background(),
		`INSERT INTO sessions (user_id, refresh_token, expires_at) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, created_at`,
		session.UserID, session.RefreshToken, session.ExpiresAt,
	).Scan(&session.SessionID, &session.CreatedAt)

	return session, err
}

func GetSessionByRefreshToken(refreshToken string) (*Session, error) {
	session := &Session{}
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id, user_id, refresh_token, expires_at, created_at 
         FROM sessions 
         WHERE refresh_token = $1`,
		refreshToken,
	).Scan(&session.SessionID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt)
	return session, err
}

func DeleteSession(sessionID int) error {
	_, err := db.DB.Exec(
		context.Background(),
		"DELETE FROM sessions WHERE id = $1",
		sessionID,
	)
	return err
}

func CleanupExpiredSessions() error {
	_, err := db.DB.Exec(
		context.Background(),
		"DELETE FROM sessions WHERE expires_at < NOW()",
	)
	return err
}
