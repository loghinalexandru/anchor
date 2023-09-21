package command

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type syncCmd ff.Command

func NewSync(rootFlags *ff.FlagSet) *syncCmd {
	var cmd *syncCmd
	flags := ff.NewFlagSet("sync").SetParent(rootFlags)

	cmd = &syncCmd{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return cmd
}

func (*syncCmd) handle(context.Context, []string) error {
	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	err = storage.PushWithSSH(dir)
	if err != nil {
		return err
	}

	return nil
}
