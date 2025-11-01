package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
    jwtSecret []byte
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
    return &AuthMiddleware{
        jwtSecret: []byte(jwtSecret),
    }
}

// Claims represents JWT claims
type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// Authenticate validates JWT tokens in requests
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "No token provided", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return m.jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        // Parse user ID from claims
        userID := claims.UserID
        if userID == "" {
            http.Error(w, "User ID not found in token", http.StatusUnauthorized)
            return
        }

        // Add user ID to request context using the defined key
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
