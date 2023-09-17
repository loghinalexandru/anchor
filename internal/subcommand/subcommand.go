package subcommand

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func Execute(args []string) error {
	rootFlags := ff.NewFlagSet("anchor")
	rootCmd := &ff.Command{
		Name:  "anchor",
		Usage: "anchor [FLAGS] <SUBCOMMAND>",
		Flags: rootFlags,
	}

	RegisterInit(rootCmd, rootFlags)
	RegisterCreate(rootCmd, rootFlags)
	RegisterGet(rootCmd, rootFlags)
	RegisterDelete(rootCmd, rootFlags)
	RegisterSync(rootCmd, rootFlags)
	RegisterImport(rootCmd, rootFlags)
	RegisterTree(rootCmd, rootFlags)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err := rootCmd.ParseAndRun(ctx, args)

	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		_, _ = fmt.Fprint(os.Stdout, ffhelp.Command(rootCmd))
		return nil
	}

	return err
}
