package command

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/peterbourgon/ff/v4"
)

const (
	addName      = "add"
	addUsage     = "anchor add [FLAGS]"
	addShortHelp = "append a bookmark entry with set labels"
	addLongHelp  = `  Append a bookmark to a file on the backing storage determined by the 
  flatten hierarchy of the provided labels. Order of the flags matter when storing the entry.

  If no label is provided via the -l flag, all the entries will be added
  to the default "root" label.

  By default it tries to fetch the "title" content from the provided URL. If it fails
  to do so, it will store the entry with same title as the URL. You can provide a specific
  title with the flag -t and it overwrites the behaviour mentioned above.

EXAMPLES:
  # Append to default label
  anchor add "https://www.youtube.com/"

  # Append to a label "programming" with a sublabel "go"
  anchor add -l programming -l go "https://gobyexample.com/"
`
)

const (
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
		Usage:     addUsage,
		ShortHelp: addShortHelp,
		LongHelp:  addLongHelp,
		Flags:     flags,
		Exec:      add.handle,
	}
}

func (add *addCmd) handle(_ context.Context, args []string) error {
	if len(args) == 0 {
		return ErrMissingURL
	}

	file, err := label.Open(add.labels, os.O_APPEND|os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}

	b, err := model.NewBookmark(
		args[0],
		model.WithTitle(add.title),
		model.WithClient(&http.Client{Timeout: clientTimeout}))
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
