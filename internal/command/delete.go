package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/loghinalexandru/anchor/internal/command/util/text"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteLabel     = "You are about to delete the label and associated bookmarks. Proceed?"
	msgDeleteBookmarks = "You are about to delete matched bookmark(s). Proceed?"
	msgDeleteEmpty     = "No match found. Skipping..."
)

type deleteCmd struct {
	command ff.Command
	labels  []string
}

func newDelete(rootFlags *ff.FlagSet) *deleteCmd {
	var cmd deleteCmd

	flags := ff.NewFlagSet("delete").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add label in order of appearance")

	cmd.command = ff.Command{
		Name:      "delete",
		Usage:     "delete",
		ShortHelp: "remove a bookmark",
		Flags:     flags,
		Exec:      cmd.handle,
	}

	return &cmd
}

func (del *deleteCmd) handle(_ context.Context, args []string) (err error) {
	path := label.Filepath(del.labels)

	if len(args) == 0 {
		return deleteFile(path)
	}

	err = label.Validate(del.labels)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(path, os.O_RDWR, config.StdFileMode)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, fh.Close())
	}()

	newContent, err := deleteContent(fh, args[0])
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
	ok := output.Confirmation(msgDeleteLabel, os.Stdin, os.Stdout)
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

	ll := text.FindLines(content, pattern)
	if len(ll) == 0 {
		fmt.Println(msgDeleteEmpty)
		return content, nil
	}

	for _, l := range ll {
		bm, err := bookmark.NewFromLine(string(l))
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stdout, "%q\n", bm.Name)
	}

	ok := output.Confirmation(msgDeleteBookmarks, os.Stdin, os.Stdout)
	if !ok {
		return content, nil
	}

	return text.DeleteLines(content, pattern), nil
}
