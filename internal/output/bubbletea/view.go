package bubbletea

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View struct {
	bookmarks list.Model
	style     lipgloss.Style
}

func NewView(bookmarks []list.Item) *View {
	viewList := list.New(bookmarks, newItemDelegate(), 0, 0)
	viewList.Title = "Bookmarks"
	viewList.InfiniteScrolling = true
	viewList.Paginator.Type = paginator.Arabic
	viewList.Paginator.ArabicFormat = "%d/%d \u2693"
	viewList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})
	viewList.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})

	return &View{
		bookmarks: viewList,
		style:     lipgloss.NewStyle().Margin(0, 0, 0, 2),
	}
}

func (v *View) Init() tea.Cmd {
	return nil
}

func (v *View) View() string {
	return v.style.Render(v.bookmarks.View())
}

func (v *View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return v, tea.Quit
		}
	case tea.WindowSizeMsg:
		x, y := v.style.GetFrameSize()
		v.bookmarks.SetSize(msg.Width-x, msg.Height-y)
	}

	var cmd tea.Cmd
	v.bookmarks, cmd = v.bookmarks.Update(msg)
	return v, cmd
}
