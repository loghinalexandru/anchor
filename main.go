package main

import (
	"fmt"
	"github.com/loghinalexandru/anchor/internal/types"
	"os"
)

func main() {
	err := types.Execute(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
