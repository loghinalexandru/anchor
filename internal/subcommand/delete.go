package subcommand

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteConfirmation = "You are about to deleteCmd %s. Proceed?"
)

type deleteCmd struct {
	command ff.Command
	labels  []string
}

func RegisterDelete(root *ff.Command, rootFlags *ff.FlagSet) {
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

	root.Subcommands = append(root.Subcommands, &cmd.command)
}

func (del *deleteCmd) handle(_ context.Context, args []string) (err error) {

	dir, err := rootDir()
	if err != nil {
		return err
	}

	err = validate(del.labels)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, fileFrom(del.labels))
	if len(args) == 0 {
		ok := confirmation(fmt.Sprintf(msgDeleteConfirmation, path), os.Stdin)
		if ok {
			err = os.Remove(path)
			return err
		}
		return nil
	}

	fh, err := os.OpenFile(path, os.O_RDWR, stdFileMode)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, fh.Close())
	}()

	content, _ := io.ReadAll(fh)
	ll := findLines(content, args[0])
	ok := confirmation(fmt.Sprintf(msgDeleteConfirmation, fmt.Sprintf("%del line(s)", len(ll))), os.Stdin)

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
