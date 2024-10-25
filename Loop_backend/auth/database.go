package auth

import (
	db "Loop/database"
	"context"
	"time"
)

func CreateUser(email string, hashedPassword string) (db.User, error) {
	var user db.User
	err := db.DB.QueryRow(
		context.Background(),
		"INSERT INTO users (email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id, email, created_at",
		email, hashedPassword, time.Now(),
	).Scan(&user.ID, &user.Email, &user.CreatedAt)
	return user, err
}

func GetUserByEmail(email string) (db.User, error) {
	var user db.User
	err := db.DB.QueryRow(
		context.Background(),
		"SELECT id, email, password_hash, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
