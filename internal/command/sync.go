package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type Differ interface {
	Diff() (string, error)
}

type Updater interface {
	Update() error
}

type Storer interface {
	Store() error
}

const (
	msgNothingToSync    = "Nothing to sync, there are no local changes."
	msgSyncConfirmation = "Sync changes with remote?"
)

type syncCmd struct {
	command ff.Command
	storer  Storer
}

func newSync(rootFlags *ff.FlagSet) *syncCmd {
	var cmd syncCmd

	flags := ff.NewFlagSet("sync").SetParent(rootFlags)
	cmd.command = ff.Command{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     flags,
		Exec:      cmd.handle,
	}
	cmd.storer, _ = storage.NewGitStorage()

	return &cmd
}

func (sync *syncCmd) handle(context.Context, []string) error {
	if u, ok := sync.storer.(Updater); ok {
		err := u.Update()
		if err != nil {
			return err
		}
	}

	if d, ok := sync.storer.(Differ); ok {
		status, err := d.Diff()
		if err != nil {
			return err
		}

		if status == "" {
			fmt.Println(msgNothingToSync)
			return nil
		}

		_, _ = fmt.Fprint(os.Stdout, status)
		if ok := output.Confirmation(msgSyncConfirmation, os.Stdin, os.Stdout); !ok {
			return nil
		}
	}

	err := sync.storer.Store()
	if err != nil {
		return err
	}

	return nil
}
