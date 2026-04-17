package query

import (
	"context"
	"errors"

	"github.com/meshenka/npcgenerator"
)

type GetDescriptors struct {
	Lang string
	N    int
}

type GetDescriptorsHandler struct{}

func NewGetDescriptorsHandler() GetDescriptorsHandler {
	return GetDescriptorsHandler{}
}

func (h GetDescriptorsHandler) Handle(ctx context.Context, query GetDescriptors) ([]string, error) {
	if query.N < 1 || query.N > 10 {
		return nil, errors.New("n must be between 1 and 10")
	}

	descriptors := npcgenerator.DescriptorsWithLocale(ctx, query.Lang, query.N)
	if len(descriptors) == 0 {
		return nil, errors.New("no descriptors available")
	}

	return descriptors, nil
}
