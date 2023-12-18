package command

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/loghinalexandru/anchor/internal/bookmark"
	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea"
	"github.com/peterbourgon/ff/v4"
)

const (
	msgDeleteBookmarks = "You are about to delete %d bookmark(s) from previous operation. Proceed?"
)

type viewCmd struct {
	labels []string
}

func (v *viewCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("view").SetParent(parent)
	flags.StringSetVar(&v.labels, 'l', "label", "specify label hierarchy")

	return &ff.Command{
		Name:      "view",
		Usage:     "view",
		ShortHelp: "view existing bookmarks",
		Flags:     flags,
		Exec:      v.handle,
	}
}

func (v *viewCmd) handle(_ context.Context, _ []string) error {
	err := label.Validate(v.labels)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(label.Filepath(v.labels), os.O_RDWR, config.StdFileMode)
	if err != nil {
		return err
	}

	defer fh.Close()

	var bookmarks []list.Item

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		bk, err := bookmark.NewFromLine(scanner.Text())
		if err != nil {
			return err
		}

		bookmarks = append(bookmarks, bk)
	}

	runner := tea.NewProgram(bubbletea.NewView(bookmarks))
	model, err := runner.Run()
	if err != nil {
		return err
	}

	view := model.(*bubbletea.View)
	if len(view.Bookmarks()) < len(bookmarks) && !output.Confirmation(fmt.Sprintf(msgDeleteBookmarks, len(bookmarks)-len(view.Bookmarks())), os.Stdin, os.Stdout) {
		return nil
	}

	err = fh.Truncate(0)
	if err != nil {
		return err
	}

	_, err = fh.Seek(0, 0)
	if err != nil {
		return err
	}

	for _, b := range view.Bookmarks() {
		_, err := fh.WriteString(b.String())
		if err != nil {
			return err
		}
	}

	return nil
}
