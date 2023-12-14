package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgNothingToSync    = "Nothing to sync, there are no local changes."
	msgSyncConfirmation = "Sync changes with remote?"
)

type Differ interface {
	Diff() (string, error)
}

type syncCmd struct{}

func (sync *syncCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     ff.NewFlagSet("sync").SetParent(parent),
		Exec:      sync.handle,
	}
}

func (sync *syncCmd) handle(ctx context.Context, _ []string) error {
	storer := ctx.Value(storerContextKey{}).(storage.Storer)

	if d, ok := storer.(Differ); ok {
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

	if ok := output.Confirmation(msgSyncConfirmation, os.Stdin, os.Stdout); !ok {
		return nil
	}

	return storer.Store()
}
