// Package game provides game domain logic
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/WinnerResolver.java
package game

import (
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

// WinnerResolver determines the winner(s) of a poker hand (Domain Service)
type WinnerResolver struct {
	handEvaluator HandEvaluator
}

// NewWinnerResolver creates a new WinnerResolver
func NewWinnerResolver(handEvaluator HandEvaluator) *WinnerResolver {
	return &WinnerResolver{
		handEvaluator: handEvaluator,
	}
}

// DetermineWinners determines the winner(s) from a list of players
func (w *WinnerResolver) DetermineWinners(players []*player.Player, communityCards []card.Card) ([]*player.Player, error) {
	// Filter out folded players
	activePlayers := []*player.Player{}
	for _, p := range players {
		if p.Status() != player.Folded {
			activePlayers = append(activePlayers, p)
		}
	}

	if len(activePlayers) == 0 {
		return nil, nil
	}

	// If only one active player, they win by default
	if len(activePlayers) == 1 {
		return activePlayers, nil
	}

	// Evaluate each player's hand
	handResults := make(map[*player.Player]vo.HandResult)
	for _, p := range activePlayers {
		playerCards := p.Hand().Cards()
		result, err := w.handEvaluator.Evaluate(playerCards, communityCards)
		if err != nil {
			return nil, err
		}
		handResults[p] = result
	}

	// Find the best hand
	var bestHand vo.HandResult
	var winners []*player.Player

	for _, p := range activePlayers {
		hand := handResults[p]

		if len(winners) == 0 {
			// First player
			bestHand = hand
			winners = []*player.Player{p}
		} else {
			comparison := hand.CompareTo(bestHand)
			if comparison > 0 {
				// New best hand
				bestHand = hand
				winners = []*player.Player{p}
			} else if comparison == 0 {
				// Tie
				winners = append(winners, p)
			}
		}
	}

	return winners, nil
}

// CompareHands compares two hands and returns the winner
// Returns: -1 if hand1 wins, 0 if tie, 1 if hand2 wins
func (w *WinnerResolver) CompareHands(hand1 vo.HandResult, hand2 vo.HandResult) int {
	// hand1.CompareTo returns: -1 if hand1 < hand2, 0 if equal, 1 if hand1 > hand2
	// We need to return: -1 if hand1 wins, 0 if tie, 1 if hand2 wins
	// So invert the result
	return -hand1.CompareTo(hand2)
}
