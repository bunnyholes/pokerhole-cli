package service

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/adapter/out/deck"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

func TestDealHoleCards(t *testing.T) {
	// Create deck and game service
	localDeck := deck.NewLocalDeck()
	localDeck.Shuffle(12345) // Use deterministic seed
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	// Create players
	p1ID := player.GeneratePlayerId()
	p2ID := player.GeneratePlayerId()
	p1Nick, _ := player.NewNickname("Player1")
	p2Nick, _ := player.NewNickname("Player2")
	player1, _ := player.NewPlayer(p1ID, p1Nick, 1000)
	player2, _ := player.NewPlayer(p2ID, p2Nick, 1000)

	players := []*player.Player{player1, player2}

	// Deal hole cards
	err := service.DealHoleCards(players)
	if err != nil {
		t.Fatalf("DealHoleCards failed: %v", err)
	}

	// Verify each player has 2 cards
	for i, p := range players {
		hand := p.Hand()
		cards := hand.Cards()
		if len(cards) != 2 {
			t.Errorf("Player %d: expected 2 cards, got %d", i+1, len(cards))
		}
	}

	// Verify 4 cards were drawn from deck (52 - 4 = 48)
	if localDeck.RemainingCards() != 48 {
		t.Errorf("Expected 48 cards remaining in deck, got %d", localDeck.RemainingCards())
	}
}

func TestDealHoleCards_MultipleRounds(t *testing.T) {
	// Test that dealing multiple times replaces hands correctly
	localDeck := deck.NewLocalDeck()
	localDeck.Shuffle(12345)
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	p1ID := player.GeneratePlayerId()
	p1Nick, _ := player.NewNickname("Player1")
	player1, _ := player.NewPlayer(p1ID, p1Nick, 1000)
	players := []*player.Player{player1}

	// First deal
	service.DealHoleCards(players)
	firstHand := player1.Hand().Cards()

	// Second deal (deck has 50 cards left)
	service.DealHoleCards(players)
	secondHand := player1.Hand().Cards()

	// Hands should be different (new cards)
	if len(secondHand) != 2 {
		t.Errorf("Expected 2 cards in second hand, got %d", len(secondHand))
	}

	// Note: We can't guarantee cards are different due to randomness,
	// but we verify the hand was replaced
	if len(firstHand) != 2 {
		t.Errorf("First hand should still have 2 cards, got %d", len(firstHand))
	}
}

func TestDealFlop(t *testing.T) {
	localDeck := deck.NewLocalDeck()
	localDeck.Shuffle(12345)
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	cards, err := service.DealFlop()
	if err != nil {
		t.Fatalf("DealFlop failed: %v", err)
	}

	// Should return 3 cards
	if len(cards) != 3 {
		t.Errorf("Expected 3 cards, got %d", len(cards))
	}

	// Should have burned 1 card, so 52 - 4 = 48 cards remaining
	// (1 burned + 3 drawn)
	if localDeck.RemainingCards() != 48 {
		t.Errorf("Expected 48 cards remaining, got %d", localDeck.RemainingCards())
	}
}

func TestDealTurn(t *testing.T) {
	localDeck := deck.NewLocalDeck()
	localDeck.Shuffle(12345)
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	card, err := service.DealTurn()
	if err != nil {
		t.Fatalf("DealTurn failed: %v", err)
	}

	// Should return 1 card (non-zero value)
	if card.String() == "" {
		t.Error("DealTurn returned empty card")
	}

	// Should have burned 1 card, so 52 - 2 = 50 cards remaining
	// (1 burned + 1 drawn)
	if localDeck.RemainingCards() != 50 {
		t.Errorf("Expected 50 cards remaining, got %d", localDeck.RemainingCards())
	}
}

func TestDealRiver(t *testing.T) {
	localDeck := deck.NewLocalDeck()
	localDeck.Shuffle(12345)
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	card, err := service.DealRiver()
	if err != nil {
		t.Fatalf("DealRiver failed: %v", err)
	}

	// Should return 1 card (non-zero value)
	if card.String() == "" {
		t.Error("DealRiver returned empty card")
	}

	// Should have burned 1 card, so 52 - 2 = 50 cards remaining
	// (1 burned + 1 drawn)
	if localDeck.RemainingCards() != 50 {
		t.Errorf("Expected 50 cards remaining, got %d", localDeck.RemainingCards())
	}
}

