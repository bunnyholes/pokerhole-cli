package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"github.com/bunnyholes/pokerhole/client/internal/identity"
	"github.com/bunnyholes/pokerhole/client/internal/network"
	"github.com/bunnyholes/pokerhole/client/internal/ui"
)

func main() {
	// Generate new UUID for each session (for dev/testing with multiple clients)
	// TODO: Use GetOrCreateUUID() for production to maintain user identity
	clientUUID := uuid.New().String()
	nickname := identity.GenerateNickname()

	// Create WebSocket client
	serverURL := getServerURL()

	client := network.NewClient(serverURL, clientUUID, nickname)

	// Connect to server with 3 second timeout
	var isOnline bool
	if err := client.ConnectWithTimeout(3 * time.Second); err != nil {
		isOnline = false
	} else {
		defer client.Close()
		isOnline = true
	}

	// Create and run Bubble Tea program (works in both online and offline modes)
	p := tea.NewProgram(
		ui.NewModel(client, isOnline, nickname),
		tea.WithAltScreen(),
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	)
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}

	fmt.Println("\n오늘도 편안한 하루 보내세요.")

	// 명시적 종료 (defer 함수 실행 후 즉시 종료)
	os.Exit(0)
}

// getServerURL returns the WebSocket server URL
func getServerURL() string {
	// Check environment variable first
	if url := os.Getenv("POKERHOLE_SERVER"); url != "" {
		return url
	}

	// Default to localhost with new protocol endpoint
	return "ws://localhost:8080/ws/game"
}
