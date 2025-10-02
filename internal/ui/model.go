package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bunnyholes/pokerhole/client/internal/network"
	"github.com/bunnyholes/pokerhole/client/internal/state"
)

// ViewMode represents the current view
type ViewMode int

const (
	ViewSplash ViewMode = iota
	ViewConnecting
	ViewMenu
	ViewGame
)

// Model represents the UI model
type Model struct {
	spinner     spinner.Model
	mode        ViewMode
	client      *network.Client
	gameState   *state.GameState
	isConnected bool
	statusMsg   string
	width       int
	height      int
}

// NewModel creates a new UI model
func NewModel(client *network.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		spinner:   s,
		mode:      ViewSplash,
		client:    client,
		gameState: state.NewGameState(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForConnection(m.client),
		listenForMessages(m.client),
	)
}

// ServerMessageMsg wraps a server message
type ServerMessageMsg struct {
	Message network.ServerMessage
}

// ConnectionEstablishedMsg indicates connection success
type ConnectionEstablishedMsg struct{}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case ConnectionEstablishedMsg:
		m.isConnected = true
		m.mode = ViewMenu
		m.statusMsg = "Connected to server!"
		return m, nil

	case ServerMessageMsg:
		return m.handleServerMessage(msg.Message)
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	switch m.mode {
	case ViewSplash:
		return m.renderSplash()
	case ViewConnecting:
		return m.renderConnecting()
	case ViewMenu:
		return m.renderMenu()
	case ViewGame:
		return m.renderGame()
	}
	return ""
}

// renderSplash renders the splash screen
func (m Model) renderSplash() string {
	banner := `
    ____        __            __  __      __   
   / __ \____  / /_____  ____/ / / /___  / /__ 
  / /_/ / __ \/ //_/ _ \/ __/ /_/ / __ \/ / _ \
 / ____/ /_/ / ,< /  __/ / / __  / /_/ / /  __/
/_/    \____/_/|_|\___/_/ /_/ /_/\____/_/\___/ 
`
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	return "\n" + style.Render(banner) + "\n\n" + m.spinner.View() + " Loading...\n"
}

// renderConnecting renders the connecting screen
func (m Model) renderConnecting() string {
	return "\n" + m.spinner.View() + " Connecting to server...\n"
}

// renderMenu renders the menu
func (m Model) renderMenu() string {
	var s strings.Builder
	s.WriteString("\n🎮 PokerHole - Main Menu\n\n")
	s.WriteString(fmt.Sprintf("Status: %s\n\n", m.statusMsg))
	s.WriteString("Available options:\n")
	s.WriteString("  - Join Random Match (coming soon)\n")
	s.WriteString("  - Join Code Match (coming soon)\n")
	s.WriteString("  - Quit (q)\n\n")
	s.WriteString("Press 'q' to quit\n")
	return s.String()
}

// renderGame renders the game view
func (m Model) renderGame() string {
	snapshot := m.gameState.GetSnapshot()
	var s strings.Builder
	
	s.WriteString("\n🎲 PokerHole - Game in Progress\n\n")
	s.WriteString(fmt.Sprintf("Game ID: %s\n", snapshot.GameID))
	s.WriteString(fmt.Sprintf("Round: %s\n", snapshot.Round))
	s.WriteString(fmt.Sprintf("Pot: %d\n", snapshot.Pot))
	s.WriteString(fmt.Sprintf("Current Bet: %d\n\n", snapshot.CurrentBet))
	
	if len(snapshot.CommunityCards) > 0 {
		s.WriteString("Community Cards: ")
		s.WriteString(strings.Join(snapshot.CommunityCards, " "))
		s.WriteString("\n\n")
	}
	
	s.WriteString("Players:\n")
	for _, p := range snapshot.Players {
		s.WriteString(fmt.Sprintf("  %s - Chips: %d, Bet: %d, Status: %s\n",
			p.Nickname, p.Chips, p.Bet, p.Status))
	}
	
	return s.String()
}

// handleServerMessage processes server messages
func (m Model) handleServerMessage(msg network.ServerMessage) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case network.ServerRegisterSuccess:
		m.statusMsg = "Registration successful!"
	case network.ServerGameStateUpdate:
		m.gameState.Update(msg)
		m.mode = ViewGame
	case network.ServerError:
		if payload := msg.Payload; payload != nil {
			if errMsg, ok := payload["message"].(string); ok {
				m.statusMsg = "Error: " + errMsg
			}
		}
	}
	return m, listenForMessages(m.client)
}

// waitForConnection waits for connection to be established
func waitForConnection(client *network.Client) tea.Cmd {
	return func() tea.Msg {
		for !client.IsConnected() {
			time.Sleep(100 * time.Millisecond)
		}
		return ConnectionEstablishedMsg{}
	}
}

// listenForMessages listens for server messages
func listenForMessages(client *network.Client) tea.Cmd {
	return func() tea.Msg {
		msg := <-client.Receive()
		return ServerMessageMsg{Message: msg}
	}
}
