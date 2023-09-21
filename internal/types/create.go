package types

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/command"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/peterbourgon/ff/v4"
)

type createCmd struct {
	Command ff.Command
	labels  []string
	title   string
}

func NewCreate(rootFlags *ff.FlagSet) *createCmd {
	cmd := createCmd{}

	flags := ff.NewFlagSet("create").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add labels in order of appearance")
	_ = flags.StringVar(&cmd.title, 't', "title", "", "add custom title")

	cmd.Command = ff.Command{
		Name:      "create",
		Usage:     "create",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return &cmd
}

func (crt *createCmd) handle(ctx context.Context, args []string) error {
	dir, err := command.RootDir()
	if err != nil {
		return err
	}

	b, err := bookmark.New(crt.title, args[0])
	if err != nil {
		return err
	}

	if crt.title == "" {
		err = b.TitleFromURL(ctx)

		if err != nil {
			return err
		}
	}

	err = command.Validate(crt.labels)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, command.FileFrom(crt.labels))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, command.StdFileMode)
	if err != nil {
		return err
	}

	err = b.Write(file)
	err = errors.Join(err, file.Close())

	if err != nil {
		return err
	}

	return nil
}
