package auth

import (
	db "Loop/database"
	"Loop/models"
	"context"
	"fmt"
	"time"
)

const (
	SessionDuration = 7 * 24 * time.Hour
	TokenDuration   = 15 * time.Minute
)

func CreateSession(userID int) (*models.Session, error) {
	refreshToken, err := GenerateJWT(userID)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
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

func GetSessionByRefreshToken(refreshToken string) (*models.Session, error) {
	fmt.Println(refreshToken)
	session := &models.Session{}
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id, user_id, refresh_token, expires_at, created_at 
         FROM sessions 
         WHERE refresh_token = $1`,
		refreshToken,
	).Scan(&session.SessionID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt)
	return session, err
}

func GetSessionByUserID(userID int) (*models.Session, error) {
	session := &models.Session{}
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id, user_id, refresh_token, expires_at, created_at FROM sessions WHERE user_id = $1`,
		userID,
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

// func CleanupExpiredSessions() error {
// 	_, err := db.DB.Exec(
// 		context.Background(),
// 		"DELETE FROM sessions WHERE expires_at < NOW()",
// 	)
// 	return err
// }