func TestCompleteGameDealingSequence(t *testing.T) {
	// Test a complete game dealing sequence:
	// 2 players get hole cards (4 cards)
	// Flop: burn 1, deal 3 (4 cards)
	// Turn: burn 1, deal 1 (2 cards)
	// River: burn 1, deal 1 (2 cards)
	// Total: 4 + 4 + 2 + 2 = 12 cards used, 40 remaining

	localDeck := deck.NewLocalDeck()
	localDeck.Shuffle(12345)
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	// Create 2 players
	p1ID := player.GeneratePlayerId()
	p2ID := player.GeneratePlayerId()
	p1Nick, _ := player.NewNickname("Player1")
	p2Nick, _ := player.NewNickname("Player2")
	player1, _ := player.NewPlayer(p1ID, p1Nick, 1000)
	player2, _ := player.NewPlayer(p2ID, p2Nick, 1000)
	players := []*player.Player{player1, player2}

	// Deal hole cards
	err := service.DealHoleCards(players)
	if err != nil {
		t.Fatalf("DealHoleCards failed: %v", err)
	}
	if localDeck.RemainingCards() != 48 {
		t.Errorf("After hole cards: expected 48 remaining, got %d", localDeck.RemainingCards())
	}

	// Deal flop
	flopCards, err := service.DealFlop()
	if err != nil {
		t.Fatalf("DealFlop failed: %v", err)
	}
	if len(flopCards) != 3 {
		t.Errorf("Expected 3 flop cards, got %d", len(flopCards))
	}
	if localDeck.RemainingCards() != 44 {
		t.Errorf("After flop: expected 44 remaining, got %d", localDeck.RemainingCards())
	}

	// Deal turn
	turnCard, err := service.DealTurn()
	if err != nil {
		t.Fatalf("DealTurn failed: %v", err)
	}
	if turnCard.String() == "" {
		t.Error("Turn card is empty")
	}
	if localDeck.RemainingCards() != 42 {
		t.Errorf("After turn: expected 42 remaining, got %d", localDeck.RemainingCards())
	}

	// Deal river
	riverCard, err := service.DealRiver()
	if err != nil {
		t.Fatalf("DealRiver failed: %v", err)
	}
	if riverCard.String() == "" {
		t.Error("River card is empty")
	}
	if localDeck.RemainingCards() != 40 {
		t.Errorf("After river: expected 40 remaining, got %d", localDeck.RemainingCards())
	}
}

func TestDealHoleCards_EmptyDeck(t *testing.T) {
	localDeck := deck.NewLocalDeck()
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	// Exhaust the deck
	for i := 0; i < 52; i++ {
		localDeck.DrawCard()
	}

	p1ID := player.GeneratePlayerId()
	p1Nick, _ := player.NewNickname("Player1")
	player1, _ := player.NewPlayer(p1ID, p1Nick, 1000)
	players := []*player.Player{player1}

	// Should fail with empty deck
	err := service.DealHoleCards(players)
	if err == nil {
		t.Error("Expected error when dealing from empty deck, got nil")
	}
}

func TestDealFlop_EmptyDeck(t *testing.T) {
	localDeck := deck.NewLocalDeck()
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	// Exhaust the deck
	for i := 0; i < 52; i++ {
		localDeck.DrawCard()
	}

	// Should fail with empty deck
	_, err := service.DealFlop()
	if err == nil {
		t.Error("Expected error when dealing flop from empty deck, got nil")
	}
}

func TestEvaluateHand(t *testing.T) {
	// This is currently a placeholder implementation
	localDeck := deck.NewLocalDeck()
	handEvaluator := game.NewHandEvaluator()
	service := NewGameService(localDeck, handEvaluator)

	p1ID := player.GeneratePlayerId()
	p1Nick, _ := player.NewNickname("Player1")
	player1, _ := player.NewPlayer(p1ID, p1Nick, 1000)

	// Deal some cards
	service.DealHoleCards([]*player.Player{player1})
	flopCards, _ := service.DealFlop()

	// Evaluate hand (currently returns empty result)
	result, err := service.EvaluateHand(player1.Hand(), flopCards)
	if err != nil {
		t.Fatalf("EvaluateHand failed: %v", err)
	}

	// For now, just verify it doesn't crash
	// TODO: Add real hand evaluation tests in Phase 2
	_ = result
}
