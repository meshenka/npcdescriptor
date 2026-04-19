package query

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/json"
	"errors"
	"math/big"
)

//go:embed data/data.en.json
var npcEn []byte

//go:embed data/data.fr.json
var npcFr []byte

var locales = map[string][]string{}

func init() {
	var en []string
	if err := json.Unmarshal(npcEn, &en); err != nil {
		panic(err)
	}
	locales["en"] = en

	var fr []string
	if err := json.Unmarshal(npcFr, &fr); err != nil {
		panic(err)
	}
	locales["fr"] = fr
}

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

	descriptors := descriptorsWithLocale(query.Lang, query.N)
	if len(descriptors) == 0 {
		return nil, errors.New("no descriptors available")
	}

	return descriptors, nil
}

func descriptorsWithLocale(locale string, n int) []string {
	descriptors, ok := locales[locale]
	if !ok {
		descriptors = locales["en"]
	}

	if n <= 0 {
		return nil
	}

	count := len(descriptors)
	if n > count {
		n = count
	}

	indices := make([]int, count)
	for i := range indices {
		indices[i] = i
	}

	// Shuffle indices
	for i := count - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			// If random fails, we stop shuffling and return what we have
			break
		}
		j := int(jBig.Int64())
		indices[i], indices[j] = indices[j], indices[i]
	}

	res := make([]string, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, descriptors[indices[i]])
	}
	return res
}
