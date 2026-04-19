package api

import (
	"net/http"

	"github.com/meshenka/npcgenerator/internal/app"
	"github.com/meshenka/npcgenerator/internal/logging"
)

// RegisterRoutes registers the API routes to the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, app app.Application) {
	h := NewDescriptorsHandler(app)
	mux.HandleFunc("GET /api/descriptors", h.GetDescriptors)
}

// MuxMiddleware wraps the entire ServeMux to add base logging and request ID.
func MuxMiddleware(mux *http.ServeMux) http.Handler {
	return logging.LoggingMiddleware(mux)
}
