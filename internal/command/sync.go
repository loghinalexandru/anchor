package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea/style"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgNothingToSync    = "Nothing to sync, there are no local changes."
	msgSyncConfirmation = "Sync changes with remote?"
)

type Differ interface {
	Diff() (string, error)
}

type syncCmd struct {
	msg string
}

func (sync *syncCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("sync").SetParent(parent)
	flags.StringVar(&sync.msg, 'm', "message", config.StdSyncMsg, "Optional sync message")

	return &ff.Command{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			return sync.handle(ctx.(rootContext), args)
		},
	}
}

func (sync *syncCmd) handle(ctx rootContext, _ []string) error {
	if d, ok := ctx.storer.(Differ); ok {
		status, err := d.Diff()
		if err != nil {
			return err
		}

		if status == "" {
			fmt.Println(msgNothingToSync)
			return nil
		}

		_, _ = fmt.Fprint(os.Stdout, status)
	}

	if ok := output.Confirmation(msgSyncConfirmation, os.Stdin, os.Stdout, style.Nop); !ok {
		return nil
	}

	return ctx.storer.Store(sync.msg)
}
