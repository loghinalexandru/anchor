package subcommand

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/storage"
	"github.com/peterbourgon/ff/v4"
)

var (
	ErrInvalidURL = errors.New("not a valid URL")
)

type initCmd ff.Command

func RegisterInit(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *initCmd
	flags := ff.NewFlags("init").SetParent(rootFlags)
	_ = flags.Bool('r', "repository", false, "used in order to init a git repository")

	cmd = &initCmd{
		Name:      "init",
		Usage:     "init",
		ShortHelp: "initialize a new empty home for anchor",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(args, res)

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

func (c *initCmd) handle(args []string, res chan<- error) {
	defer close(res)

	repo, _ := c.Flags.GetFlag("repository")
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	path := filepath.Join(home, dir.GetValue())
	if repo.GetValue() == "true" {
		if len(args) == 0 {
			res <- ErrInvalidURL
			return
		}

		err := storage.CloneWithSSH(path, args[0])
		if err != nil {
			res <- err
		}

		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, fs.ModePerm)

		if err != nil {
			res <- err
			return
		}
	}
}
