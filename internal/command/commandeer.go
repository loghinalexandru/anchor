package command

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
	rootCmd := NewExec()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err := rootCmd.ParseAndRun(ctx, args)

	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		_, _ = fmt.Fprint(os.Stdout, ffhelp.Command(rootCmd))
		return nil
	}

	return err
}
