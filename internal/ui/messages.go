package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bunnyholes/pokerhole/client/internal/network"
)

type animationTickMsg struct{}

type aiTurnMsg struct{}

type statusClearMsg struct {
	seq int
}

type serverMessageMsg struct {
	Message network.ServerMessage
}

func animationTickCmd() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(time.Time) tea.Msg {
		return animationTickMsg{}
	})
}

func listenForMessages(client *network.Client) tea.Cmd {
	if client == nil {
		return nil
	}

	return func() tea.Msg {
		msg, ok := <-client.Receive()
		if !ok {
			return nil
		}
		return serverMessageMsg{Message: msg}
	}
}
