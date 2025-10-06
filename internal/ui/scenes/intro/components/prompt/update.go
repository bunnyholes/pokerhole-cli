package prompt

import tea "github.com/charmbracelet/bubbletea"

// Update implements tea.Model.
// Prompt is stateless, so it doesn't respond to any messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
