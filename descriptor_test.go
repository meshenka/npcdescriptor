package npcgenerator_test

import (
	"testing"

	"github.com/meshenka/npcgenerator"
	"github.com/stretchr/testify/assert"
)

func TestDescriptor(t *testing.T) {
	ctx := t.Context()

	desc1, desc2 := npcgenerator.Descriptor(ctx), npcgenerator.Descriptor(ctx)
	assert.NotZero(t, desc1)
	assert.NotZero(t, desc2)
}
