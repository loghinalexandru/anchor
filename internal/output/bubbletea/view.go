package bubbletea

import (
	"fmt"
	"os/exec"
	"runtime"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/loghinalexandru/anchor/internal/model"
	"github.com/loghinalexandru/anchor/internal/output/bubbletea/style"
)

const (
	msgStatus = "Deleted %q"
)

type View struct {
	input     textinput.Model
	bookmarks list.Model
	dirty     bool
}

func NewView(bookmarks []list.Item) *View {
	del := list.NewDefaultDelegate()
	style.ApplyToDelegate(&del)

	viewList := list.New(bookmarks, del, 0, 0)
	style.ApplyToList(&viewList)

	return &View{
		input:     textinput.New(),
		bookmarks: viewList,
	}
}

func (v *View) Bookmarks() []*model.Bookmark {
	res := make([]*model.Bookmark, len(v.bookmarks.Items()))

	for i, bk := range v.bookmarks.Items() {
		res[i] = bk.(*model.Bookmark)
	}

	return res
}

func (v *View) Dirty() bool {
	return v.dirty
}

func (v *View) Init() tea.Cmd {
	return nil
}

func (v *View) View() string {
	if v.input.Focused() {
		v.bookmarks.SetShowPagination(false)
		v.bookmarks.SetShowHelp(false)
		return style.Default().Render(v.bookmarks.View() + "\n" + v.input.View())
	}

	v.bookmarks.SetShowPagination(true)
	v.bookmarks.SetShowHelp(true)
	return style.Default().Render(v.bookmarks.View())
}

func (v *View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var inputCmd tea.Cmd
	var viewCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		x, y := style.Default().GetFrameSize()
		width := msg.Width - x
		height := msg.Height - y
		v.input.Width = width - x
		v.bookmarks.SetSize(width, height)
		v.bookmarks, inputCmd = v.bookmarks.Update(msg)
		v.input, viewCmd = v.input.Update(msg)
	case tea.KeyMsg:
		// Bypass "msg" input pipeline if setting filter value.
		if v.bookmarks.SettingFilter() {
			v.bookmarks, viewCmd = v.bookmarks.Update(msg)
			break
		}

		if v.input.Focused() {
			v.input, inputCmd = v.handleInput(msg)
			break
		}

		v.bookmarks, viewCmd = v.handleList(msg)
	default:
		v.input, inputCmd = v.input.Update(msg)
		v.bookmarks, viewCmd = v.bookmarks.Update(msg)
	}

	return v, tea.Batch(inputCmd, viewCmd)
}

func (v *View) handleInput(msg tea.KeyMsg) (textinput.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", "q":
		item := v.bookmarks.SelectedItem().(*model.Bookmark)
		if v.input.Value() != item.Title() {
			item.Update(v.input.Value())
			v.dirty = true
		}

		v.input.Reset()
		v.input.Blur()
		return v.input, tea.ClearScreen
	}

	return v.input.Update(msg)
}

func (v *View) handleList(msg tea.KeyMsg) (list.Model, tea.Cmd) {
	item, ok := v.bookmarks.SelectedItem().(*model.Bookmark)
	if !ok {
		return v.bookmarks.Update(msg)
	}

	switch msg.String() {
	case "enter", " ":
		_ = open(item.Description())
	case "d", "delete":
		v.bookmarks.RemoveItem(slices.Index(v.bookmarks.Items(), v.bookmarks.SelectedItem()))
		v.dirty = true
		return v.bookmarks, v.bookmarks.NewStatusMessage(fmt.Sprintf(msgStatus, item.Title()))
	case "r":
		v.input.SetValue(item.Title())
		v.input.Focus()
		return v.bookmarks, textinput.Blink
	}

	return v.bookmarks.Update(msg)
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
