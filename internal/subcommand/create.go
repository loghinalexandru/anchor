package subcommand

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
)

type create struct {
	command ff.Command
	labels  []string
	title   string
}

func RegisterCreate(root *ff.Command, rootFlags *ff.FlagSet) {
	cmd := create{}

	flags := ff.NewFlagSet("create").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add labels in order of appearance")
	_ = flags.StringVar(&cmd.title, 't', "title", "", "add custom title")

	cmd.command = ff.Command{
		Name:      "create",
		Usage:     "crate",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec: func(ctx context.Context, args []string) error {
			res := make(chan error, 1)
			go cmd.handle(ctx, args, res)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-res:
				return err
			}
		},
	}

	root.Subcommands = append(root.Subcommands, &cmd.command)
}

func (c *create) handle(ctx context.Context, args []string, res chan<- error) {
	defer close(res)

	rootDir, err := rootDir()
	if err != nil {
		res <- err
		return
	}

	b, err := bookmark.New(c.title, args[0])
	if err != nil {
		res <- err
		return
	}

	if c.title == "" {
		err = b.TitleFromURL(ctx)

		if err != nil {
			res <- err
			return
		}
	}

	err = validate(c.labels)
	if err != nil {
		res <- err
		return
	}

	tree := fileFrom(c.labels)
	path := filepath.Join(rootDir, tree)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, stdFileMode)
	if err != nil {
		res <- err
		return
	}

	err = b.Write(file)
	err = errors.Join(err, file.Close())

	if err != nil {
		res <- err
	}
}
