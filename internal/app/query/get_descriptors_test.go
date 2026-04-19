package query_test

import (
	"context"
	"testing"

	"github.com/meshenka/npcgenerator/internal/app/query"
	"github.com/stretchr/testify/assert"
)

func TestGetDescriptorsHandler_Handle(t *testing.T) {
	h := query.NewGetDescriptorsHandler()
	ctx := context.Background()

	t.Run("success en", func(t *testing.T) {
		q := query.GetDescriptors{Lang: "en", N: 3}
		res, err := h.Handle(ctx, q)
		assert.NoError(t, err)
		assert.Len(t, res, 3)

		seen := make(map[string]bool)
		for _, d := range res {
			assert.NotEmpty(t, d)
			assert.False(t, seen[d], "Found duplicate descriptor: %s", d)
			seen[d] = true
		}
	})

	t.Run("success fr", func(t *testing.T) {
		q := query.GetDescriptors{Lang: "fr", N: 1}
		res, err := h.Handle(ctx, q)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("invalid N too high", func(t *testing.T) {
		q := query.GetDescriptors{Lang: "en", N: 11}
		res, err := h.Handle(ctx, q)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "n must be between 1 and 10", err.Error())
	})

	t.Run("invalid N too low", func(t *testing.T) {
		q := query.GetDescriptors{Lang: "en", N: 0}
		res, err := h.Handle(ctx, q)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "n must be between 1 and 10", err.Error())
	})

	t.Run("unsupported language falls back to en", func(t *testing.T) {
		q := query.GetDescriptors{Lang: "de", N: 3}
		res, err := h.Handle(ctx, q)
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})

	t.Run("empty language defaults to en", func(t *testing.T) {
		q := query.GetDescriptors{Lang: "", N: 3}
		res, err := h.Handle(ctx, q)
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})
}
