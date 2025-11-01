package middleware

// ContextKey is used for context value keys
type ContextKey string

// Context keys
const (
    UserIDKey       ContextKey = "user_id"
    ValidatedDTOKey ContextKey = "validatedDTO"
)
