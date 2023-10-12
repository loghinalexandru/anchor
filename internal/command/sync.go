package command

import (
	"context"
	"fmt"
	"os"

	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type Storer interface {
	Read() error
	Write() error
	Status() (string, error)
}

const (
	msgNothingToSync = "Nothing to sync, there are no local changes."
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
	err := sync.storer.Read()
	if err != nil {
		return err
	}

	status, err := sync.storer.Status()
	if err != nil {
		return err
	}

	if status == "" {
		fmt.Println(msgNothingToSync)
		return nil
	}

	_, _ = fmt.Fprintln(os.Stdout, status)
	if ok := output.Confirmation("Sync changes with remote?", os.Stdin, os.Stdout); !ok {
		return nil
	}

	err = sync.storer.Write()
	if err != nil {
		return err
	}

	return nil
}
