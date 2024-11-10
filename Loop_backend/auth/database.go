package auth

import (
	db "Loop/database"
	"context"
	"strings"
	"errors"
)

var ErrDuplicateEmail = errors.New("email already exists")

func CreateUser(name string, email string, hashedPassword string) (db.User, error) {
	var user db.User
	err := db.DB.QueryRow(
		context.Background(),
		"INSERT INTO users (name, email, hashed_password) VALUES ($1, $2, $3) RETURNING id,name,email, hashed_password",
		name, email, hashedPassword,
	).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword)
	
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return db.User{}, ErrDuplicateEmail
		}
		return db.User{}, err // return other types of errors as is
	}
	return user, nil
}

func GetUserByEmail(email string) (db.User, error) {
	var user db.User
	err := db.DB.QueryRow(
		context.Background(),
		"SELECT id, email, hashed_password, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt)
	return user, err
}
