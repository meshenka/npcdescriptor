package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// renderError sends a JSON-formatted error response.
func renderError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}
