package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondWithJSON is a helper function to send JSON responses
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	fmt.Println("JSON response", code, string(response))
}

// RespondWithError is a helper function to send error responses
func RespondWithError(w http.ResponseWriter, code int, message string) {
	fmt.Println("RESPONDED with ",message)
	RespondWithJSON(w, code, map[string]string{"error": message})
}
