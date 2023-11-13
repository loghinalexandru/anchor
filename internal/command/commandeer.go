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
	root, err := newRoot(args)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = root.Run(ctx)
	if errors.Is(err, ff.ErrHelp) || errors.Is(err, ff.ErrNoExec) {
		_, _ = fmt.Fprint(os.Stdout, ffhelp.Command(root))
		return nil
	}

	return err
}
