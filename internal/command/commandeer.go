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
	root := newRoot()
	err := root.bootstrap(args)
	if err != nil {
		return err
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err = root.cmd.Run(ctx)
	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		_, _ = fmt.Fprint(os.Stdout, ffhelp.Command(&root.cmd))
		return nil
	}

	return err
}
