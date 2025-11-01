package middleware

import (
	"Loop_backend/internal/response"
	"net/http"
	"os"
	"strings"
)

// APITokenMiddleware handles API token validation for tool endpoints
type APITokenMiddleware struct {
	apiToken string
}

// NewAPITokenMiddleware creates a new APITokenMiddleware
func NewAPITokenMiddleware(apiToken string) *APITokenMiddleware {
	// If no token is provided, try to get it from environment
	if apiToken == "" {
		apiToken = os.Getenv("GO_BACKEND_API_TOKEN")
	}

	return &APITokenMiddleware{
		apiToken: apiToken,
	}
}

// ValidateAPIToken validates the API token in requests
func (m *APITokenMiddleware) ValidateAPIToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If no API token is configured, log a warning and proceed
		if m.apiToken == "" {
			// In a production environment, you might want to reject the request instead
			// log.Warn("No API token configured for tool endpoints")
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "No token provided")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(w, "Invalid token format")
			return
		}

		// Simple token comparison for API token validation
		if tokenString != m.apiToken {
			response.Unauthorized(w, "Invalid API token")
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}