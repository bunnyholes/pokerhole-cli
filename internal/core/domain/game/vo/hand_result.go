// Package vo provides game value objects
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/vo/HandResult.java
package vo

import "github.com/bunnyholes/pokerhole/client/internal/core/domain/card"

// Tier represents the ranking tier of a poker hand
type Tier int

const (
	HighCard Tier = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// HandResult represents the result of evaluating a poker hand
type HandResult struct {
	tier       Tier
	bestCards  []card.Card // Best 5 cards
	tieBreaker []int       // Tie-breaking values
}

// NewHandResult creates a new HandResult
func NewHandResult(tier Tier, bestCards []card.Card, tieBreaker []int) HandResult {
	// TODO: validate bestCards length (must be 5)
	return HandResult{
		tier:       tier,
		bestCards:  bestCards,
		tieBreaker: tieBreaker,
	}
}

// Tier returns the hand tier
func (h HandResult) Tier() Tier {
	return h.tier
}

// BestCards returns the best 5 cards
func (h HandResult) BestCards() []card.Card {
	// Return defensive copy
	cards := make([]card.Card, len(h.bestCards))
	copy(cards, h.bestCards)
	return cards
}

// TieBreaker returns tie-breaking values
func (h HandResult) TieBreaker() []int {
	// Return defensive copy
	values := make([]int, len(h.tieBreaker))
	copy(values, h.tieBreaker)
	return values
}

// CompareTo compares two hand results
// Returns: -1 if h < other, 0 if equal, 1 if h > other
func (h HandResult) CompareTo(other HandResult) int {
	// First compare tiers
	if h.tier < other.tier {
		return -1
	}
	if h.tier > other.tier {
		return 1
	}

	// Same tier, compare tiebreakers
	minLen := len(h.tieBreaker)
	if len(other.tieBreaker) < minLen {
		minLen = len(other.tieBreaker)
	}

	for i := 0; i < minLen; i++ {
		if h.tieBreaker[i] < other.tieBreaker[i] {
			return -1
		}
		if h.tieBreaker[i] > other.tieBreaker[i] {
			return 1
		}
	}

	return 0 // Equal
}

// GetRankCards returns only the cards that make up the hand rank
// For example, Two Pair returns only the 4 cards (2 pairs), not the kicker
func (h HandResult) GetRankCards() []card.Card {
	switch h.tier {
	case HighCard:
		// Return only the highest card
		if len(h.bestCards) > 0 {
			return []card.Card{h.bestCards[0]}
		}
		return []card.Card{}

	case OnePair:
		// Return only the pair (2 cards with same rank)
		if len(h.tieBreaker) == 0 {
			return []card.Card{}
		}
		pairRank := h.tieBreaker[0]
		return filterCardsByRank(h.bestCards, pairRank)

	case TwoPair:
		// Return both pairs (4 cards)
		if len(h.tieBreaker) < 2 {
			return []card.Card{}
		}
		highPairRank := h.tieBreaker[0]
		lowPairRank := h.tieBreaker[1]
		result := []card.Card{}
		result = append(result, filterCardsByRank(h.bestCards, highPairRank)...)
		result = append(result, filterCardsByRank(h.bestCards, lowPairRank)...)
		return result

	case ThreeOfAKind:
		// Return only the three of a kind
		if len(h.tieBreaker) == 0 {
			return []card.Card{}
		}
		tripsRank := h.tieBreaker[0]
		return filterCardsByRank(h.bestCards, tripsRank)

	case FourOfAKind:
		// Return only the four of a kind
		if len(h.tieBreaker) == 0 {
			return []card.Card{}
		}
		quadsRank := h.tieBreaker[0]
		return filterCardsByRank(h.bestCards, quadsRank)

	case Straight, Flush, FullHouse, StraightFlush, RoyalFlush:
		// Return all 5 cards
		return h.BestCards()

	default:
		return h.BestCards()
	}
}

// filterCardsByRank filters cards by rank value
func filterCardsByRank(cards []card.Card, rankValue int) []card.Card {
	result := []card.Card{}
	for _, c := range cards {
		if c.Rank().Value() == rankValue {
			result = append(result, c)
		}
	}
	return result
}

// String returns string representation
func (h HandResult) String() string {
	tierNames := []string{
		"High Card",
		"One Pair",
		"Two Pair",
		"Three of a Kind",
		"Straight",
		"Flush",
		"Full House",
		"Four of a Kind",
		"Straight Flush",
		"Royal Flush",
	}

	if h.tier >= 0 && int(h.tier) < len(tierNames) {
		return tierNames[h.tier]
	}
	return "Unknown"
}
