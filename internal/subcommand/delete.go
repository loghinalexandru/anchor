package subcommand

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/regex"
	"github.com/peterbourgon/ff/v4"
)

var (
	ErrInvalidPattern = errors.New("invalid pattern")
)

type deleteCmd ff.Command

func RegisterDelete(root *ff.Command, rootFlags *ff.CoreFlags) {
	var cmd *deleteCmd
	flags := ff.NewFlags("delete").SetParent(rootFlags)
	_ = flags.StringSet('l', "label", "add label in order of appearance")

	cmd = &deleteCmd{
		Name:      "delete",
		Usage:     "delete",
		ShortHelp: "remove a bookmark",
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

func (c *deleteCmd) handle(args []string, res chan<- error) {
	defer close(res)

	labelFlag, _ := c.Flags.GetFlag("label")
	dir, _ := c.Flags.GetFlag("root-dir")
	home, err := os.UserHomeDir()
	if err != nil {
		res <- err
		return
	}

	tree, err := flatten(labelFlag.GetValue())
	if err != nil {
		res <- err
		return
	}

	path := filepath.Join(home, dir.GetValue(), tree)
	fh, err := os.OpenFile(path, os.O_RDWR, fs.ModePerm)
	if err != nil {
		res <- err
		return
	}

	defer fh.Close()

	if len(args) == 0 {
		res <- ErrInvalidPattern
		return
	}

	content, _ := io.ReadAll(fh)
	ll := regex.FindLines(content, args[0])
	for _, l := range ll {
		l = append(l, byte('\n'))
		content = bytes.ReplaceAll(content, l, []byte(""))
	}

	// Refactor this to be more efficient
	_ = fh.Truncate(0)
	_, _ = fh.Seek(0, 0)
	_, err = fh.Write(content)

	if err != nil {
		res <- err
		return
	}
}
