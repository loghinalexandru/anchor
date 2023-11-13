package command

import (
	"context"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/command/util/text"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea"
	"github.com/peterbourgon/ff/v4"
)

type getCmd struct {
	labels []string
}

func (get *getCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("get").SetParent(parent)
	_ = flags.StringSetVar(&get.labels, 'l', "label", "specify label hierarchy")

	return &ff.Command{
		Name:      "get",
		Usage:     "get",
		ShortHelp: "get existing bookmarks",
		Flags:     flags,
		Exec:      get.handle,
	}
}

func (get *getCmd) handle(_ context.Context, _ []string) error {
	err := label.Validate(get.labels)
	if err != nil {
		return err
	}

	path, err := os.Open(label.Filepath(get.labels))
	if err != nil {
		return err
	}

	content, err := io.ReadAll(path)
	if err != nil {
		return err
	}

	_ = path.Close()
	match := text.FindLines(content, "")
	bookmarks := make([]list.Item, len(match))

	for i, l := range match {
		bookmarks[i], err = bookmark.NewFromLine(string(l))
		if err != nil {
			return err
		}
	}

	p := tea.NewProgram(bubbletea.NewView(bookmarks))
	_, err = p.Run()

	return err
}
