package deck

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

func TestNewLocalDeck(t *testing.T) {
	deck := NewLocalDeck()

	if deck.RemainingCards() != 52 {
		t.Errorf("Expected 52 cards, got %d", deck.RemainingCards())
	}
}

func TestDrawCard(t *testing.T) {
	deck := NewLocalDeck()

	// Draw first card
	card1, err := deck.DrawCard()
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	if deck.RemainingCards() != 51 {
		t.Errorf("Expected 51 cards remaining, got %d", deck.RemainingCards())
	}

	// Draw second card
	card2, err := deck.DrawCard()
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	// Cards should be different
	if card1.Equals(card2) {
		t.Error("Drew the same card twice")
	}

	if deck.RemainingCards() != 50 {
		t.Errorf("Expected 50 cards remaining, got %d", deck.RemainingCards())
	}
}

func TestDrawCardEmptyDeck(t *testing.T) {
	deck := NewLocalDeck()

	// Draw all 52 cards
	for i := 0; i < 52; i++ {
		_, err := deck.DrawCard()
		if err != nil {
			t.Fatalf("Failed to draw card %d: %v", i, err)
		}
	}

	// Try to draw from empty deck
	_, err := deck.DrawCard()
	if err == nil {
		t.Error("Expected error when drawing from empty deck")
	}
}

func TestShuffle(t *testing.T) {
	deck1 := NewLocalDeck()
	deck2 := NewLocalDeck()

	// Get first 5 cards from unshuffled deck1
	cards1Before := make([]card.Card, 5)
	for i := 0; i < 5; i++ {
		cards1Before[i], _ = deck1.DrawCard()
	}

	// Reset and shuffle deck1
	deck1.Reset()
	deck1.Shuffle(42)

	// Get first 5 cards from shuffled deck1
	cards1After := make([]card.Card, 5)
	for i := 0; i < 5; i++ {
		cards1After[i], _ = deck1.DrawCard()
	}

	// They should be different (very unlikely to be the same after shuffle)
	allSame := true
	for i := 0; i < 5; i++ {
		if !cards1Before[i].Equals(cards1After[i]) {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Shuffle did not change card order")
	}

	// Shuffle deck2 with same seed
	deck2.Shuffle(42)

	// Get first 5 cards from deck2
	cards2 := make([]card.Card, 5)
	for i := 0; i < 5; i++ {
		cards2[i], _ = deck2.DrawCard()
	}

	// deck1 and deck2 should have same cards (deterministic shuffle)
	for i := 0; i < 5; i++ {
		if !cards1After[i].Equals(cards2[i]) {
			t.Errorf("Deterministic shuffle failed: card %d different", i)
		}
	}
}

func TestReset(t *testing.T) {
	deck := NewLocalDeck()

	// Draw 10 cards
	for i := 0; i < 10; i++ {
		deck.DrawCard()
	}

	if deck.RemainingCards() != 42 {
		t.Errorf("Expected 42 cards, got %d", deck.RemainingCards())
	}

	// Reset
	deck.Reset()

	if deck.RemainingCards() != 52 {
		t.Errorf("Expected 52 cards after reset, got %d", deck.RemainingCards())
	}
}
