package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type initCmd struct {
	command ff.Command
	storer  storage.Storer
}

func newInit(rootFlags *ff.FlagSet, storer storage.Storer) *initCmd {
	var cmd initCmd

	cmd.command = ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     ff.NewFlagSet("init").SetParent(rootFlags),
		Exec:      cmd.handle,
	}
	cmd.storer = storer

	return &cmd
}

func (init *initCmd) handle(_ context.Context, args []string) error {
	return init.storer.Init(args[0])
}
