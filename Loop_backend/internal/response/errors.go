package response

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

// Error types
const (
	ErrorTypeValidation   ErrorType = "validation_error"
	ErrorTypeAuth         ErrorType = "authentication_error"
	ErrorTypePermission   ErrorType = "permission_error"
	ErrorTypeNotFound     ErrorType = "not_found_error"
	ErrorTypeInternal     ErrorType = "internal_error"
	ErrorTypeDatabase     ErrorType = "database_error"
	ErrorTypeExternalAPI  ErrorType = "external_api_error"
	ErrorTypeRateLimit    ErrorType = "rate_limit_error"
	ErrorTypeInvalidInput ErrorType = "invalid_input_error"
	ErrorTypeSearch       ErrorType = "search_error"
)

// AppError represents an application error
type AppError struct {
	Type       ErrorType              `json:"error_type"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	HTTPStatus int                    `json:"-"` // Not serialized
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewAppError creates a new AppError
func NewAppError(errorType ErrorType, message string, httpStatus int) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// RespondWithAppError sends an AppError response
func RespondWithAppError(w http.ResponseWriter, err *AppError) {
	RespondWithJSON(w, err.HTTPStatus, err)
}

// Common error constructors

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return NewAppError(ErrorTypeValidation, message, http.StatusBadRequest)
}

// NewAuthError creates an authentication error
func NewAuthError(message string) *AppError {
	return NewAppError(ErrorTypeAuth, message, http.StatusUnauthorized)
}

// NewPermissionError creates a permission error
func NewPermissionError(message string) *AppError {
	return NewAppError(ErrorTypePermission, message, http.StatusForbidden)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return NewAppError(ErrorTypeNotFound, message, http.StatusNotFound)
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *AppError {
	return NewAppError(ErrorTypeInternal, message, http.StatusInternalServerError)
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string) *AppError {
	return NewAppError(ErrorTypeDatabase, message, http.StatusInternalServerError)
}

// NewExternalAPIError creates an external API error
func NewExternalAPIError(message string) *AppError {
	return NewAppError(ErrorTypeExternalAPI, message, http.StatusBadGateway)
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(message string) *AppError {
	return NewAppError(ErrorTypeRateLimit, message, http.StatusTooManyRequests)
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string) *AppError {
	return NewAppError(ErrorTypeInvalidInput, message, http.StatusBadRequest)
}

// NewSearchError creates a search error
func NewSearchError(message string) *AppError {
	return NewAppError(ErrorTypeSearch, message, http.StatusInternalServerError)
}

// HandleError handles an error and sends an appropriate response
// It converts regular errors to AppErrors with sensible defaults
func HandleError(w http.ResponseWriter, err error) {
	// If it's already an AppError, use it directly
	if appErr, ok := err.(*AppError); ok {
		RespondWithAppError(w, appErr)
		return
	}

	// Otherwise, create a generic internal error
	// In production, you might want to log the original error but not expose it
	appErr := NewInternalError("An unexpected error occurred")
	RespondWithAppError(w, appErr)
}