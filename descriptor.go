package npcgenerator

import (
	"context"
	_ "embed"
	"encoding/json"
	"math/rand"
	"time"
)

//go:embed data/character.json
var npc []byte
var npcDescriptors []string

func init() {
	rand.Seed(time.Now().UnixNano())
	if err := json.Unmarshal(npc, &npcDescriptors); err != nil {
		panic(err)
	}
}

func Descriptor(context.Context) string {
	return choose(npcDescriptors)
}

func choose(items []string) string {
	return items[rand.Intn(len(items))]
}
