package main

import (
	"context"
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
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	desc1, err := npcgenerator.Descriptor(ctx)
	if err != nil {
		return err
	}
	desc2, err := npcgenerator.Descriptor(ctx)
	if err != nil {
		return err
	}
	desc3, err := npcgenerator.Descriptor(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("This character can be described as %s, %s and %s\n", desc1, desc2, desc3)
	return nil
}
