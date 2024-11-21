package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	. "Loop/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type contextKey string

const UserContextKey contextKey = "user_id"

var JwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	if secretKey := os.Getenv("JWT_SECRET"); secretKey != "" {
		JwtSecret = []byte(secretKey)
	} else {
		panic("JWT_SECRET environment variable is not set")
	}
}

func GenerateJWT(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func SetUserContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserContextKey, userID)
}

func GetUserFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserContextKey).(int)
	return userID, ok
}
