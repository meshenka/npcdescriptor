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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Descriptor(context.Context) (string, error) {
	var npcDesc []string
	if err := json.Unmarshal(npc, &npcDesc); err != nil {
		return "", err
	}

	return choose(npcDesc), nil
}

func choose(items []string) string {
	return items[rand.Intn(len(items))]
}
