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
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(res)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-res:
				return err
			}
		},
	}

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (c *syncCmd) handle(res chan<- error) {
	defer close(res)

	rootDir, err := rootDir()
	if err != nil {
		res <- err
		return
	}

	err = storage.PushWithSSH(rootDir)
	if err != nil {
		res <- err
		return
	}
}
