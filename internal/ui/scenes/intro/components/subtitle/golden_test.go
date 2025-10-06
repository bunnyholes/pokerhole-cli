package subtitle

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

// TestView_Golden_OpacityProgression tests subtitle rendering at different opacity levels.
// Each opacity level triggers different styling (color and bold).
func TestView_Golden_OpacityProgression(t *testing.T) {
	// Force TrueColor for consistent golden snapshots
	lipgloss.SetColorProfile(termenv.TrueColor)
	tests := []struct {
		name    string
		opacity float64
	}{
		{"opacity_0.0", 0.0},   // ColorTextSecondary
		{"opacity_0.2", 0.2},   // ColorVintageGoldDim
		{"opacity_0.5", 0.5},   // ColorVintageGold (no bold)
		{"opacity_0.8", 0.8},   // ColorVintageGold (bold)
		{"opacity_1.0", 1.0},   // ColorVintageGold (bold)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(80, "T E X A S   H O L D ' E M")
			m = m.SetOpacity(tt.opacity)
			output := m.View()
			golden.RequireEqual(t, []byte(output))
		})
	}
}
