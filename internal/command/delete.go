package command

import (
	"context"
	"os"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteLabel = "You are about to delete the label and associated bookmarks. Proceed?"
)

type deleteCmd struct {
	labels []string
}

func (del *deleteCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("delete").SetParent(parent)
	flags.StringSetVar(&del.labels, 'l', "label", "add label in order of appearance")

	return &ff.Command{
		Name:      "delete",
		Usage:     "delete",
		ShortHelp: "remove a bookmark",
		Flags:     flags,
		Exec:      del.handle,
	}
}

func (del *deleteCmd) handle(_ context.Context, _ []string) (err error) {
	path := label.Filepath(del.labels)

	ok := output.Confirmation(msgDeleteLabel, os.Stdin, os.Stdout)
	if !ok {
		return nil
	}

	err = os.Remove(path)
	if !os.IsNotExist(err) {
		return err
	}

	return nil
}
