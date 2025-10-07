package service

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
)

// TestShowdownWinnerEvaluation tests that showdown correctly determines winner
func TestShowdownWinnerEvaluation(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	err := game.Start()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Progress through all rounds to showdown
	// PRE_FLOP
	game.PlayerAction(0, vo.Call, 0)
	game.PlayerAction(1, vo.Check, 0)
	game.ProgressRound()

	// FLOP
	game.PlayerAction(0, vo.Check, 0)
	game.PlayerAction(1, vo.Check, 0)
	game.ProgressRound()

	// TURN
	game.PlayerAction(0, vo.Check, 0)
	game.PlayerAction(1, vo.Check, 0)
	game.ProgressRound()

	// RIVER
	game.PlayerAction(0, vo.Check, 0)
	game.PlayerAction(1, vo.Check, 0)
	game.ProgressRound() // This should trigger showdown

	// Get game state after showdown
	gameState := game.GetGameState()

	// Verify round is showdown
	if gameState.Round != "SHOWDOWN" {
		t.Errorf("Expected round to be SHOWDOWN, got %s", gameState.Round)
	}

	// Verify winner was determined
	if gameState.WinnerIndex < 0 || gameState.WinnerIndex >= len(gameState.Players) {
		t.Errorf("Invalid winner index: %d", gameState.WinnerIndex)
	}

	// Verify winner hand rank is set
	if gameState.WinnerHandRank == "" {
		t.Error("Winner hand rank should not be empty")
	}

	t.Logf("Winner: Player %d with %s", gameState.WinnerIndex, gameState.WinnerHandRank)

	// Verify each player has a hand rank and rank cards
	for i, p := range gameState.Players {
		if p.HandRank == "" {
			t.Errorf("Player %d should have a hand rank", i)
		}
		if len(p.BestCards) == 0 {
			t.Errorf("Player %d should have rank cards", i)
		}
		t.Logf("Player %d (%s): %s - %s", i, p.Nickname, p.Hand, p.HandRank)
		t.Logf("  Rank cards (%d): %v", len(p.BestCards), p.BestCards)
	}

	// Verify community cards are all dealt (should be 5)
	if len(gameState.CommunityCards) != 5 {
		t.Errorf("Expected 5 community cards, got %d", len(gameState.CommunityCards))
	}

	t.Logf("Community cards: %v", gameState.CommunityCards)
}
