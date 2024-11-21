package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateAndProcessToken validates a JWT and optionally handles refresh tokens
func ValidateAndProcessToken(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("no token provided")
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		// Handle expired token and refresh logic
		if err == jwt.ErrTokenExpired {
			refreshToken := r.Header.Get("X-Refresh-Token")
			if refreshToken != "" {
				session, err := GetSessionByRefreshToken(refreshToken)
				if err == nil && session.ExpiresAt.After(time.Now()) {
					newToken, err := GenerateJWT(session.UserID)
					if err == nil {
						w.Header().Set("X-New-Token", newToken)
						return jwt.MapClaims{"user_id": float64(session.UserID)}, nil
					}
				}
			}
		}
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=== Starting Auth Middleware ===")

		claims, err := ValidateAndProcessToken(w, r)
		if err != nil {
			fmt.Printf("❌ Token validation failed: %v\n", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		fmt.Printf("✅ Valid token for user ID: %d\n", userID)

		r = r.WithContext(SetUserContext(r.Context(), userID))

		next(w, r)
	}
}
