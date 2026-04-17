package api

import (
	"net/http"

	"github.com/meshenka/npcgenerator/internal/app"
)

// RegisterRoutes registers the API routes to the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, app app.Application) {
	h := NewDescriptorsHandler(app)
	mux.HandleFunc("GET /api/descriptors", h.GetDescriptors)
}
