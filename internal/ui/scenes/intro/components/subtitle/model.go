package subtitle

import tea "github.com/charmbracelet/bubbletea"

// Model represents the subtitle component with opacity animation.
type Model struct {
	width   int
	text    string
	opacity float64
}

// New creates a new subtitle component.
func New(width int, text string) Model {
	return Model{
		width:   width,
		text:    text,
		opacity: 0.0,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil // Animation is driven by parent
}

// Opacity returns the current opacity (for testing).
func (m Model) Opacity() float64 {
	return m.opacity
}

// SetOpacity updates the opacity value.
func (m Model) SetOpacity(opacity float64) Model {
	m.opacity = opacity
	return m
}

// IsComplete returns true when opacity reaches 1.0.
func (m Model) IsComplete() bool {
	return m.opacity >= 1.0
}
