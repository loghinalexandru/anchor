package subcommand

import (
	"context"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

type syncCmd ff.Command

func RegisterSync(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *syncCmd
	flags := ff.NewFlags("sync").SetParent(rootFlags)

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

	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	path := filepath.Join(home, dir.GetValue())
	err = storage.PushWithSSH(path)

	if err != nil {
		res <- err
		return
	}
}
