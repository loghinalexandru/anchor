package bubbletea

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type Input struct {
	model textinput.Model
	err   error
}

func NewInput(placeholder string) Input {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()

	return Input{
		model: ti,
		err:   nil,
	}
}

func (input Input) Init() tea.Cmd {
	return textinput.Blink
}

func (input Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return input, tea.Quit
		}

	case errMsg:
		input.err = msg
		return input, nil
	}

	input.model, cmd = input.model.Update(msg)
	return input, cmd
}

func (input Input) View() string {
	return input.model.View()
}
