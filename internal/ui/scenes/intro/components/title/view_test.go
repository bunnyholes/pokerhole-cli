package title

import (
	"strings"
	"testing"
)

// TestView_Empty tests rendering with no characters revealed
func TestView_Empty(t *testing.T) {
	m := New(80, "TEST")

	view := m.View()

	if view != "" {
		t.Errorf("expected empty view, got: %s", view)
	}
}

// TestView_Partial tests rendering with partial word
func TestView_Partial(t *testing.T) {
	m := New(80, "POKER")
	m = m.SetCharsRevealed(3) // "POK"

	view := m.View()

	if view == "" {
		t.Fatal("expected non-empty view")
	}

	// Should contain ASCII art representation (standard font uses _/\| characters)
	if !strings.Contains(view, "_") && !strings.Contains(view, "|") {
		t.Error("expected view to contain ASCII art characters from standard font")
	}

	// Should be multiple lines
	lines := strings.Split(view, "\n")
	if len(lines) < 5 {
		t.Errorf("expected at least 5 lines, got %d", len(lines))
	}
}

// TestView_Complete tests rendering with full word
func TestView_Complete(t *testing.T) {
	m := New(80, "TEST")
	m = m.SetCharsRevealed(4)

	view := m.View()

	if view == "" {
		t.Fatal("expected non-empty view")
	}

	// Verify it contains ASCII art (standard font uses _/\| characters)
	if !strings.Contains(view, "_") && !strings.Contains(view, "|") {
		t.Error("expected view to contain ASCII art characters from standard font")
	}
}

// TestView_DynamicWords tests rendering different words
func TestView_DynamicWords(t *testing.T) {
	words := []string{"HELLO", "WORLD", "GO", "TEST123"}

	for _, word := range words {
		t.Run(word, func(t *testing.T) {
			m := New(80, word)
			m = m.SetCharsRevealed(len(word))

			view := m.View()

			if view == "" {
				t.Errorf("expected non-empty view for word %s", word)
			}

			// Standard font should produce ASCII art
			lines := strings.Split(view, "\n")
			if len(lines) == 0 {
				t.Errorf("expected multiple lines for word %s", word)
			}
		})
	}
}

// TestView_Width tests that view respects width
func TestView_Width(t *testing.T) {
	widths := []int{60, 80, 100}

	for _, width := range widths {
		m := New(width, "HI")
		m = m.SetCharsRevealed(2)

		view := m.View()

		if view == "" {
			t.Fatalf("expected non-empty view for width %d", width)
		}

		// Check that lines don't exceed width (accounting for ANSI codes)
		// This is a rough check since lipgloss adds styling
		lines := strings.Split(view, "\n")
		for i, line := range lines {
			// Remove ANSI codes for width check
			stripped := stripAnsi(line)
			if len(stripped) > width+10 { // +10 for padding tolerance
				t.Errorf("line %d exceeds width %d: got %d chars", i, width, len(stripped))
			}
		}
	}
}

// stripAnsi removes ANSI escape codes (simple version)
func stripAnsi(s string) string {
	// Simple ANSI code removal for testing
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			result += string(r)
		}
	}
	return result
}
