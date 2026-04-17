package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"strings"
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
	n := flag.Int("n", 3, "Number of descriptors to generate")
	flag.Parse()

	if *n < 1 || *n > 10 {
	        return fmt.Errorf("n must be between 1 and 10")
	}
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	descriptors := npcgenerator.DescriptorsWithLocale(ctx, *lang, *n)
	if len(descriptors) == 0 {
		return fmt.Errorf("no descriptors available")
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
