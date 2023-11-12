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

type syncCmd struct {
	command ff.Command
	storer  *storage.Storer
}

func newSync(rootFlags *ff.FlagSet, storer *storage.Storer) *syncCmd {
	var cmd syncCmd

	flags := ff.NewFlagSet("sync").SetParent(rootFlags)
	cmd.command = ff.Command{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     flags,
		Exec:      cmd.handle,
	}
	cmd.storer = storer

	return &cmd
}

func (sync *syncCmd) handle(context.Context, []string) error {
	if d, ok := (*sync.storer).(Differ); ok {
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

	return (*sync.storer).Store()
}
