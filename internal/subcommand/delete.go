package subcommand

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteConfirmation = "You are about to delete %s. Proceed?"
)

type delete struct {
	command ff.Command
	labels  []string
}

func RegisterDelete(root *ff.Command, rootFlags *ff.FlagSet) {
	cmd := delete{}

	flags := ff.NewFlagSet("delete").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add label in order of appearance")

	cmd.command = ff.Command{
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

	root.Subcommands = append(root.Subcommands, &cmd.command)
}

func (d *delete) handle(args []string, res chan<- error) {
	defer close(res)

	rootDir, err := rootDir()
	if err != nil {
		res <- err
		return
	}

	err = validate(d.labels)
	if err != nil {
		res <- err
		return
	}

	file := fileFrom(d.labels)
	path := filepath.Join(rootDir, file)

	if len(args) == 0 {
		ok := confirmation(fmt.Sprintf(msgDeleteConfirmation, path), os.Stdin)
		if ok {
			err = os.Remove(path)
			res <- err
		}
		return
	}

	fh, err := os.OpenFile(path, os.O_RDWR, stdFileMode)
	if err != nil {
		res <- err
		return
	}

	defer fh.Close()

	content, _ := io.ReadAll(fh)
	ll := findLines(content, args[0])
	ok := confirmation(fmt.Sprintf(msgDeleteConfirmation, fmt.Sprintf("%d line(s)", len(ll))), os.Stdin)

	if !ok {
		return
	}

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
