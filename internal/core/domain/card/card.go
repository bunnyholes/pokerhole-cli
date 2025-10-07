// Package card provides core card domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/card/Card.java
package card

import "fmt"

// Card represents a single playing card (immutable value object)
type Card struct {
	suit Suit
	rank Rank
}

// NewCard creates a new Card
func NewCard(suit Suit, rank Rank) (Card, error) {
	// TODO: validate suit and rank
	return Card{suit: suit, rank: rank}, nil
}

// Suit returns the suit of the card
func (c Card) Suit() Suit {
	return c.suit
}

// Rank returns the rank of the card
func (c Card) Rank() Rank {
	return c.rank
}

// String returns the string representation (e.g., "♠A")
func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.suit, c.rank)
}

// CompareTo compares two cards by rank, then by suit
// Returns: -1 if c < other, 0 if equal, 1 if c > other
func (c Card) CompareTo(other Card) int {
	// Compare by rank first
	if c.rank.Value() < other.rank.Value() {
		return -1
	}
	if c.rank.Value() > other.rank.Value() {
		return 1
	}
	// Same rank, compare by suit
	if c.suit < other.suit {
		return -1
	}
	if c.suit > other.suit {
		return 1
	}
	return 0
}

// Equals checks if two cards are equal
func (c Card) Equals(other Card) bool {
	return c.suit == other.suit && c.rank == other.rank
}
