// Package deck provides deck adapter implementations
// Mirror of: pokerhole-server/adapter/out/deck/StandardDeck.java
package deck

import (
	"errors"
	"math/rand"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

// LocalDeck is an offline deck implementation (Adapter)
// Implements: card.DeckPort
type LocalDeck struct {
	cards []card.Card
	rng   *rand.Rand
}

// NewLocalDeck creates a new LocalDeck with 52 cards
func NewLocalDeck() *LocalDeck {
	deck := &LocalDeck{
		cards: make([]card.Card, 0, 52),
	}
	deck.Reset()
	return deck
}

// Compile-time check: LocalDeck implements card.DeckPort
var _ card.DeckPort = (*LocalDeck)(nil)

// DrawCard draws a single card from the deck
func (d *LocalDeck) DrawCard() (card.Card, error) {
	if len(d.cards) == 0 {
		return card.Card{}, errors.New("deck is empty")
	}

	// Draw from top (first card)
	drawnCard := d.cards[0]
	d.cards = d.cards[1:]

	return drawnCard, nil
}

// Shuffle shuffles the deck using Fisher-Yates algorithm (deterministic)
// IMPORTANT: Must use same algorithm as Java server for Golden Tests!
func (d *LocalDeck) Shuffle(seed int64) error {
	d.rng = rand.New(rand.NewSource(seed))

	// Fisher-Yates shuffle
	n := len(d.cards)
	for i := n - 1; i > 0; i-- {
		j := d.rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}

	return nil
}

// RemainingCards returns the number of cards left
func (d *LocalDeck) RemainingCards() int {
	return len(d.cards)
}

// Reset resets the deck to 52 cards
func (d *LocalDeck) Reset() error {
	d.cards = make([]card.Card, 0, 52)

	// Create all 52 cards (4 suits x 13 ranks)
	suits := []card.Suit{card.Clubs, card.Diamonds, card.Hearts, card.Spades}
	ranks := []card.Rank{
		card.Two, card.Three, card.Four, card.Five, card.Six, card.Seven,
		card.Eight, card.Nine, card.Ten, card.Jack, card.Queen, card.King, card.Ace,
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			c, _ := card.NewCard(suit, rank)
			d.cards = append(d.cards, c)
		}
	}

	return nil
}
