package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"strings"
	"syscall"

	"github.com/meshenka/npcgenerator/internal/app"
	"github.com/meshenka/npcgenerator/internal/app/query"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	lang := flag.String("lang", "en", "Language to use (en|fr)")
	n := flag.Int("n", 3, "Number of descriptors to generate")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	application := app.NewApplication()
	descriptors, err := application.Queries.GetDescriptors.Handle(ctx, query.GetDescriptors{
		Lang: *lang,
		N:    *n,
	})
	if err != nil {
		return err
	}

	var joined strings.Builder
	for i, d := range descriptors {
		if i > 0 {
			if i == len(descriptors)-1 {
				if *lang == "fr" {
					joined.WriteString(" et ")
				} else {
					joined.WriteString(" and ")
				}
			} else {
				joined.WriteString(", ")
			}
		}
		joined.WriteString(d)
	}

	switch *lang {
	case "fr":
		fmt.Printf("Ce personnage peut être décrit comme étant %s\n", joined.String())
	default:
		fmt.Printf("This character can be described as %s\n", joined.String())
	}
	return nil
}
