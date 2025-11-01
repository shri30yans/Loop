package middleware

import (
    "net/http"
)

// HandlerWithMethod wraps a handler with HTTP method validation
func HandlerWithMethod(method string, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != method {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        next(w, r)
    }
}

// HandlerWithAuth wraps a handler with authentication
func HandlerWithAuth(authMiddleware func(http.HandlerFunc) http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
    return authMiddleware(next)
}

// HandlerWithMethodAndAuth combines method validation and authentication
func HandlerWithMethodAndAuth(method string, authMiddleware func(http.HandlerFunc) http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
    return HandlerWithAuth(authMiddleware, HandlerWithMethod(method, next))
}
