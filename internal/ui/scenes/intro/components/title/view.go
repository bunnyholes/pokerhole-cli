package title

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/bunnyholes/pokerhole/client/internal/ui/constants"
)

// View implements tea.Model.
// Renders the title with typing animation based on charsRevealed.
func (m Model) View() string {
	// Get visible portion of the word
	visibleWord := ""
	if m.charsRevealed > 0 && m.charsRevealed <= len(m.word) {
		visibleWord = m.word[:m.charsRevealed]
	}

	if visibleWord == "" {
		return ""
	}

	// Generate ASCII art using go-figure (standard font)
	fig := figure.NewFigure(visibleWord, "standard", true)
	asciiArt := fig.String()

	// Split into lines
	lines := strings.Split(strings.TrimSuffix(asciiArt, "\n"), "\n")

	// Apply vintage gold styling
	goldStyle := lipgloss.NewStyle().Foreground(constants.ColorVintageGold).Bold(true)
	var styled []string
	for _, line := range lines {
		styled = append(styled, goldStyle.Render(line))
	}

	// Left align
	leftAlign := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.width)

	var aligned []string
	for _, line := range styled {
		aligned = append(aligned, leftAlign.Render(line))
	}

	return strings.Join(aligned, "\n")
}

// Height returns the component's height in lines.
// Standard font typically produces 6 lines.
func (m Model) Height() int {
	return 6
}
