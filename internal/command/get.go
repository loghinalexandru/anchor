package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/peterbourgon/ff/v4"
)

const (
	formatShort = "%d. %s\n"
	formatLong  = "%d. %s %s\n"
)

var (
	ErrInvalidLine = errors.New("invalid line specified")
)

type getCmd struct {
	command ff.Command
	labels  []string
	pattern string
	full    bool
}

func newGet(rootFlags *ff.FlagSet) *getCmd {
	cmd := getCmd{}

	flags := ff.NewFlagSet("get").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "specify label hierarchy")
	_ = flags.StringVar(&cmd.pattern, 'p', "pattern", "", "match for pattern in bookmark title")
	_ = flags.BoolVar(&cmd.full, 'f', "full", "show full bookmark entry")

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
	match := FindLines(content, get.pattern)

	// Redesign this
	if len(args) >= 1 {
		line, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		if int(line) > len(match) || line == 0 {
			return ErrInvalidLine
		}

		bmk, err := bookmark.NewFromLine(string(match[line-1]))
		if err != nil {
			return err
		}

		err = Open(bmk.URL)
		if err != nil {
			return err
		}

		return nil
	}

	for i, l := range match {
		bmk, err := bookmark.NewFromLine(string(l))
		if err != nil {
			return err
		}

		if get.full {
			_, _ = fmt.Fprintf(os.Stdout, formatLong, i+1, bmk.Title, bmk.URL)
		} else {
			_, _ = fmt.Fprintf(os.Stdout, formatShort, i+1, bmk.Title)
		}
	}

	return nil
}
