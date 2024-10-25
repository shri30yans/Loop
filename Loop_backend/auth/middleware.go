package auth

import (
    "fmt"
    "net/http"
    "github.com/golang-jwt/jwt/v5"
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Remove Bearer prefix
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        // Parse and validate token
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return JwtSecret, nil
        })

        if err != nil {
            if err == jwt.ErrTokenExpired {
                http.Error(w, "Token expired", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user ID to request context
        r = r.WithContext(SetUserContext(r.Context(), claims.UserID))
        next.ServeHTTP(w, r)
    }
}