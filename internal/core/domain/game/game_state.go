// Package game provides game domain logic
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/GameState.java
package game

// GameState represents the overall state of a game
type GameState int

const (
	Waiting  GameState = iota // Waiting for players
	Playing                   // Game in progress
	Finished                  // Game finished
)

var gameStateNames = [...]string{
	"WAITING",
	"PLAYING",
	"FINISHED",
}

// String returns the string representation
func (g GameState) String() string {
	if g < Waiting || g > Finished {
		return "UNKNOWN"
	}
	return gameStateNames[g]
}

// IsPlaying returns true if game is in progress
func (g GameState) IsPlaying() bool {
	return g == Playing
}
