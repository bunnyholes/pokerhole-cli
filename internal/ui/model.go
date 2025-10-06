package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bunnyholes/pokerhole/client/internal/core/application/service"
	"github.com/bunnyholes/pokerhole/client/internal/network"
	intro "github.com/bunnyholes/pokerhole/client/internal/ui/scenes/intro"
)

// screenID represents a primary screen in the CLI application.
type screenID string

const (
	screenIntro screenID = "intro"
	screenHome  screenID = "home"
	screenGame  screenID = "game"
)

// modalID represents modal overlays rendered above the primary screen.
type modalID string

const (
	modalNone     modalID = ""
	modalHelp     modalID = "help"
	modalAbout    modalID = "about"
	modalShowdown modalID = "showdown"
)

// statusLevel controls the accent color of the status bar.
type statusLevel int

const (
	statusNeutral statusLevel = iota
	statusInfo
	statusSuccess
	statusWarning
	statusError
)

// homeAction enumerates possible primary menu actions.
type homeAction int

const (
	homeActionOffline homeAction = iota
	homeActionOnlineMatch
	homeActionQuit
)

type menuItem struct {
	title       string
	description string
	action      homeAction
	disabled    bool
	disabledMsg string
}

type homeState struct {
	items    []menuItem
	selected int
}

type gameState struct {
	offlineGame *service.OfflineGame
	snapshot    service.GameStateSnapshot
}

type statusState struct {
	message string
	level   statusLevel
	seq     int
}

// Model is the root Bubble Tea model for the PokerHole CLI.
type Model struct {
	client     *network.Client
	playerName string
	online     bool

	spinner spinner.Model

	screen screenID
	modal  modalID

	width  int
	height int

	introModel intro.Model // Updated: Now using intro.Model instead of intro.State
	home       homeState
	game       gameState

	status statusState
}

// NewModel constructs the CLI application model.
func NewModel(client *network.Client, online bool, playerName string) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle()

	m := Model{
		client:     client,
		playerName: playerName,
		online:     online,
		spinner:    sp,
		screen:     screenIntro,
		modal:      modalNone,
		introModel: intro.NewModel(80), // Updated: Initialize with intro.NewModel
	}

	m.home.items = m.buildHomeMenu()

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.introModel.Init(), // Updated: Initialize intro model
		animationTickCmd(),
	}

	if m.online && m.client != nil {
		cmds = append(cmds, listenForMessages(m.client))
	}

	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Forward to intro model
		if m.screen == screenIntro {
			newIntroModel, _ := m.introModel.Update(msg)
			m.introModel = newIntroModel.(intro.Model)
		}

		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case animationTickMsg:
		return m.handleAnimationTick()

	case serverMessageMsg:
		return m.handleServerMessage(msg)

	case aiTurnMsg:
		return m.handleAITurn()

	case statusClearMsg:
		if msg.seq == m.status.seq {
			m.status.message = ""
			m.status.level = statusNeutral
		}
		return m, nil
	}

	return m, nil
}

// View renders the current screen (and modal if present).
func (m Model) View() string {
	var content string

	switch m.screen {
	case screenIntro:
		// Updated: Render intro scene model directly
		body := m.introModel.View()
		content = m.applyShell(body)
	case screenHome:
		content = m.viewHome()
	case screenGame:
		content = m.viewOfflineGame()
	default:
		content = ""
	}

	if m.modal != modalNone && m.screen != screenGame {
		return m.renderModalOverlay(content)
	}

	return content
}

// --- High level handlers ----------------------------------------------------

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key := msg.String(); key == "ctrl+c" {
		return m, tea.Quit
	}

	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
		switch msg.Runes[0] {
		case '?':
			m.modal = modalHelp
			return m, nil
		case 'h', 'H':
			m.modal = modalAbout
			return m, nil
		}
	}

	if m.modal != modalNone {
		return m.handleModalKey(msg)
	}

	switch m.screen {
	case screenIntro:
		// Updated: Handle intro key press through introModel
		newIntroModel, cmd := m.introModel.Update(msg)
		m.introModel = newIntroModel.(intro.Model)

		// If DoneMsg received, transition to home
		if cmd != nil {
			if doneMsg := cmd(); doneMsg != nil {
				if _, ok := doneMsg.(intro.DoneMsg); ok {
					m.screen = screenHome
					return m, nil
				}
			}
		}
		return m, cmd

	case screenHome:
		return m.handleHomeKey(msg)

	case screenGame:
		return m.handleGameKey(msg)
	}

	return m, nil
}

func (m Model) handleAnimationTick() (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenIntro:
		return m.advanceIntroAnimation()
	case screenGame:
		m.game.snapshot = m.currentSnapshot()
	}
	return m, animationTickCmd()
}

func (m Model) handleServerMessage(msg serverMessageMsg) (tea.Model, tea.Cmd) {
	// Placeholder: update status bar for now.
	m = m.withStatus(statusInfo, fmt.Sprintf("서버 이벤트: %s", msg.Message.Type), 3*time.Second)

	cmds := []tea.Cmd{
		listenForMessages(m.client),
	}

	if cmd := m.statusCommand(3 * time.Second); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleAITurn() (tea.Model, tea.Cmd) {
	if m.game.offlineGame == nil {
		return m, nil
	}

	updated, cmd := m.performAITurn()
	return updated, cmd
}

// --- Helper state mutators ---------------------------------------------------

func (m Model) withStatus(level statusLevel, message string, duration time.Duration) Model {
	m.status.level = level
	m.status.message = message
	m.status.seq++
	return m
}

func (m Model) statusCommand(duration time.Duration) tea.Cmd {
	seq := m.status.seq
	if duration <= 0 {
		return nil
	}
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return statusClearMsg{seq: seq}
	})
}

func (m Model) currentSnapshot() service.GameStateSnapshot {
	if m.game.offlineGame == nil {
		return service.GameStateSnapshot{}
	}
	return m.game.offlineGame.GetGameState()
}
