// Package vo provides game value objects
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/vo/SidePot.java
package vo

import "github.com/bunnyholes/pokerhole/client/internal/core/domain/player"

// SidePot represents a side pot created when a player goes all-in
type SidePot struct {
	amount             int
	eligiblePlayerIDs  []player.PlayerId
	capPerPlayer       int
}

// NewSidePot creates a new side pot
func NewSidePot(amount int, eligiblePlayerIDs []player.PlayerId, capPerPlayer int) SidePot {
	// TODO: validate inputs
	return SidePot{
		amount:            amount,
		eligiblePlayerIDs: eligiblePlayerIDs,
		capPerPlayer:      capPerPlayer,
	}
}

// Amount returns the pot amount
func (s SidePot) Amount() int {
	return s.amount
}

// EligiblePlayerIDs returns IDs of players eligible to win this pot
func (s SidePot) EligiblePlayerIDs() []player.PlayerId {
	// TODO: return defensive copy
	return nil
}

// CapPerPlayer returns the max contribution per player
func (s SidePot) CapPerPlayer() int {
	return s.capPerPlayer
}

// IsPlayerEligible checks if a player is eligible for this pot
func (s SidePot) IsPlayerEligible(playerID player.PlayerId) bool {
	// TODO: implement
	return false
}

// String returns string representation
func (s SidePot) String() string {
	// TODO: implement
	return ""
}
