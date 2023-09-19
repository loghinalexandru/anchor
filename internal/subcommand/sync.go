package subcommand

import (
	"context"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type syncCmd ff.Command

func RegisterSync(root *ff.Command, rootFlags *ff.FlagSet) {
	var cmd *syncCmd
	flags := ff.NewFlagSet("sync").SetParent(rootFlags)

	cmd = &syncCmd{
		Name:      "sync",
		Usage:     "sync",
		ShortHelp: "sync changes with configured remote",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (*syncCmd) handle(context.Context, []string) error {
	dir, err := rootDir()
	if err != nil {
		return err
	}

	err = storage.PushWithSSH(dir)
	if err != nil {
		return err
	}

	return nil
}
