package game

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
)

func TestHandEvaluator_HighCard(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.King),
		mustMakeCard(card.Diamonds, card.Ten),
		mustMakeCard(card.Clubs, card.Seven),
		mustMakeCard(card.Spades, card.Three),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.HighCard {
		t.Errorf("Expected HighCard, got %v", result.Tier())
	}
}

func TestHandEvaluator_OnePair(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Ace),
		mustMakeCard(card.Diamonds, card.Ten),
		mustMakeCard(card.Clubs, card.Seven),
		mustMakeCard(card.Spades, card.Three),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.OnePair {
		t.Errorf("Expected OnePair, got %v", result.Tier())
	}
}

func TestHandEvaluator_TwoPair(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Ace),
		mustMakeCard(card.Diamonds, card.King),
		mustMakeCard(card.Clubs, card.King),
		mustMakeCard(card.Spades, card.Three),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.TwoPair {
		t.Errorf("Expected TwoPair, got %v", result.Tier())
	}
}

func TestHandEvaluator_ThreeOfAKind(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Ace),
		mustMakeCard(card.Diamonds, card.Ace),
		mustMakeCard(card.Clubs, card.King),
		mustMakeCard(card.Spades, card.Three),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.ThreeOfAKind {
		t.Errorf("Expected ThreeOfAKind, got %v", result.Tier())
	}
}

func TestHandEvaluator_Straight(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ten),
		mustMakeCard(card.Hearts, card.Nine),
		mustMakeCard(card.Diamonds, card.Eight),
		mustMakeCard(card.Clubs, card.Seven),
		mustMakeCard(card.Spades, card.Six),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.Straight {
		t.Errorf("Expected Straight, got %v", result.Tier())
	}
}

func TestHandEvaluator_Flush(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Spades, card.King),
		mustMakeCard(card.Spades, card.Ten),
		mustMakeCard(card.Spades, card.Seven),
		mustMakeCard(card.Spades, card.Three),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.Flush {
		t.Errorf("Expected Flush, got %v", result.Tier())
	}
}

func TestHandEvaluator_FullHouse(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Ace),
		mustMakeCard(card.Diamonds, card.Ace),
		mustMakeCard(card.Clubs, card.King),
		mustMakeCard(card.Spades, card.King),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.FullHouse {
		t.Errorf("Expected FullHouse, got %v", result.Tier())
	}
}

func TestHandEvaluator_FourOfAKind(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Ace),
		mustMakeCard(card.Diamonds, card.Ace),
		mustMakeCard(card.Clubs, card.Ace),
		mustMakeCard(card.Spades, card.King),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.FourOfAKind {
		t.Errorf("Expected FourOfAKind, got %v", result.Tier())
	}
}

func TestHandEvaluator_StraightFlush(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ten),
		mustMakeCard(card.Spades, card.Nine),
		mustMakeCard(card.Spades, card.Eight),
		mustMakeCard(card.Spades, card.Seven),
		mustMakeCard(card.Spades, card.Six),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.StraightFlush {
		t.Errorf("Expected StraightFlush, got %v", result.Tier())
	}
}

func TestHandEvaluator_RoyalFlush(t *testing.T) {
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Spades, card.King),
		mustMakeCard(card.Spades, card.Queen),
		mustMakeCard(card.Spades, card.Jack),
		mustMakeCard(card.Spades, card.Ten),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.RoyalFlush {
		t.Errorf("Expected RoyalFlush, got %v", result.Tier())
	}
}

func TestHandEvaluator_WheelStraight(t *testing.T) {
	// Test A-2-3-4-5 (wheel/bicycle)
	evaluator := NewHandEvaluator()

	cards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Five),
		mustMakeCard(card.Diamonds, card.Four),
		mustMakeCard(card.Clubs, card.Three),
		mustMakeCard(card.Spades, card.Two),
	}

	result, err := evaluator.EvaluateFiveCards(cards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	if result.Tier() != vo.Straight {
		t.Errorf("Expected Straight (wheel), got %v", result.Tier())
	}
}

func TestHandEvaluator_CompareHands(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Flush vs Full House (Full House wins)
	flush := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Spades, card.King),
		mustMakeCard(card.Spades, card.Ten),
		mustMakeCard(card.Spades, card.Seven),
		mustMakeCard(card.Spades, card.Three),
	}

	fullHouse := []card.Card{
		mustMakeCard(card.Hearts, card.Ace),
		mustMakeCard(card.Diamonds, card.Ace),
		mustMakeCard(card.Clubs, card.Ace),
		mustMakeCard(card.Hearts, card.King),
		mustMakeCard(card.Diamonds, card.King),
	}

	flushResult, _ := evaluator.EvaluateFiveCards(flush)
	fullHouseResult, _ := evaluator.EvaluateFiveCards(fullHouse)

	comparison := fullHouseResult.CompareTo(flushResult)
	if comparison <= 0 {
		t.Errorf("Expected Full House to beat Flush, got comparison %d", comparison)
	}
}

func TestHandEvaluator_EvaluateWithCommunityCards(t *testing.T) {
	evaluator := NewHandEvaluator()

	playerCards := []card.Card{
		mustMakeCard(card.Spades, card.Ace),
		mustMakeCard(card.Hearts, card.Ace),
	}

	communityCards := []card.Card{
		mustMakeCard(card.Diamonds, card.Ace),
		mustMakeCard(card.Clubs, card.King),
		mustMakeCard(card.Spades, card.King),
		mustMakeCard(card.Hearts, card.Seven),
		mustMakeCard(card.Diamonds, card.Three),
	}

	result, err := evaluator.Evaluate(playerCards, communityCards)
	if err != nil {
		t.Fatalf("Evaluation failed: %v", err)
	}

	// Should find Full House (Ace Ace Ace King King)
	if result.Tier() != vo.FullHouse {
		t.Errorf("Expected FullHouse, got %v", result.Tier())
	}
}

// Helper function to create card without error handling
func mustMakeCard(suit card.Suit, rank card.Rank) card.Card {
	c, _ := card.NewCard(suit, rank)
	return c
}
