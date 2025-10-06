package prompt

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/bunnyholes/pokerhole/client/internal/ui/constants"
)

// View implements tea.Model.
// Renders the prompt text.
func (m Model) View() string {
	style := lipgloss.NewStyle().
		Foreground(constants.ColorTextSecondary).
		Align(lipgloss.Center).
		Width(m.width)

	return style.Render(m.text)
}

// Height returns the component's height in lines.
func (m Model) Height() int {
	return 1
}
