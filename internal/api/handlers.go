package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meshenka/npcgenerator/internal/app"
	"github.com/meshenka/npcgenerator/internal/app/query"
)

// DescriptorResponse represents the JSON response for the descriptor endpoint.
type DescriptorResponse struct {
	Descriptors []string `json:"descriptors"`
}

type DescriptorsHandler struct {
	app app.Application
}

func NewDescriptorsHandler(app app.Application) *DescriptorsHandler {
	return &DescriptorsHandler{app: app}
}

// GetDescriptors returns random NPC descriptors.
// @Summary Get NPC descriptors
// @Description returns random NPC descriptors as a JSON list. Optional query param 'n' (1-10) sets count (default 3). 'lang' sets locale (default 'en').
// @Produce  json
// @Param n query int false "Number of descriptors" default(3) minimum(1) maximum(10)
// @Param lang query string false "Locale" default("en")
// @Success 200 {object} DescriptorResponse
// @Router /api/descriptors [get]
func (h *DescriptorsHandler) GetDescriptors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	n := 3
	if nStr := r.URL.Query().Get("n"); nStr != "" {
		val, err := strconv.Atoi(nStr)
		if err != nil {
			http.Error(w, "invalid n parameter", http.StatusBadRequest)
			return
		}
		n = val
	}

	lang := r.URL.Query().Get("lang")

	descriptors, err := h.app.Queries.GetDescriptors.Handle(ctx, query.GetDescriptors{
		Lang: lang,
		N:    n,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := DescriptorResponse{
		Descriptors: descriptors,
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
