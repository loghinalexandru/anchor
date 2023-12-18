package main

import (
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/command"
)

func main() {
	err := command.Execute(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
