// Package service provides application services
package service

import (
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

// GameService orchestrates game logic (Application Service)
type GameService struct {
	deck          card.DeckPort      // Port (interface)
	HandEvaluator game.HandEvaluator // Domain service (exported for access)
}

// NewGameService creates a new GameService
func NewGameService(deck card.DeckPort, handEvaluator game.HandEvaluator) *GameService {
	return &GameService{
		deck:          deck,
		HandEvaluator: handEvaluator,
	}
}

// DealHoleCards deals 2 cards to each player
func (s *GameService) DealHoleCards(players []*player.Player) error {
	for _, p := range players {
		// Draw 2 cards
		card1, err := s.deck.DrawCard()
		if err != nil {
			return err
		}
		card2, err := s.deck.DrawCard()
		if err != nil {
			return err
		}

		// Create hand and set to player
		hand := card.NewHand([]card.Card{card1, card2})
		p.SetHand(hand)
	}
	return nil
}

// DealFlop deals 3 community cards
func (s *GameService) DealFlop() ([]card.Card, error) {
	// Burn 1 card
	_, err := s.deck.DrawCard()
	if err != nil {
		return nil, err
	}

	// Draw 3 cards
	cards := make([]card.Card, 3)
	for i := 0; i < 3; i++ {
		c, err := s.deck.DrawCard()
		if err != nil {
			return nil, err
		}
		cards[i] = c
	}

	return cards, nil
}

// DealTurn deals 1 turn card
func (s *GameService) DealTurn() (card.Card, error) {
	// Burn 1 card
	_, err := s.deck.DrawCard()
	if err != nil {
		return card.Card{}, err
	}

	// Draw 1 card
	return s.deck.DrawCard()
}

// DealRiver deals 1 river card
func (s *GameService) DealRiver() (card.Card, error) {
	// Burn 1 card
	_, err := s.deck.DrawCard()
	if err != nil {
		return card.Card{}, err
	}

	// Draw 1 card
	return s.deck.DrawCard()
}

// EvaluateHand evaluates a player's hand
func (s *GameService) EvaluateHand(playerHand card.Hand, communityCards []card.Card) (vo.HandResult, error) {
	// For now, return empty result (hand evaluator not implemented yet)
	// This will be implemented in Phase 2
	return vo.HandResult{}, nil
}
