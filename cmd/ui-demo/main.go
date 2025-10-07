package main

import (
	"fmt"

	"github.com/bunnyholes/pokerhole/client/internal/core/application/service"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
)

func main() {
	// Create offline game
	game := service.NewOfflineGame("TestPlayer")
	game.Start()

	// Simulate one action to show game state
	gameState := game.GetGameState()

	// Display game state snapshot
	fmt.Println("\n=== NEW UI DESIGN PREVIEW ===")
	fmt.Println("\n■ Game State:")
	fmt.Printf("   Round: %s\n", gameState.Round)
	fmt.Printf("   Pot: %d chips\n", gameState.Pot)
	fmt.Printf("   Current Bet: %d chips\n", gameState.CurrentBet)
	fmt.Printf("   Current Player: %d\n", gameState.CurrentPlayer)

	fmt.Println("\n[Community Cards]")
	if len(gameState.CommunityCards) == 0 {
		fmt.Println("   (No community cards yet - PRE_FLOP)")
	} else {
		for _, card := range gameState.CommunityCards {
			fmt.Printf("   %s\n", card)
		}
	}

	fmt.Println("\n[Players]")
	for i, p := range gameState.Players {
		marker := "  "
		if i == gameState.CurrentPlayer {
			marker = "▶"
		}
		fmt.Printf("   %s Player %d: %s\n", marker, i+1, p.Nickname)
		fmt.Printf("      ◆ Chips: %d | ○ Bet: %d | Status: %s\n", p.Chips, p.Bet, p.Status)
		fmt.Printf("      Hand: %s\n", p.Hand)
	}

	// Progress to FLOP to show more UI elements
	fmt.Println("\n--- Calling to progress to FLOP ---")
	game.PlayerAction(0, vo.Call, 0)
	game.PlayerAction(1, vo.Check, 0)
	game.ProgressRound()

	gameState = game.GetGameState()
	fmt.Println("\n■ Game State After FLOP:")
	fmt.Printf("   Round: %s\n", gameState.Round)
	fmt.Printf("   Pot: %d chips\n", gameState.Pot)

	fmt.Println("\n[Community Cards (FLOP)]")
	for _, card := range gameState.CommunityCards {
		fmt.Printf("   %s\n", card)
	}

	fmt.Println("\n* UI Features Implemented:")
	fmt.Println("   • Korean interface (라운드, 칩, 베팅, etc.)")
	fmt.Println("   • Multiple border styles (DoubleBorder, RoundedBorder, ThickBorder)")
	fmt.Println("   • Color scheme (7 different colors for elements)")
	fmt.Println("   • Special character indicators (▶ active, ● status, etc.)")
	fmt.Println("   • Active player highlighting with green border")
	fmt.Println("   • Showdown display with golden double border")
	fmt.Println("   • Organized layout with sections")

	fmt.Println("\n[Controls]")
	fmt.Println("   [F] Fold   [C] Call   [R] Raise   [A] All-in   [K] Check   [N] New Game")
}
