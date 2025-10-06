package prompt

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

// TestView_Golden_DifferentTexts tests prompt rendering with various text content.
// Tests centering with different text lengths and character sets.
func TestView_Golden_DifferentTexts(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	tests := []struct {
		name string
		text string
	}{
		{"korean_default", "아무 키나 눌러 바로 시작하기"},
		{"english", "Press any key to start"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(80, tt.text)
			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}
