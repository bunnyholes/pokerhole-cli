// Package card provides core card domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/card/Deck.java
package card

// DeckPort defines the interface for a card deck (Port for hexagonal architecture)
// Implementations: LocalDeck (offline), RemoteDeck (online)
type DeckPort interface {
	// DrawCard draws a single card from the deck
	DrawCard() (Card, error)

	// Shuffle shuffles the deck with the given seed (deterministic)
	Shuffle(seed int64) error

	// RemainingCards returns the number of cards left in the deck
	RemainingCards() int

	// Reset resets the deck to 52 cards
	Reset() error
}
