package bubbletea

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/loghinalexandru/anchor/internal/model"
)

const (
	msgStatus = "Deleted %q"
)

type View struct {
	input     textinput.Model
	bookmarks list.Model
	style     lipgloss.Style
	dirty     bool
}

func NewView(bookmarks []list.Item) *View {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.NoColor{})
	d.Styles.SelectedDesc = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.NoColor{})

	viewList := list.New(bookmarks, d, 0, 0)
	viewList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	viewList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	viewList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("enter", "open")),
			key.NewBinding(key.WithKeys("delete", "d"), key.WithHelp("d/del", "delete")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename")),
		}
	}

	viewList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("enter/space", "open in browser")),
			key.NewBinding(key.WithKeys("delete", "d"), key.WithHelp("d/del", "remove bookmark")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename bookmark")),
		}
	}

	viewList.Title = "Bookmarks"
	viewList.InfiniteScrolling = true
	viewList.Paginator.Type = paginator.Arabic
	viewList.Paginator.ArabicFormat = "%d/%d \u2693"
	viewList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})
	viewList.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})

	return &View{
		input:     textinput.New(),
		bookmarks: viewList,
		style:     lipgloss.NewStyle().Margin(2, 0, 2, 2),
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
		return v.style.Render(v.bookmarks.View() + "\n" + v.input.View())
	}

	v.bookmarks.SetShowPagination(true)
	v.bookmarks.SetShowHelp(true)
	return v.style.Render(v.bookmarks.View())
}

func (v *View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var inputCmd tea.Cmd
	var viewCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		x, y := v.style.GetFrameSize()
		v.bookmarks.SetSize(msg.Width-x, msg.Height-y)
		v.bookmarks, viewCmd = v.bookmarks.Update(msg)
		return v, viewCmd
	case tea.KeyMsg:
		if v.input.Focused() {
			v.input, inputCmd = v.handleInput(msg)
		} else {
			v.bookmarks, viewCmd = v.handleList(msg)
		}
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
	switch msg.String() {
	case "enter", "space":
		item := v.bookmarks.SelectedItem().(*model.Bookmark)
		_ = open(item.Description())
	case "d", "delete":
		item := v.bookmarks.SelectedItem().(*model.Bookmark)
		v.bookmarks.RemoveItem(v.bookmarks.Index())
		v.dirty = true
		return v.bookmarks, v.bookmarks.NewStatusMessage(fmt.Sprintf(msgStatus, item.Title()))
	case "r":
		item := v.bookmarks.SelectedItem().(*model.Bookmark)
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
