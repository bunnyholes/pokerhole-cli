// Package game provides game domain logic
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/PotDistributor.java
package game

import (
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

// PotDistributor handles pot distribution logic (Domain Service)
type PotDistributor struct {
	// Stateless
}

// NewPotDistributor creates a new PotDistributor
func NewPotDistributor() *PotDistributor {
	return &PotDistributor{}
}

// DistributePot distributes the main pot to winners
func (p *PotDistributor) DistributePot(pot vo.Pot, winners []*player.Player) map[player.PlayerId]int {
	// TODO: implement
	// 1. Split pot amount evenly among winners
	// 2. Handle remainder chips (give to player closest to dealer button)
	// 3. Return map of playerId -> chips won
	return nil
}

// DistributeSidePots distributes multiple side pots
func (p *PotDistributor) DistributeSidePots(sidePots []vo.SidePot, players []*player.Player, handResults map[player.PlayerId]vo.HandResult) map[player.PlayerId]int {
	// TODO: implement
	// 1. For each side pot:
	//    a. Find eligible players
	//    b. Determine winner(s) among eligible players
	//    c. Distribute pot
	// 2. Return total winnings per player
	return nil
}

// CreateSidePots creates side pots when players go all-in
func (p *PotDistributor) CreateSidePots(players []*player.Player) []vo.SidePot {
	// TODO: implement
	// 1. Sort players by bet amount (ascending)
	// 2. For each all-in player:
	//    a. Create side pot up to their bet amount
	//    b. Mark eligible players
	// 3. Create main pot with remaining
	return nil
}
