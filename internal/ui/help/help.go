package help

// KeyBinding represents a single key binding and its description.
type KeyBinding struct {
	Key         string
	Description string
}

// Provider is an interface for components that can provide help bindings.
type Provider interface {
	Bindings() []KeyBinding
}

// Render formats key bindings for display (future implementation).
// This is a placeholder for the future help UI renderer.
func Render(bindings []KeyBinding) string {
	// TODO: Implement help UI rendering with lipgloss
	return ""
}
