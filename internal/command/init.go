package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type initCmd struct {
	storer storage.Storer
}

func newInit() *initCmd {
	return &initCmd{
		storer: storage.New(storage.Local),
	}
}

func (init *initCmd) withStorage(storer storage.Storer) {
	init.storer = storer
}

func (init *initCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     ff.NewFlagSet("init").SetParent(parent),
		Exec:      init.handle,
	}
}

func (init *initCmd) handle(_ context.Context, args []string) error {
	return init.storer.Init(args...)
}
