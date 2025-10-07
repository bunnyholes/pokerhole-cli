// Package game provides game domain logic
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/HandEvaluator.java
package game

import (
	"errors"
	"sort"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
)

// HandEvaluator evaluates poker hands (Domain Service)
type HandEvaluator interface {
	// Evaluate evaluates the best 5-card hand from player's hole cards + community cards
	Evaluate(playerCards []card.Card, communityCards []card.Card) (vo.HandResult, error)

	// EvaluateFiveCards evaluates exactly 5 cards
	EvaluateFiveCards(fiveCards []card.Card) (vo.HandResult, error)
}

// handEvaluatorImpl is the implementation of HandEvaluator
type handEvaluatorImpl struct {
	// No state (stateless domain service)
}

// NewHandEvaluator creates a new HandEvaluator
func NewHandEvaluator() HandEvaluator {
	return &handEvaluatorImpl{}
}

// Evaluate evaluates the best poker hand
func (h *handEvaluatorImpl) Evaluate(playerCards []card.Card, communityCards []card.Card) (vo.HandResult, error) {
	// Combine all cards
	allCards := append(playerCards, communityCards...)

	if len(allCards) < 5 {
		return vo.HandResult{}, errors.New("insufficient cards for evaluation")
	}

	// If exactly 5 cards, evaluate directly
	if len(allCards) == 5 {
		return h.EvaluateFiveCards(allCards)
	}

	// Generate all 5-card combinations and find the best
	var bestResult vo.HandResult
	combinations := generateCombinations(allCards, 5)

	for i, combo := range combinations {
		result, err := h.EvaluateFiveCards(combo)
		if err != nil {
			continue
		}

		if i == 0 || result.CompareTo(bestResult) > 0 {
			bestResult = result
		}
	}

	return bestResult, nil
}

