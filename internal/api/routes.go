package api

import (
	"net/http"
)

// RegisterRoutes registers the API routes to the given ServeMux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/descriptors", GetDescriptorsHandler)
}
