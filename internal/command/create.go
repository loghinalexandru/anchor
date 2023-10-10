package command

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/peterbourgon/ff/v4"
)

type createCmd struct {
	command ff.Command
	labels  []string
	title   string
}

func newCreate(rootFlags *ff.FlagSet) *createCmd {
	cmd := createCmd{}

	flags := ff.NewFlagSet("create").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "add labels in order of appearance")
	_ = flags.StringVar(&cmd.title, 't', "title", "", "add custom title")

	cmd.command = ff.Command{
		Name:      "create",
		Usage:     "create",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return &cmd
}

func (crt *createCmd) handle(ctx context.Context, args []string) error {
	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	b, err := model.New(crt.title, args[0])
	if err != nil {
		return err
	}

	if crt.title == "" {
		err = b.TitleFromURL(ctx)

		if err != nil {
			return err
		}
	}

	err = Validate(crt.labels)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, FileFrom(crt.labels))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, config.StdFileMode)
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
