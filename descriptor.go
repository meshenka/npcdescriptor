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

// Descriptor returns a random descriptor in English (default).
func Descriptor(ctx context.Context) string {
	return DescriptorWithLocale(ctx, "en")
}

// DescriptorWithLocale returns a random descriptor for the given locale.
// Falls back to "en" if locale is not found.
func DescriptorWithLocale(ctx context.Context, locale string) string {
	descriptors, ok := locales[locale]
	if !ok {
		descriptors = locales["en"]
	}
	return choose(descriptors)
}

func choose(items []string) string {
	if len(items) == 0 {
		return ""
	}
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		panic(err)
	}
	return items[i.Int64()]
}
