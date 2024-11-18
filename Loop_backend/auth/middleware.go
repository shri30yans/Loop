// auth/middleware.go
package auth

import (
	. "Loop/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Remove Bearer prefix
		if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			if err == jwt.ErrTokenExpired {
				// Check if refresh token is present
				refreshToken := r.Header.Get("X-Refresh-Token")
				if refreshToken != "" {
					if session, err := GetSessionByRefreshToken(refreshToken); err == nil && session.ExpiresAt.After(time.Now()) {
						newToken, err := GenerateJWT(session.UserID)
						if err == nil {
							w.Header().Set("X-New-Token", newToken)
							r = r.WithContext(SetUserContext(r.Context(), session.UserID))
							next.ServeHTTP(w, r)
							return
						} else {
							fmt.Printf("Error generating new token: %v\n", err)
						}
					} else {
						fmt.Println("Invalid or expired refresh token")
					}
				}
			} else {
				fmt.Println("Invalid token")
			}
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed with request
		r = r.WithContext(SetUserContext(r.Context(), claims.UserID))
		next.ServeHTTP(w, r)
	}
}
