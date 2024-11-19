// auth/middleware.go
package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=== Starting Auth Middleware ===")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			fmt.Println("❌ No Authorization header found")
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		fmt.Printf("🔍 Processing token: %s\n", tokenString)
		fmt.Printf("🔑 Using JWT_SECRET: %s\n", os.Getenv("JWT_SECRET"))

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			fmt.Printf("📝 Token Header: %v\n", token.Header)
			fmt.Printf("🔐 Signing Method: %T\n", token.Method)

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("❌ Invalid signing method: %v\n", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			fmt.Printf("❌ Token validation failed: %v\n", err)
			if err == jwt.ErrTokenExpired {
				fmt.Println("🕒 Token is expired, checking refresh token...")
				refreshToken := r.Header.Get("X-Refresh-Token")
				if refreshToken != "" {
					fmt.Printf("🔄 Found refresh token: %s\n", refreshToken)
					if session, err := GetSessionByRefreshToken(refreshToken); err == nil && session.ExpiresAt.After(time.Now()) {
						fmt.Printf("✅ Valid session found for user: %d\n", session.UserID)
						newToken, err := GenerateJWT(session.UserID)
						if err == nil {
							fmt.Printf("✅ Generated new token: %s\n", newToken)
							w.Header().Set("X-New-Token", newToken)
							r = r.WithContext(SetUserContext(r.Context(), session.UserID))
							next(w, r)
							return
						}
						fmt.Printf("❌ Failed to generate new token: %v\n", err)
					}
					fmt.Printf("❌ Refresh token validation failed: %v\n", err)
				}
			}
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("❌ Failed to parse token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		fmt.Printf("📄 Token claims: %+v\n", claims)

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			fmt.Printf("❌ Invalid user_id type in claims. Value: %v, Type: %T\n", claims["user_id"], claims["user_id"])
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)
		if !ok {
			fmt.Printf("❌ Invalid user_id type in claims. Value: %v, Type: %T\n", claims["user_id"], claims["user_id"])
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		fmt.Printf("✅ Successfully validated token for user: %d\n", userID)
		r = r.WithContext(SetUserContext(r.Context(), userID))
		next(w, r)
	}
}
