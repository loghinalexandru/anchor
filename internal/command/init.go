package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type initCmd struct {
	command  ff.Command
	repoFlag bool
	storer   storage.Storer
}

func newInit(rootFlags *ff.FlagSet) *initCmd {
	var cmd initCmd

	flags := ff.NewFlagSet("init").SetParent(rootFlags)
	_ = flags.BoolVar(&cmd.repoFlag, 'r', "repository", "used in order to init the storage")

	cmd.command = ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     flags,
		Exec:      cmd.handle,
	}

	return &cmd
}

func (init *initCmd) withStorage(storer storage.Storer) {
	init.storer = storer
}

func (init *initCmd) handle(_ context.Context, args []string) error {
	return init.storer.Init(args[0])
}
