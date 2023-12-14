package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type initCmd struct{}

func (init *initCmd) manifest(parent *ff.FlagSet) *ff.Command {
	return &ff.Command{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "init a new empty home for anchor",
		Flags:     ff.NewFlagSet("init").SetParent(parent),
		Exec: func(ctx context.Context, args []string) error {
			storer := ctx.Value(storerContextKey{}).(storage.Storer)
			return storer.Init(args...)
		},
	}
}
