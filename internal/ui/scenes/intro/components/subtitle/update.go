package subtitle

import tea "github.com/charmbracelet/bubbletea"

// TickMsg is sent to advance the opacity animation.
type TickMsg struct{}

// OpacityStep is how much opacity increases per tick.
const OpacityStep = 0.05

// Update implements tea.Model.
// The parent scene controls the animation timing by sending TickMsg.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		if !m.IsComplete() {
			m.opacity += OpacityStep
			if m.opacity > 1.0 {
				m.opacity = 1.0
			}
		}
		return m, nil
	}
	return m, nil
}
