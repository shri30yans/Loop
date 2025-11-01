package response

import (
    "encoding/json"
    "net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// RespondWithJSON sends a JSON response with status code
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// RespondWithError sends an error response with status code
func RespondWithError(w http.ResponseWriter, status int, message string) {
    RespondWithJSON(w, status, ErrorResponse{
        Message: message,
    })
}

// RespondWithErrorDetails sends an error response with additional details
func RespondWithErrorDetails(w http.ResponseWriter, status int, message string, details map[string]interface{}) {
    RespondWithJSON(w, status, ErrorResponse{
        Message: message,
        Details: details,
    })
}

// JSON is an alias for RespondWithJSON for new code
func JSON(w http.ResponseWriter, status int, data interface{}) {
    RespondWithJSON(w, status, data)
}

// Error is an alias for RespondWithError for new code
func Error(w http.ResponseWriter, status int, message string) {
    RespondWithError(w, status, message)
}

// Success sends an OK response
func Success(w http.ResponseWriter, data interface{}) {
    RespondWithJSON(w, http.StatusOK, data)
}

// Created sends a created response
func Created(w http.ResponseWriter, data interface{}) {
    RespondWithJSON(w, http.StatusCreated, data)
}

// NoContent sends a no content response
func NoContent(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a bad request response
func BadRequest(w http.ResponseWriter, message string) {
    RespondWithError(w, http.StatusBadRequest, message)
}

// Unauthorized sends an unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
    RespondWithError(w, http.StatusUnauthorized, message)
}

// Forbidden sends a forbidden response
func Forbidden(w http.ResponseWriter, message string) {
    RespondWithError(w, http.StatusForbidden, message)
}

// NotFound sends a not found response
func NotFound(w http.ResponseWriter, message string) {
    RespondWithError(w, http.StatusNotFound, message)
}

// InternalError sends an internal server error response
func InternalError(w http.ResponseWriter, message string) {
    RespondWithError(w, http.StatusInternalServerError, message)
}
