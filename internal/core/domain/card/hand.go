// Package card provides core card domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/card/Hand.java
package card

// Hand represents a player's hand (e.g., 2 hole cards in Texas Hold'em)
type Hand struct {
	cards []Card
}

// NewHand creates a new hand with the given cards
func NewHand(cards []Card) Hand {
	// TODO: validate cards (e.g., no duplicates)
	return Hand{cards: cards}
}

// Cards returns a copy of the cards in the hand
func (h Hand) Cards() []Card {
	// Return defensive copy
	cardsCopy := make([]Card, len(h.cards))
	copy(cardsCopy, h.cards)
	return cardsCopy
}

// Size returns the number of cards in the hand
func (h Hand) Size() int {
	return len(h.cards)
}

// AddCard adds a card to the hand
func (h *Hand) AddCard(card Card) error {
	h.cards = append(h.cards, card)
	return nil
}

// String returns string representation of the hand
func (h Hand) String() string {
	if len(h.cards) == 0 {
		return "[]"
	}
	result := "["
	for i, card := range h.cards {
		if i > 0 {
			result += " "
		}
		result += card.String()
	}
	result += "]"
	return result
}
