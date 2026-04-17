package npcgenerator

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/json"
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

// DescriptorsWithLocale returns n unique random descriptors for the given locale.
// Falls back to "en" if locale is not found.
// If n is greater than the total number of available descriptors, all descriptors are returned.
func DescriptorsWithLocale(ctx context.Context, locale string, n int) []string {
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
