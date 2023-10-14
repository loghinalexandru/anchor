package bubbletea

import (
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/loghinalexandru/anchor/internal/bookmark"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.NoColor{})
	d.Styles.SelectedDesc = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.NoColor{})
	d.UpdateFunc = update

	return d
}

func update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			item := m.SelectedItem().(*bookmark.Bookmark)
			_ = open(item.URL)
		}
	}

	return nil
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
