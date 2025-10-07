// Package vo provides game value objects
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/vo/PlayerAction.java
package vo

// PlayerAction represents an action a player can take
type PlayerAction int

const (
	Fold   PlayerAction = iota // Forfeit hand
	Check                      // Pass (only if currentBet == 0)
	Call                       // Match current bet
	Raise                      // Bet higher than current
	AllIn                      // Bet all remaining chips
)

var playerActionNames = [...]string{
	"FOLD",
	"CHECK",
	"CALL",
	"RAISE",
	"ALL_IN",
}

// String returns the string representation
func (a PlayerAction) String() string {
	if a < Fold || a > AllIn {
		return "UNKNOWN"
	}
	return playerActionNames[a]
}

// IsValid checks if the action is valid
func (a PlayerAction) IsValid() bool {
	return a >= Fold && a <= AllIn
}

// RequiresBetAmount returns true if the action requires a bet amount
func (a PlayerAction) RequiresBetAmount() bool {
	// TODO: implement (RAISE and ALL_IN require amount)
	return false
}
