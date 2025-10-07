// Package card provides core card domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/card/Suit.java
package card

// Suit represents the suit of a card (♠ ♥ ♦ ♣)
type Suit int

const (
	Clubs    Suit = iota // ♣
	Diamonds             // ♦
	Hearts               // ♥
	Spades               // ♠
)

var suitSymbols = [...]string{"♣", "♦", "♥", "♠"}

// String returns the symbol representation of the suit
func (s Suit) String() string {
	if s < Clubs || s > Spades {
		return "?"
	}
	return suitSymbols[s]
}

// IsRed returns true if the suit is red (Hearts or Diamonds)
func (s Suit) IsRed() bool {
	return s == Hearts || s == Diamonds
}

// IsBlack returns true if the suit is black (Clubs or Spades)
func (s Suit) IsBlack() bool {
	return s == Clubs || s == Spades
}
