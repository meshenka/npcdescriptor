package npcgenerator_test

import (
	"context"
	"testing"

	"github.com/meshenka/npcgenerator"
	"github.com/stretchr/testify/assert"
)

func TestDescriptor(t *testing.T) {
	ctx := context.Background()

	desc1, desc2 := npcgenerator.Descriptor(ctx), npcgenerator.Descriptor(ctx)
	assert.NotZero(t, desc1)
	assert.NotZero(t, desc2)
	assert.NotEqual(t, desc1, desc2)
}
