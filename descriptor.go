package npcgenerator

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/json"
	"math/big"
)

//go:embed data/character.json
var npc []byte
var npcDescriptors []string

func init() {
	if err := json.Unmarshal(npc, &npcDescriptors); err != nil {
		panic(err)
	}
}

func Descriptor(context.Context) string {
	return choose(npcDescriptors)
}

func choose(items []string) string {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		panic(err)
	}
	return items[i.Int64()]
}
