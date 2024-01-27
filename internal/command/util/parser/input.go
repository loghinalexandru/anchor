package parser

import (
	"fmt"
)

// First always returns the first argument in the provided list
// if it exists; otherwise it reads it from stdin.
// Enables chaining of commands via UNIX pipes.
func First(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	var result string
	_, err := fmt.Scanln(&result)
	if err != nil {
		return ""
	}

	return result
}
