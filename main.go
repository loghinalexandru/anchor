package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/loghinalexandru/anchor/internal/subcommand"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	rootFlags := ff.NewFlagSet("anchor")
	rootCmd := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	subcommand.RegisterInit(rootCmd, rootFlags)
	subcommand.RegisterCreate(rootCmd, rootFlags)
	subcommand.RegisterGet(rootCmd, rootFlags)
	subcommand.RegisterDelete(rootCmd, rootFlags)
	subcommand.RegisterSync(rootCmd, rootFlags)
	subcommand.RegisterImport(rootCmd, rootFlags)
	subcommand.RegisterTree(rootCmd, rootFlags)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err := rootCmd.ParseAndRun(ctx, os.Args[1:])

	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		_, _ = fmt.Fprint(os.Stdout, ffhelp.Command(rootCmd))
		os.Exit(0)
	} else if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
