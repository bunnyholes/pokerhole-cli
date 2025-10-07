// Package player provides player domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/player/vo/PlayerStatus.java
package player

// PlayerStatus represents the current status of a player
type PlayerStatus int

const (
	Active  PlayerStatus = iota // Player is actively playing
	Folded                      // Player has folded this round
	AllIn                       // Player is all-in
	SitOut                      // Player is sitting out
	Waiting                     // Player is waiting for next round
)

var playerStatusNames = [...]string{
	"ACTIVE",
	"FOLDED",
	"ALL_IN",
	"SIT_OUT",
	"WAITING",
}

// String returns the string representation
func (s PlayerStatus) String() string {
	if s < Active || s > Waiting {
		return "UNKNOWN"
	}
	return playerStatusNames[s]
}

// IsActive returns true if player is actively playing
func (s PlayerStatus) IsActive() bool {
	return s == Active
}

// CanAct returns true if player can take an action
func (s PlayerStatus) CanAct() bool {
	// TODO: implement (Active players can act)
	return false
}
