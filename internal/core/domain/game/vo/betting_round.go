// Package vo provides game value objects
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/vo/BettingRound.java
package vo

// BettingRound represents the current betting round in Texas Hold'em
type BettingRound int

const (
	PreFlop BettingRound = iota // Before flop (2 hole cards dealt)
	Flop                        // After 3 community cards
	Turn                        // After 4th community card
	River                       // After 5th community card
	Showdown                    // Revealing hands
)

var bettingRoundNames = [...]string{
	"PRE_FLOP",
	"FLOP",
	"TURN",
	"RIVER",
	"SHOWDOWN",
}

// String returns the string representation
func (b BettingRound) String() string {
	if b < PreFlop || b > Showdown {
		return "UNKNOWN"
	}
	return bettingRoundNames[b]
}

// Next returns the next betting round
func (b BettingRound) Next() BettingRound {
	// TODO: implement
	return PreFlop
}

// IsValid checks if the betting round is valid
func (b BettingRound) IsValid() bool {
	return b >= PreFlop && b <= Showdown
}
