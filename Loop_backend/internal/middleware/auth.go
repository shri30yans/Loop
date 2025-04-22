package middleware

import (
	"context"
	"net/http"

	"Loop_backend/internal/response"
	"Loop_backend/internal/services/auth"
)

// WithAuth is a middleware that ensures a user is authenticated
func WithAuth(next http.HandlerFunc, authService auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.RespondWithError(w, http.StatusUnauthorized, "No token provided")
			return
		}

		claims, err := authService.ValidateToken(authHeader)
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
