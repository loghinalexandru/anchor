package command

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea"
	"github.com/peterbourgon/ff/v4"
)

type getCmd struct {
	command ff.Command
	labels  []string
}

func newGet(rootFlags *ff.FlagSet) *getCmd {
	cmd := getCmd{}

	flags := ff.NewFlagSet("get").SetParent(rootFlags)
	_ = flags.StringSetVar(&cmd.labels, 'l', "label", "specify label hierarchy")

	cmd.command = ff.Command{
		Name:      "get",
		Usage:     "get",
		ShortHelp: "get existing bookmarks",
		Flags:     flags,
		Exec:      cmd.handle,
	}

	return &cmd
}

func (get *getCmd) handle(_ context.Context, _ []string) error {
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
	match := FindLines(content, "")
	bookmarks := make([]list.Item, len(match))

	for i, l := range match {
		bookmarks[i], err = model.NewFromLine(string(l))
		if err != nil {
			return err
		}
	}

	p := tea.NewProgram(bubbletea.NewView(bookmarks))
	_, err = p.Run()

	return err
}
