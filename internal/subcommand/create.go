package subcommand

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
)

var (
	ErrInvalidURL   = errors.New("not a valid url")
	ErrDuplicate    = errors.New("duplicate")
	ErrInvalidTitle = errors.New("could not infer title and no flag was set")
)

type createCmd ff.Command

func RegisterCreate(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *createCmd
	flags := ff.NewFlags("create").SetParent(rootFlags)
	_ = flags.String('l', "label", "root", "add label in order of appearance split by ','")
	_ = flags.String('t', "title", "", "add custom title")

	cmd = &createCmd{
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

	root.Subcommands = append(root.Subcommands, (*ff.Command)(cmd))
}

func (cmd *createCmd) handle(ctx context.Context, args []string, res chan<- error) {
	defer close(res)

	labelFlag, _ := cmd.Flags.GetFlag("label")
	titleFlag, _ := cmd.Flags.GetFlag("title")
	dir, _ := cmd.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()

	if err != nil {
		res <- err
		return
	}

	b, err := bookmark.New(titleFlag.GetValue(), args[0])
	if err != nil {
		res <- err
		return
	}

	if titleFlag.GetValue() == titleFlag.GetDefault() {
		err = b.TitleFromURL(ctx)

		if err != nil {
			res <- err
			return
		}
	}

	hierarchy := strings.Split(labelFlag.GetValue(), ",")
	path := filepath.Join(home, dir.GetValue(), strings.Join(hierarchy, "."))
	_, err = bookmark.Append(*b, path)

	if err != nil {
		res <- err
	}
}
