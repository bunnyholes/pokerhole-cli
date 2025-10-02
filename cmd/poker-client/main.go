package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bunnyholes/pokerhole/client/internal/identity"
	"github.com/bunnyholes/pokerhole/client/internal/network"
	"github.com/bunnyholes/pokerhole/client/internal/ui"
)

func main() {
	// Set up logging
	logFile, err := os.OpenFile("/tmp/pokerhole-client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	// Get or create UUID and nickname
	uuid := identity.GetOrCreateUUID()
	nickname := identity.GenerateNickname()

	log.Printf("Starting client with UUID: %s, Nickname: %s", uuid, nickname)

	// Create WebSocket client
	serverURL := getServerURL()
	client := network.NewClient(serverURL, uuid, nickname)

	// Connect to server
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	log.Printf("Connected to server: %s", serverURL)

	// Create and run Bubble Tea program
	p := tea.NewProgram(ui.NewModel(client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
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
