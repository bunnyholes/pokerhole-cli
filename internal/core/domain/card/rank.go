// Package card provides core card domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/card/Rank.java
package card

// Rank represents the rank of a card (2-10, J, Q, K, A)
type Rank int

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var rankSymbols = [...]string{
	"2", "3", "4", "5", "6", "7", "8", "9", "10",
	"J", "Q", "K", "A",
}

// String returns the string representation of the rank
func (r Rank) String() string {
	if r < Two || r > Ace {
		return "?"
	}
	return rankSymbols[r]
}

// Value returns the numeric value of the rank (2-14, Ace=14)
func (r Rank) Value() int {
	return int(r) + 2 // Two=0 → 2, Three=1 → 3, ..., Ace=12 → 14
}
