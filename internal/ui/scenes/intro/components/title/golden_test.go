package title

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

// TestView_Golden_TypingProgression tests the typing animation progression
func TestView_Golden_TypingProgression(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	word := "POKERHOLE"
	width := 80

	tests := []struct {
		name          string
		charsRevealed int
	}{
		{"empty", 0},
		{"P", 1},
		{"PO", 2},
		{"POK", 3},
		{"POKE", 4},
		{"POKER", 5},
		{"POKERH", 6},
		{"POKERHO", 7},
		{"POKERHOL", 8},
		{"POKERHOLE", 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(width, word)
			m = m.SetCharsRevealed(tt.charsRevealed)

			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}

// TestView_Golden_DifferentWords tests rendering different words
func TestView_Golden_DifferentWords(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	width := 80

	tests := []struct {
		name string
		word string
	}{
		{"HELLO", "HELLO"},
		{"WORLD", "WORLD"},
		{"GO", "GO"},
		{"POKER", "POKER"},
		{"TEST", "TEST"},
		{"ABC123", "ABC123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(width, tt.word)
			m = m.SetCharsRevealed(len(tt.word))

			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}

// TestView_Golden_DifferentWidths tests rendering at different widths
func TestView_Golden_DifferentWidths(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	word := "TEST"

	tests := []struct {
		name  string
		width int
	}{
		{"width_60", 60},
		{"width_80", 80},
		{"width_100", 100},
		{"width_120", 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.width, word)
			m = m.SetCharsRevealed(len(word))

			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}
