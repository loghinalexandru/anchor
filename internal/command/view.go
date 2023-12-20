package command

import (
	"bufio"
	"context"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/loghinalexandru/anchor/internal/command/util/label"
	"github.com/loghinalexandru/anchor/internal/config"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/loghinalexandru/anchor/internal/output"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea/style"
	"github.com/peterbourgon/ff/v4"
)

const (
	viewName        = "view"
	msgApplyChanges = "You are about to apply changes from previous operation. Proceed?"
)

type viewCmd struct {
	labels []string
}

func (v *viewCmd) manifest(parent *ff.FlagSet) *ff.Command {
	flags := ff.NewFlagSet("view").SetParent(parent)
	flags.StringSetVar(&v.labels, 'l', "label", "specify label hierarchy")

	return &ff.Command{
		Name:      viewName,
		Usage:     "view",
		ShortHelp: "view existing bookmarks",
		Flags:     flags,
		Exec:      v.handle,
	}
}

func (v *viewCmd) handle(ctx context.Context, _ []string) error {
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
		bk, err := model.BookmarkLine(scanner.Text())
		if err != nil {
			return err
		}

		bookmarks = append(bookmarks, bk)
	}

	runner := tea.NewProgram(bubbletea.NewView(bookmarks), tea.WithContext(ctx))
	state, err := runner.Run()
	if err != nil {
		return err
	}

	view := state.(*bubbletea.View)
	confirmer := output.Confirmer{
		MaxRetries: 3,
		Renderer:   style.Prompt,
	}

	if view.Dirty() && !confirmer.Confirm(msgApplyChanges, os.Stdin, os.Stdout) {
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
