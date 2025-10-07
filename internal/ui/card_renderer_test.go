package ui

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

func TestParseCardString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantRank card.Rank
		wantSuit card.Suit
		wantNil  bool
	}{
		// Rank + Suit format (legacy)
		{"Ace of Spades (rank+suit)", "A♠", card.Ace, card.Spades, false},
		{"King of Hearts (rank+suit)", "K♥", card.King, card.Hearts, false},
		{"Ten of Spades (rank+suit)", "10♠", card.Ten, card.Spades, false},
		// Suit + Rank format (current Card.String() format)
		{"Ace of Spades (suit+rank)", "♠A", card.Ace, card.Spades, false},
		{"King of Hearts (suit+rank)", "♥K", card.King, card.Hearts, false},
		{"Queen of Diamonds (suit+rank)", "♦Q", card.Queen, card.Diamonds, false},
		{"Jack of Clubs (suit+rank)", "♣J", card.Jack, card.Clubs, false},
		{"Ten of Spades (suit+rank)", "♠10", card.Ten, card.Spades, false},
		{"Nine of Hearts (suit+rank)", "♥9", card.Nine, card.Hearts, false},
		{"Two of Clubs (suit+rank)", "♣2", card.Two, card.Clubs, false},
		{"Empty string", "", card.Ace, card.Spades, true},
		{"Too short", "A", card.Ace, card.Spades, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCardString(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("parseCardString(%q) = %v, want nil", tt.input, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("parseCardString(%q) = nil, want card", tt.input)
			}
			if got.Rank() != tt.wantRank {
				t.Errorf("parseCardString(%q).Rank() = %v, want %v", tt.input, got.Rank(), tt.wantRank)
			}
			if got.Suit() != tt.wantSuit {
				t.Errorf("parseCardString(%q).Suit() = %v, want %v", tt.input, got.Suit(), tt.wantSuit)
			}
		})
	}
}

func TestParseHand(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
	}{
		{"Two cards", "[A♠ K♥]", 2},
		{"Empty hand", "[]", 0},
		{"No cards string", "No cards", 0},
		{"Empty string", "", 0},
		{"Single card", "[A♠]", 1},
		{"Three cards", "[A♠ K♥ Q♦]", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseHand(tt.input)
			if len(got) != tt.wantCount {
				t.Errorf("parseHand(%q) returned %d cards, want %d", tt.input, len(got), tt.wantCount)
			}
		})
	}
}

func TestParseHandWithRealCards(t *testing.T) {
	input := "[A♠ K♥]"
	cards := parseHand(input)

	if len(cards) != 2 {
		t.Fatalf("parseHand(%q) returned %d cards, want 2", input, len(cards))
	}

	// Check first card (Ace of Spades)
	if cards[0].Rank() != card.Ace {
		t.Errorf("First card rank = %v, want Ace", cards[0].Rank())
	}
	if cards[0].Suit() != card.Spades {
		t.Errorf("First card suit = %v, want Spades", cards[0].Suit())
	}

	// Check second card (King of Hearts)
	if cards[1].Rank() != card.King {
		t.Errorf("Second card rank = %v, want King", cards[1].Rank())
	}
	if cards[1].Suit() != card.Hearts {
		t.Errorf("Second card suit = %v, want Hearts", cards[1].Suit())
	}
}
