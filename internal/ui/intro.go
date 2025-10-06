package ui

import (
	intro "github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) advanceIntroAnimation() (tea.Model, tea.Cmd) {
	// Forward TickMsg to intro model
	newIntroModel, cmd := m.introModel.Update(intro.TickMsg{})
	m.introModel = newIntroModel.(intro.Model)

	// Check if intro is done
	if cmd != nil {
		if doneMsg := cmd(); doneMsg != nil {
			if _, ok := doneMsg.(intro.DoneMsg); ok {
				m.screen = screenHome
				return m, animationTickCmd()
			}
		}
	}

	return m, tea.Batch(cmd, animationTickCmd())
}
