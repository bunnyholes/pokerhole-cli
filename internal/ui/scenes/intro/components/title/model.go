package title

import tea "github.com/charmbracelet/bubbletea"

// Model represents the title component with typing animation.
type Model struct {
	width         int
	word          string
	charsRevealed int
}

// New creates a new title component.
func New(width int, word string) Model {
	return Model{
		width:         width,
		word:          word,
		charsRevealed: 0,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil // Animation is driven by parent
}

// CharsRevealed returns the current number of revealed characters (for testing).
func (m Model) CharsRevealed() int {
	return m.charsRevealed
}

// SetCharsRevealed updates the number of revealed characters.
func (m Model) SetCharsRevealed(chars int) Model {
	m.charsRevealed = chars
	return m
}

// IsComplete returns true when all characters are revealed.
func (m Model) IsComplete() bool {
	return m.charsRevealed >= len(m.word)
}
