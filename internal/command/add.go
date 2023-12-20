package command

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/peterbourgon/ff/v4"
)

const (
	addName       = "add"
	clientTimeout = 5 * time.Second
)

var (
	ErrMissingURL = errors.New("missing bookmark URL from arguments")
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
		Name:      addName,
		Usage:     "anchor add [FLAGS]",
		ShortHelp: "add a bookmark with set labels",
		Flags:     flags,
		Exec:      add.handle,
	}
}

func (add *addCmd) handle(_ context.Context, args []string) error {
	if len(args) == 0 {
		return ErrMissingURL
	}

	b, err := model.NewBookmark(
		args[0],
		model.WithTitle(add.title),
		model.WithClient(&http.Client{Timeout: clientTimeout}))

	if err != nil {
		return err
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
