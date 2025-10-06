package intro

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro/components/prompt"
	"github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro/components/subtitle"
	"github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro/components/title"
)

// TickMsg is sent periodically to advance the animation.
type TickMsg time.Time

// DoneMsg is sent when the intro animation completes or is skipped.
type DoneMsg struct{}

// Update implements tea.Model.
// Handles animation ticks and key presses (any key skips intro).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Any key press skips the intro
		return m, func() tea.Msg { return DoneMsg{} }

	case TickMsg:
		// Advance components based on current phase
		switch m.phase {
		case PhaseTyping:
			// Advance title typing animation
			updatedTitle, _ := m.titleModel.Update(title.TickMsg{})
			m.titleModel = updatedTitle.(title.Model)

			// Check if typing is complete
			if m.titleModel.IsComplete() {
				m.phase = PhaseSubtitle
			}

		case PhaseSubtitle:
			// Advance subtitle opacity animation
			updatedSubtitle, _ := m.subtitleModel.Update(subtitle.TickMsg{})
			m.subtitleModel = updatedSubtitle.(subtitle.Model)

			// Check if subtitle is complete
			if m.subtitleModel.IsComplete() {
				m.phase = PhaseHold
			}

		case PhaseHold:
			m.holdTicks++
			if m.holdTicks >= 90 {
				m.phase = PhaseDone
			}

		case PhaseDone:
			return m, func() tea.Msg { return DoneMsg{} }
		}

		return m, tickCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.titleModel = title.New(msg.Width, "POKERHOLE")
		m.subtitleModel = subtitle.New(msg.Width, "T E X A S   H O L D ' E M")
		m.promptModel = prompt.New(msg.Width, "아무 키나 눌러 바로 시작하기")
		return m, nil
	}

	return m, nil
}

// tickCmd returns a command that sends a TickMsg after the animation interval.
func tickCmd() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
