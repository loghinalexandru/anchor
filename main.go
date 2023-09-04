package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/loghinalexandru/anchor/internal/subcommand"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	rootFlags := ff.NewFlags("anchor")
	_ = rootFlags.Bool('v', "verbose", false, "increase log verbosity")
	rootCmd := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	err := subcommand.RegisterInit(rootCmd, rootFlags)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	err = rootCmd.ParseAndRun(ctx, os.Args[1:])

	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		fmt.Fprint(os.Stdout, ffhelp.Command(rootCmd))
		os.Exit(0)
	} else if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
