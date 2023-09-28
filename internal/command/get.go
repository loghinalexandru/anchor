package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/peterbourgon/ff/v4"
)

type getCmd struct {
	command  ff.Command
	labels   []string
	fullFlag bool
	openFlag bool
}

func newGet(rootFlags *ff.FlagSet) *getCmd {
	cmd := getCmd{}

	flags := ff.NewFlagSet("get").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "specify label hierarchy for each")
	_ = flags.BoolVar(&cmd.fullFlag, 'f', "full", "show full bookmark entry")
	_ = flags.BoolVar(&cmd.openFlag, 'o', "open", "open specified link")

	cmd.command = ff.Command{
		Name:      "get",
		Usage:     "get",
		ShortHelp: "get existing bookmarks",
		Flags:     flags,
		Exec:      handlerMiddleware(cmd.handle),
	}

	return &cmd
}

func (get *getCmd) handle(_ context.Context, args []string) error {
	dir, err := config.RootDir()
	if err != nil {
		return err
	}

	err = Validate(get.labels)
	if err != nil {
		return err
	}

	path, err := os.Open(filepath.Join(dir, FileFrom(get.labels)))
	if err != nil {
		return err
	}

	content, err := io.ReadAll(path)
	if err != nil {
		return err
	}

	_ = path.Close()
	var pattern string
	if len(args) >= 1 {
		pattern = args[0]
	}

	for _, l := range FindLines(content, pattern) {
		b, err := bookmark.NewFromLine(string(l))
		if err != nil {
			return err
		}

		if get.openFlag {
			err = Open(b.URL)
			if err != nil {
				return err
			}

			return nil
		}

		if get.fullFlag {
			_, _ = fmt.Fprintln(os.Stdout, b.Title, b.URL)
		} else {
			_, _ = fmt.Fprintln(os.Stdout, b.Title)
		}
	}

	return nil
}
