package command

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/peterbourgon/ff/v4"
)

const (
	clientTimeout = 5 * time.Second
)

var (
	ErrInvalidArgument = errors.New("missing bookmark URL from arguments")
)

type addCmd struct {
	labels []string
	title  string
}

func (add *addCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("add").SetParent(parent)
	flags.StringSetVar(&add.labels, 'l', "label", "add labels in order of appearance")
	flags.StringVar(&add.title, 't', "title", "", "add custom title")

	return &ff.Command{
		Name:      "add",
		Usage:     "add",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec:      add.handle,
	}
}

func (add *addCmd) handle(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return ErrInvalidArgument
	}

	client := &http.Client{Timeout: clientTimeout}
	b, err := bookmark.New(add.title, args[0], bookmark.WithClient(client))

	if err != nil {
		return err
	}

	if add.title == "" {
		err = b.TitleFromURL(ctx)
		if err != nil {
			return err
		}
	}

	err = label.Validate(add.labels)
	if err != nil {
		return err
	}

	path := label.Filepath(add.labels)
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
