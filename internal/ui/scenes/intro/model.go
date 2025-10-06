package intro

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro/components/prompt"
	"github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro/components/subtitle"
	"github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro/components/title"
)

// Phase represents the current step of the intro animation.
type Phase int

const (
	PhaseTyping Phase = iota
	PhaseSubtitle
	PhaseHold
	PhaseDone
)

// Model represents the intro scene that orchestrates sub-components.
type Model struct {
	phase        Phase
	width        int
	holdTicks    int
	titleModel   title.Model
	subtitleModel subtitle.Model
	promptModel  prompt.Model
}

// NewModel creates a new intro scene model with the specified terminal width.
func NewModel(width int) Model {
	return Model{
		phase:         PhaseTyping,
		width:         width,
		holdTicks:     0,
		titleModel:    title.New(width, "POKERHOLE"),
		subtitleModel: subtitle.New(width, "T E X A S   H O L D ' E M"),
		promptModel:   prompt.New(width, "아무 키나 눌러 바로 시작하기"),
	}
}

// Init implements tea.Model.
// Returns a command to start animation ticks.
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// Phase returns the current phase (for testing).
func (m Model) Phase() Phase {
	return m.phase
}

// TitleModel returns the title component (for testing).
func (m Model) TitleModel() title.Model {
	return m.titleModel
}

// SubtitleModel returns the subtitle component (for testing).
func (m Model) SubtitleModel() subtitle.Model {
	return m.subtitleModel
}
