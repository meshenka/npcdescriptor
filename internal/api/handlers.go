package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meshenka/npcgenerator"
)

// DescriptorResponse represents the JSON response for the descriptor endpoint.
type DescriptorResponse struct {
	Descriptors []string `json:"descriptors"`
}

// GetDescriptorsHandler returns random NPC descriptors.
// @Summary Get NPC descriptors
// @Description returns random NPC descriptors as a JSON list. Optional query param 'n' (1-10) sets count (default 3). 'lang' sets locale (default 'en').
// @Produce  json
// @Param n query int false "Number of descriptors" default(3) minimum(1) maximum(10)
// @Param lang query string false "Locale" default("en")
// @Success 200 {object} DescriptorResponse
// @Router /api/descriptors [get]
func GetDescriptorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	n := 3
	if nStr := r.URL.Query().Get("n"); nStr != "" {
		if val, err := strconv.Atoi(nStr); err == nil {
			n = val
		}
	}

	if n < 1 {
		n = 1
	}
	if n > 10 {
		n = 10
	}

	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}

	descriptors := make([]string, n)
	for i := 0; i < n; i++ {
		descriptors[i] = npcgenerator.DescriptorWithLocale(ctx, lang)
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
