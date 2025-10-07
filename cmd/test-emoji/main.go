package main

import (
	"fmt"
	"strings"

	"github.com/bunnyholes/pokerhole/client/internal/core/application/service"
)

func main() {
	// Create and start game
	game := service.NewOfflineGame("TestPlayer")
	game.Start()

	// Get game state
	gameState := game.GetGameState()

	// Check for emoji variant selectors
	fmt.Println("=== CHECKING FOR EMOJI VARIANT SELECTORS ===")
	
	for i, p := range gameState.Players {
		fmt.Printf("Player %d Hand: %s\n", i, p.Hand)
		if strings.Contains(p.Hand, "\uFE0F") {
			fmt.Printf("  ⚠ ERROR: Contains emoji variant selector!\n")
		} else {
			fmt.Printf("  ✓ OK: No emoji variant selector\n")
		}
		
		// Show hex dump
		fmt.Printf("  Hex: % X\n\n", p.Hand)
	}

	fmt.Println("Community Cards:")
	for i, card := range gameState.CommunityCards {
		fmt.Printf("  Card %d: %s\n", i, card)
		if strings.Contains(card, "\uFE0F") {
			fmt.Printf("    ⚠ ERROR: Contains emoji variant selector!\n")
		} else {
			fmt.Printf("    ✓ OK: No emoji variant selector\n")
		}
		fmt.Printf("    Hex: % X\n", card)
	}
}
