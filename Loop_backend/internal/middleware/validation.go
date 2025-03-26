package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"Loop_backend/internal/response"
)

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

var validate = validator.New()

// DecodeAndValidate is a middleware that decodes and validates the request body into a DTO
func DecodeAndValidate(handler http.HandlerFunc, dtoType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// if specified dto is nil, skip
		if dtoType == nil {
			handler(w, r)
			return
		}
		// Create a new instance of the DTO type
		dtoType := reflect.TypeOf(dtoType)
		if dtoType.Kind() == reflect.Ptr {
			dtoType = dtoType.Elem()
		}
		dto := reflect.New(dtoType).Interface()
		w.Header().Set("Content-Type", "application/json")

		// Read and restore the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading request body: %v\n", err)
			response.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid request payload", map[string]string{
				"error": "Failed to read request body",
			})
			return
		}
		r.Body.Close()
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		// Decode the JSON
		if err := json.NewDecoder(strings.NewReader(string(body))).Decode(&dto); err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			response.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid request payload", map[string]string{
				"decode_error": err.Error(),
			})
			return
		}

		// Validate the DTO
		if err := validate.Struct(dto); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			fmt.Printf("Validation errors: %+v\n", validationErrors)
			response.RespondWithErrorDetails(w, http.StatusBadRequest, "Validation failed",
				formatValidationErrors(err.(validator.ValidationErrors)))
			return
		}

		// Store the validated DTO in the request context
		var contextValue interface{}
		if reflect.TypeOf(dtoType).Kind() == reflect.Ptr {
			contextValue = dto
		} else {
			contextValue = reflect.ValueOf(dto).Elem().Interface()
		}
		ctx := context.WithValue(r.Context(), ValidatedDTOKey, contextValue)

		// Call the handler with the updated context
		handler(w, r.WithContext(ctx))
	}
}

// GetDTO retrieves the validated DTO from the request context
func GetDTO[T any](r *http.Request) (T, bool) {
	if dto, ok := r.Context().Value(ValidatedDTOKey).(T); ok {
		return dto, true
	}
	var zero T
	return zero, false
}

// Helper function to format validation errors
func formatValidationErrors(errors validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		field := strings.ToLower(err.Field())
		switch err.Tag() {
		case "required":
			errorMap[field] = "This field is required"
		case "email":
			errorMap[field] = "Invalid email format"
		case "min":
			errorMap[field] = fmt.Sprintf("Value must be at least %s characters", err.Param())
		case "max":
			errorMap[field] = fmt.Sprintf("Value must be at most %s characters", err.Param())
		case "oneof":
			errorMap[field] = fmt.Sprintf("Value must be one of: %s", err.Param())
		default:
			errorMap[field] = fmt.Sprintf("Failed validation: %s=%s", err.Tag(), err.Param())
		}
	}
	return errorMap
}

// ValidateRequest is a middleware that validates request body without authentication
func ValidateRequest(handler http.HandlerFunc, dtoType interface{}) http.HandlerFunc {
	return DecodeAndValidate(handler, dtoType)
}
