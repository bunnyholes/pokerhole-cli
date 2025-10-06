package subtitle

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/bunnyholes/pokerhole/client/internal/ui/constants"
)

// View implements tea.Model.
// Renders the subtitle with opacity-based styling.
func (m Model) View() string {
	// Select style based on opacity
	var style lipgloss.Style
	switch {
	case m.opacity > 0.7:
		style = lipgloss.NewStyle().Foreground(constants.ColorVintageGold).Bold(true)
	case m.opacity > 0.4:
		style = lipgloss.NewStyle().Foreground(constants.ColorVintageGold)
	case m.opacity > 0.1:
		style = lipgloss.NewStyle().Foreground(constants.ColorVintageGoldDim)
	default:
		style = lipgloss.NewStyle().Foreground(constants.ColorTextSecondary)
	}

	// Center align
	center := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width)

	return center.Render(style.Render(m.text))
}

// Height returns the component's height in lines.
func (m Model) Height() int {
	return 1
}
