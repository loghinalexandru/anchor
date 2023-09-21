package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	config "github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteConfirmation = "You are about to delete %s. Proceed?"
)

type deleteCmd struct {
	command ff.Command
	labels  []string
}

func newDelete(rootFlags *ff.FlagSet) *deleteCmd {
	cmd := deleteCmd{}

	flags := ff.NewFlagSet("delete").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add label in order of appearance")

	cmd.command = ff.Command{
		Name:      "delete",
		Usage:     "delete",
		ShortHelp: "remove a bookmark",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return &cmd
}

func (del *deleteCmd) handle(_ context.Context, args []string) (err error) {

	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	err = Validate(del.labels)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, FileFrom(del.labels))
	if len(args) == 0 {
		ok := output.Confirmation(fmt.Sprintf(msgDeleteConfirmation, path), os.Stdin)
		if ok {
			err = os.Remove(path)
			return err
		}
		return nil
	}

	fh, err := os.OpenFile(path, os.O_RDWR, config.StdFileMode)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, fh.Close())
	}()

	content, _ := io.ReadAll(fh)
	ll := FindLines(content, args[0])
	ok := output.Confirmation(fmt.Sprintf(msgDeleteConfirmation, fmt.Sprintf("%d line(s)", len(ll))), os.Stdin)

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
		return err
	}

	return nil
}
