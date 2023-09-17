package main

import (
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/subcommand"
)

func main() {
	err := subcommand.Execute(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
