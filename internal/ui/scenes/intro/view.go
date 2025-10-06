package intro

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/bunnyholes/pokerhole/client/internal/ui/constants"
)

// View implements tea.Model.
// Composes sub-components into the final scene layout.
func (m Model) View() string {
	var components []string

	// Top spacing
	components = append(components, "")

	// Title component (typing animation)
	if titleView := m.titleModel.View(); titleView != "" {
		components = append(components, titleView)
	}

	// Subtitle component (opacity animation)
	if m.phase >= PhaseSubtitle {
		components = append(components, m.subtitleModel.View())
	} else {
		components = append(components, "")
	}

	// Spacing before prompt
	components = append(components, "")

	// Prompt component
	components = append(components, m.promptModel.View())

	// Join all components
	content := strings.Join(components, "\n")

	// Center content vertically within terminal height
	return lipgloss.Place(
		m.width,
		constants.TerminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
