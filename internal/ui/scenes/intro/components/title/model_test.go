package title

import (
	"testing"
)

// TestModel_New tests model creation
func TestModel_New(t *testing.T) {
	m := New(80, "POKER")

	if m.width != 80 {
		t.Errorf("expected width 80, got %d", m.width)
	}
	if m.word != "POKER" {
		t.Errorf("expected word POKER, got %s", m.word)
	}
	if m.charsRevealed != 0 {
		t.Errorf("expected charsRevealed 0, got %d", m.charsRevealed)
	}
}

// TestModel_SetCharsRevealed tests updating revealed characters
func TestModel_SetCharsRevealed(t *testing.T) {
	m := New(80, "TEST")
	m = m.SetCharsRevealed(2)

	if m.CharsRevealed() != 2 {
		t.Errorf("expected 2 chars revealed, got %d", m.CharsRevealed())
	}
}

// TestModel_IsComplete tests completion detection
func TestModel_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		revealed int
		complete bool
	}{
		{"not started", "POKER", 0, false},
		{"partial", "POKER", 3, false},
		{"exact", "POKER", 5, true},
		{"overflow", "POKER", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(80, tt.word)
			m = m.SetCharsRevealed(tt.revealed)

			if m.IsComplete() != tt.complete {
				t.Errorf("expected complete=%v, got %v", tt.complete, m.IsComplete())
			}
		})
	}
}

// TestModel_Update tests update with TickMsg
func TestModel_Update(t *testing.T) {
	m := New(80, "TEST")

	// Send tick message
	updated, _ := m.Update(TickMsg{})
	m = updated.(Model)

	if m.CharsRevealed() != 1 {
		t.Errorf("expected 1 char revealed after tick, got %d", m.CharsRevealed())
	}

	// Multiple ticks
	for i := 0; i < 3; i++ {
		updated, _ = m.Update(TickMsg{})
		m = updated.(Model)
	}

	if m.CharsRevealed() != 4 {
		t.Errorf("expected 4 chars revealed, got %d", m.CharsRevealed())
	}

	// Complete - no more updates
	updated, _ = m.Update(TickMsg{})
	m = updated.(Model)

	if m.CharsRevealed() != 4 {
		t.Errorf("expected to stay at 4 chars when complete, got %d", m.CharsRevealed())
	}
}
