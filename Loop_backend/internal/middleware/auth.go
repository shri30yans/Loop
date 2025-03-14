package middleware

import (
 "bytes"
 "context"
 "fmt"
 "io"
 "net/http"

 "Loop_backend/internal/response"
 "Loop_backend/internal/services"
)


type contextKey string

const UserIDKey contextKey = "userID"

// WithAuth is a middleware that ensures a user is authenticated
func WithAuth(next http.HandlerFunc, authService services.AuthService) http.HandlerFunc {
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

// LogPayload is a middleware that logs the request payload to the console
func LogPayload(next http.HandlerFunc) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  body, err := io.ReadAll(r.Body)
  if err != nil {
   fmt.Println("Error reading body:", err)
   next.ServeHTTP(w, r)
   return
  }

  // Restore the io.ReadCloser to its original state
  r.Body = io.NopCloser(bytes.NewBuffer(body))

  fmt.Printf("Request Payload: %s\n", string(body))
  next.ServeHTTP(w, r)
 }
}