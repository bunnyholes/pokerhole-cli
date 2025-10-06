package intro

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

// TestView_Golden_PhaseProgression tests intro scene rendering across all phases.
// Validates scene-level composition with lipgloss.Place centering (28x80 terminal).
func TestView_Golden_PhaseProgression(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	tests := []struct {
		name           string
		phase          Phase
		titleChars     int
		subtitleOpacity float64
	}{
		{
			name:           "PhaseTyping",
			phase:          PhaseTyping,
			titleChars:     5, // "POKER"
			subtitleOpacity: 0.0,
		},
		{
			name:           "PhaseSubtitle",
			phase:          PhaseSubtitle,
			titleChars:     9, // "POKERHOLE"
			subtitleOpacity: 0.5,
		},
		{
			name:           "PhaseHold",
			phase:          PhaseHold,
			titleChars:     9,
			subtitleOpacity: 1.0,
		},
		{
			name:           "PhaseDone",
			phase:          PhaseDone,
			titleChars:     9,
			subtitleOpacity: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(80)
			m.phase = tt.phase
			m.titleModel = m.titleModel.SetCharsRevealed(tt.titleChars)
			m.subtitleModel = m.subtitleModel.SetOpacity(tt.subtitleOpacity)
			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}

// TestView_Golden_TitleProgressionInScene tests how title typing animation
// affects the overall scene composition with vertical centering.
func TestView_Golden_TitleProgressionInScene(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	tests := []struct {
		name       string
		titleChars int
	}{
		{"empty", 0},
		{"POK", 3},
		{"POKERH", 6},
		{"POKERHOLE", 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(80)
			m.phase = PhaseTyping
			m.titleModel = m.titleModel.SetCharsRevealed(tt.titleChars)
			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}

// TestView_Golden_SubtitleProgressionInScene tests how subtitle opacity animation
// affects the overall scene composition.
func TestView_Golden_SubtitleProgressionInScene(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	tests := []struct {
		name    string
		opacity float64
	}{
		{"opacity_0.0", 0.0},
		{"opacity_0.3", 0.3},
		{"opacity_0.7", 0.7},
		{"opacity_1.0", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(80)
			m.phase = PhaseSubtitle // Ensure subtitle is visible
			m.titleModel = m.titleModel.SetCharsRevealed(9) // Full title
			m.subtitleModel = m.subtitleModel.SetOpacity(tt.opacity)
			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}
