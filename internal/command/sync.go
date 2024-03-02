package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

const (
	syncName      = "sync"
	syncUsage     = "anchor sync [FLAGS]"
	syncShortHelp = "synchronize changes with configured backing storage"
	syncLongHelp  = `  In order to persist changes in case of a remote backing storage
  this command needs to be invoked. Otherwise this will persist only on the local file system.
  This should be performed only for the write part since for reading anchor always gets the
  latest changes from the configured storage.

  Has no effect if the backing storage is set to "local".
`
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
		Name:      syncName,
		Usage:     syncUsage,
		ShortHelp: syncShortHelp,
		LongHelp:  syncLongHelp,
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			return sync.handle(ctx.(appContext), args)
		},
	}
}

func (sync *syncCmd) handle(ctx appContext, _ []string) error {
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

	if ok := output.Confirm(msgSyncConfirmation); !ok {
		return nil
	}

	return ctx.storer.Store(sync.msg)
}
