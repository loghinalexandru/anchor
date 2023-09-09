package subcommand

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/peterbourgon/ff/v4"
)

type initCmd ff.Command

var (
	ErrHomeDir   = errors.New("could not open home directory")
	ErrCreateDir = errors.New("could not create anchor directory")
)

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
	repo, _ := c.Flags.GetFlag("repository")
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- fmt.Errorf("%w with base: '%w'", ErrHomeDir, err)
		return
	}

	path := filepath.Join(home, dir.GetValue())

	if repo.GetValue() == "true" {
		repo, err := git.PlainInit(path, false)
		if err != nil {
			res <- err
			return
		}

		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name:  "origin",
			URLs:  []string{args[0]},
			Fetch: []config.RefSpec{config.RefSpec("+refs/heads/*:refs/remotes/origin/*")},
		})

		res <- err
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, fs.ModePerm)

		if err != nil {
			res <- fmt.Errorf("%w with base: '%w'", ErrCreateDir, err)
			return
		}
	}

	close(res)
}
