package middleware

// ContextKey is used for context value keys
type ContextKey string

// Context keys
const (
	UserIDKey       ContextKey = "userID"
	ValidatedDTOKey ContextKey = "validatedDTO"
)
