// auth/utils.go
package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var JwtSecret = []byte("your-secret-key")

type contextKey string

const UserContextKey contextKey = "user_id"

func GenerateJWT(userID int) (string, error) {
	claims := &Claims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
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

