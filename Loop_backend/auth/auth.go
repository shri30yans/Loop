// auth/utils.go
package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type contextKey string

const UserContextKey contextKey = "user"

func GenerateToken(userID int) (string, error) {
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

var (
	JwtSecret = []byte("your-secret-key") // In production, use environment variable
	//emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)