// EvaluateFiveCards evaluates exactly 5 cards
func (h *handEvaluatorImpl) EvaluateFiveCards(fiveCards []card.Card) (vo.HandResult, error) {
	if len(fiveCards) != 5 {
		return vo.HandResult{}, errors.New("exactly 5 cards required")
	}

	// Sort cards by rank (descending)
	sorted := make([]card.Card, len(fiveCards))
	copy(sorted, fiveCards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rank().Value() > sorted[j].Rank().Value()
	})

	// Check hands from highest to lowest rank
	if result, ok := h.checkRoyalFlush(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkStraightFlush(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkFourOfAKind(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkFullHouse(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkFlush(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkStraight(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkThreeOfAKind(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkTwoPair(sorted); ok {
		return result, nil
	}
	if result, ok := h.checkOnePair(sorted); ok {
		return result, nil
	}

	// High card
	return h.checkHighCard(sorted), nil
}

// Helper functions (private)

func (h *handEvaluatorImpl) checkRoyalFlush(cards []card.Card) (vo.HandResult, bool) {
	// Royal Flush: A-K-Q-J-10 same suit
	if !h.isFlush(cards) {
		return vo.HandResult{}, false
	}

	ranks := []int{cards[0].Rank().Value(), cards[1].Rank().Value(), cards[2].Rank().Value(),
		cards[3].Rank().Value(), cards[4].Rank().Value()}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

	if ranks[0] == 14 && ranks[1] == 13 && ranks[2] == 12 && ranks[3] == 11 && ranks[4] == 10 {
		return vo.NewHandResult(vo.RoyalFlush, cards, ranks), true
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkStraightFlush(cards []card.Card) (vo.HandResult, bool) {
	if h.isFlush(cards) && h.isStraight(cards) {
		// Check for Ace-low straight (wheel): A-2-3-4-5
		ranks := make([]int, len(cards))
		for i, c := range cards {
			ranks[i] = c.Rank().Value()
		}
		sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

		var tieBreaker []int
		if ranks[0] == 14 && ranks[1] == 5 && ranks[2] == 4 && ranks[3] == 3 && ranks[4] == 2 {
			// Ace-low straight: treat as 5-high
			tieBreaker = []int{5}
		} else {
			// Normal straight: use highest card
			tieBreaker = []int{cards[0].Rank().Value()}
		}
		return vo.NewHandResult(vo.StraightFlush, cards, tieBreaker), true
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkFourOfAKind(cards []card.Card) (vo.HandResult, bool) {
	rankCounts := h.countRanks(cards)

	for rank, count := range rankCounts {
		if count == 4 {
			// Find kicker
			kicker := 0
			for r, c := range rankCounts {
				if c == 1 {
					kicker = r
				}
			}
			tieBreaker := []int{rank, kicker}
			return vo.NewHandResult(vo.FourOfAKind, cards, tieBreaker), true
		}
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkFullHouse(cards []card.Card) (vo.HandResult, bool) {
	rankCounts := h.countRanks(cards)

	threeRank := -1
	pairRank := -1

	for rank, count := range rankCounts {
		if count == 3 {
			threeRank = rank
		} else if count == 2 {
			pairRank = rank
		}
	}

	if threeRank != -1 && pairRank != -1 {
		tieBreaker := []int{threeRank, pairRank}
		return vo.NewHandResult(vo.FullHouse, cards, tieBreaker), true
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkFlush(cards []card.Card) (vo.HandResult, bool) {
	if h.isFlush(cards) {
		// All 5 cards as tiebreaker (descending)
		tieBreaker := make([]int, 5)
		for i, c := range cards {
			tieBreaker[i] = c.Rank().Value()
		}
		return vo.NewHandResult(vo.Flush, cards, tieBreaker), true
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkStraight(cards []card.Card) (vo.HandResult, bool) {
	if h.isStraight(cards) {
		// Check for Ace-low straight (wheel): A-2-3-4-5
		ranks := make([]int, len(cards))
		for i, c := range cards {
			ranks[i] = c.Rank().Value()
		}
		sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

		var tieBreaker []int
		if ranks[0] == 14 && ranks[1] == 5 && ranks[2] == 4 && ranks[3] == 3 && ranks[4] == 2 {
			// Ace-low straight: treat as 5-high
			tieBreaker = []int{5}
		} else {
			// Normal straight: use highest card
			tieBreaker = []int{cards[0].Rank().Value()}
		}
		return vo.NewHandResult(vo.Straight, cards, tieBreaker), true
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkThreeOfAKind(cards []card.Card) (vo.HandResult, bool) {
	rankCounts := h.countRanks(cards)

	for rank, count := range rankCounts {
		if count == 3 {
			// Find kickers
			kickers := []int{}
			for r, c := range rankCounts {
				if c == 1 {
					kickers = append(kickers, r)
				}
			}
			sort.Sort(sort.Reverse(sort.IntSlice(kickers)))
			tieBreaker := append([]int{rank}, kickers...)
			return vo.NewHandResult(vo.ThreeOfAKind, cards, tieBreaker), true
		}
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkTwoPair(cards []card.Card) (vo.HandResult, bool) {
	rankCounts := h.countRanks(cards)

	pairs := []int{}
	kicker := 0

	for rank, count := range rankCounts {
		if count == 2 {
			pairs = append(pairs, rank)
		} else if count == 1 {
			kicker = rank
		}
	}

	if len(pairs) == 2 {
		sort.Sort(sort.Reverse(sort.IntSlice(pairs)))
		tieBreaker := []int{pairs[0], pairs[1], kicker}
		return vo.NewHandResult(vo.TwoPair, cards, tieBreaker), true
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkOnePair(cards []card.Card) (vo.HandResult, bool) {
	rankCounts := h.countRanks(cards)

	for rank, count := range rankCounts {
		if count == 2 {
			// Find kickers
			kickers := []int{}
			for r, c := range rankCounts {
				if c == 1 {
					kickers = append(kickers, r)
				}
			}
			sort.Sort(sort.Reverse(sort.IntSlice(kickers)))
			tieBreaker := append([]int{rank}, kickers...)
			return vo.NewHandResult(vo.OnePair, cards, tieBreaker), true
		}
	}
	return vo.HandResult{}, false
}

func (h *handEvaluatorImpl) checkHighCard(cards []card.Card) vo.HandResult {
	// All cards as tiebreaker (already sorted descending)
	tieBreaker := make([]int, 5)
	for i, c := range cards {
		tieBreaker[i] = c.Rank().Value()
	}
	return vo.NewHandResult(vo.HighCard, cards, tieBreaker)
}

// Utility functions

func (h *handEvaluatorImpl) isFlush(cards []card.Card) bool {
	suit := cards[0].Suit()
	for _, c := range cards[1:] {
		if c.Suit() != suit {
			return false
		}
	}
	return true
}

func (h *handEvaluatorImpl) isStraight(cards []card.Card) bool {
	// Get ranks in descending order
	ranks := make([]int, len(cards))
	for i, c := range cards {
		ranks[i] = c.Rank().Value()
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

	// Check for regular straight
	for i := 0; i < len(ranks)-1; i++ {
		if ranks[i]-ranks[i+1] != 1 {
			// Check for A-2-3-4-5 (wheel)
			if ranks[0] == 14 && ranks[1] == 5 && ranks[2] == 4 && ranks[3] == 3 && ranks[4] == 2 {
				return true
			}
			return false
		}
	}
	return true
}

func (h *handEvaluatorImpl) countRanks(cards []card.Card) map[int]int {
	counts := make(map[int]int)
	for _, c := range cards {
		counts[c.Rank().Value()]++
	}
	return counts
}

// generateCombinations generates all k-combinations from a slice
func generateCombinations(cards []card.Card, k int) [][]card.Card {
	var result [][]card.Card
	n := len(cards)

	var generate func(start int, combo []card.Card)
	generate = func(start int, combo []card.Card) {
		if len(combo) == k {
			temp := make([]card.Card, k)
			copy(temp, combo)
			result = append(result, temp)
			return
		}

		for i := start; i < n; i++ {
			generate(i+1, append(combo, cards[i]))
		}
	}

	generate(0, []card.Card{})
	return result
}
