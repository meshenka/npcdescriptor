package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshenka/npcgenerator/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGetDescriptorsHandler(t *testing.T) {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	t.Run("default (en)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/descriptors", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp api.DescriptorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Descriptors, 3)
		for _, desc := range resp.Descriptors {
			assert.NotEmpty(t, desc)
		}
	})

	t.Run("french (fr)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/descriptors?lang=fr", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp api.DescriptorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Descriptors, 3)
		for _, desc := range resp.Descriptors {
			assert.NotEmpty(t, desc)
		}
	})

	t.Run("uniqueness", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/descriptors?n=10", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp api.DescriptorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Descriptors, 10)

		seen := make(map[string]bool)
		for _, desc := range resp.Descriptors {
			assert.False(t, seen[desc], "Duplicate descriptor found in API response: %s", desc)
			seen[desc] = true
		}
	})
}
