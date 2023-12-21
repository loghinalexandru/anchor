package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

const (
	deleteName     = "delete"
	msgDeleteLabel = "You are about to delete the label and associated bookmarks. Proceed?"
)

type deleteCmd struct {
	labels []string
}

func (del *deleteCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("delete").SetParent(parent)
	flags.StringSetVar(&del.labels, 'l', "label", "add label in order of appearance")

	return &ff.Command{
		Name:      deleteName,
		Usage:     "anchor delete [FLAGS]",
		ShortHelp: "remove a bookmark",
		Flags:     flags,
		Exec:      del.handle,
	}
}

func (del *deleteCmd) handle(_ context.Context, _ []string) (err error) {
	ok := output.Confirm(msgDeleteLabel)
	if !ok {
		return nil
	}

	return label.Remove(del.labels)
}
