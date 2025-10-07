// Package vo provides game value objects
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/vo/Pot.java
package vo

// Pot represents the main pot in a poker game
type Pot struct {
	amount int
}

// NewPot creates a new pot with the given amount
func NewPot(amount int) Pot {
	return Pot{amount: amount}
}

// Amount returns the pot amount
func (p Pot) Amount() int {
	return p.amount
}

// Add adds chips to the pot
func (p *Pot) Add(chips int) error {
	// TODO: validate chips > 0
	// TODO: implement
	return nil
}

// IsEmpty returns true if pot is empty
func (p Pot) IsEmpty() bool {
	return p.amount == 0
}

// String returns string representation
func (p Pot) String() string {
	// TODO: implement
	return ""
}
