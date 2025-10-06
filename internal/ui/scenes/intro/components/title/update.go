package title

import tea "github.com/charmbracelet/bubbletea"

// TickMsg is sent to advance the typing animation.
type TickMsg struct{}

// Update implements tea.Model.
// The parent scene controls the animation timing by sending TickMsg.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		if !m.IsComplete() {
			m.charsRevealed++
		}
		return m, nil
	}
	return m, nil
}
