package bubbletea

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 0, 0, 2)

type View struct {
	bookmarks list.Model
}

func NewView(bookmarks []list.Item) View {
	viewList := list.New(bookmarks, newItemDelegate(), 0, 0)
	viewList.Title = "Bookmarks"

	return View{
		bookmarks: viewList,
	}
}

func (m View) Init() tea.Cmd {
	return nil
}

func (m View) View() string {
	return docStyle.Render(m.bookmarks.View())
}

func (m View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.bookmarks.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.bookmarks, cmd = m.bookmarks.Update(msg)
	return m, cmd
}
