package npcgenerator_test

import (
	"testing"

	"github.com/meshenka/npcgenerator"
	"github.com/stretchr/testify/assert"
)

func TestDescriptorsWithLocale_Uniqueness(t *testing.T) {
	ctx := t.Context()
	n := 10
	descriptors := npcgenerator.DescriptorsWithLocale(ctx, "en", n)

	assert.Len(t, descriptors, n)

	seen := make(map[string]bool)
	for _, d := range descriptors {
		assert.NotEmpty(t, d)
		assert.False(t, seen[d], "Found duplicate descriptor: %s", d)
		seen[d] = true
	}
}

func TestDescriptorsWithLocale_Cap(t *testing.T) {
	ctx := t.Context()
	// More than available in data.en.json (which has around 100)
	n := 1000
	descriptors := npcgenerator.DescriptorsWithLocale(ctx, "en", n)

	assert.Less(t, len(descriptors), n)
	assert.NotEmpty(t, descriptors)
}
