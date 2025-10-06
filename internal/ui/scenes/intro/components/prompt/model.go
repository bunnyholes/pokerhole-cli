package prompt

import tea "github.com/charmbracelet/bubbletea"

// Model represents the prompt component (stateless).
type Model struct {
	width int
	text  string
}

// New creates a new prompt component.
func New(width int, text string) Model {
	return Model{
		width: width,
		text:  text,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}
