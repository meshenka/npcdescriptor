package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/meshenka/npcgenerator"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	lang := flag.String("lang", "en", "Language to use (en|fr)")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	desc1 := npcgenerator.DescriptorWithLocale(ctx, *lang)
	desc2 := npcgenerator.DescriptorWithLocale(ctx, *lang)
	desc3 := npcgenerator.DescriptorWithLocale(ctx, *lang)

	switch *lang {
	case "fr":
		fmt.Printf("Ce personnage peut être décrit comme étant %s, %s et %s\n", desc1, desc2, desc3)
	default:
		fmt.Printf("This character can be described as %s, %s and %s\n", desc1, desc2, desc3)
	}
	return nil
}
