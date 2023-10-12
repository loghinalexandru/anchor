package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/model"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteConfirmation = "You are about to delete %s. Proceed?"
	msgDeleteEmpty        = "No match found. Skipping..."
)

type deleteCmd struct {
	command ff.Command
	labels  []string
	pattern string
}

func newDelete(rootFlags *ff.FlagSet) *deleteCmd {
	cmd := deleteCmd{}

	flags := ff.NewFlagSet("delete").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add label in order of appearance")
	_ = flags.StringVar(&cmd.pattern, 'p', "pattern", "", "delete items matching the pattern")

	cmd.command = ff.Command{
		Name:      "delete",
		Usage:     "delete",
		ShortHelp: "remove a bookmark",
		Flags:     flags,
		Exec:      cmd.handle,
	}

	return &cmd
}

func (del *deleteCmd) handle(_ context.Context, _ []string) (err error) {
	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	err = Validate(del.labels)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, FileFrom(del.labels))
	if del.pattern == "" {
		return deleteFile(path)
	}

	fh, err := os.OpenFile(path, os.O_RDWR, config.StdFileMode)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, fh.Close())
	}()

	newContent, err := deleteContent(fh, del.pattern)
	if err != nil {
		return err
	}

	err = fh.Truncate(0)
	if err != nil {
		return err
	}

	_, err = fh.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = fh.Write(newContent)

	return err
}

func deleteFile(path string) error {
	ok := output.Confirmation(fmt.Sprintf(msgDeleteConfirmation, path), os.Stdin, os.Stdout)
	if ok {
		err := os.Remove(path)
		return err
	}

	return nil
}

func deleteContent(reader io.Reader, pattern string) ([]byte, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	ll := FindLines(content, pattern)
	if len(ll) == 0 {
		fmt.Println(msgDeleteEmpty)
		return content, nil
	}

	for _, l := range ll {
		bmk, err := model.NewFromLine(string(l))
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stdout, "%q\n", bmk.Name)
	}

	ok := output.Confirmation(fmt.Sprintf(msgDeleteConfirmation, fmt.Sprintf("%d line(s)", len(ll))), os.Stdin, os.Stdout)
	if !ok {
		return content, nil
	}

	return DeleteLines(content, pattern), nil
}
