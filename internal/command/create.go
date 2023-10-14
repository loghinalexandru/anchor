package command

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/loghinalexandru/anchor/internal/command/util/label"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/peterbourgon/ff/v4"
)

const (
	clientTimeout = 5 * time.Second
)

var (
	ErrInvalidArgument = errors.New("missing bookmark argument")
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
		Exec:      cmd.handle,
	}

	return &cmd
}

func (crt *createCmd) handle(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return ErrInvalidArgument
	}

	client := &http.Client{Timeout: clientTimeout}
	b, err := bookmark.New(crt.title, args[0], bookmark.WithClient(client))

	if err != nil {
		return err
	}

	if crt.title == "" {
		err = b.TitleFromURL(ctx)

		if err != nil {
			return err
		}
	}

	err = label.Validate(crt.labels)
	if err != nil {
		return err
	}

	path := label.Filepath(crt.labels)
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
