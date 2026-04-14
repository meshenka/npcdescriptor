package api

import (
	"encoding/json"
	"net/http"

	"github.com/meshenka/npcgenerator"
)

// DescriptorResponse represents the JSON response for the descriptor endpoint.
type DescriptorResponse struct {
	Descriptors []string `json:"descriptors"`
}

// GetDescriptorsHandler returns three random NPC descriptors.
// @Summary Get NPC descriptors
// @Description returns three random NPC descriptors as a JSON list
// @Produce  json
// @Success 200 {object} DescriptorResponse
// @Router /api/descriptors [get]
func GetDescriptorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := DescriptorResponse{
		Descriptors: []string{
			npcgenerator.Descriptor(ctx),
			npcgenerator.Descriptor(ctx),
			npcgenerator.Descriptor(ctx),
		},
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
