package player

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

func TestNewPlayer(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	chips := 1000

	player, err := NewPlayer(id, nickname, chips)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if player.ID() != id {
		t.Errorf("Expected ID %v, got %v", id, player.ID())
	}
	if player.Nickname() != nickname {
		t.Errorf("Expected nickname %v, got %v", nickname, player.Nickname())
	}
	if player.Chips() != chips {
		t.Errorf("Expected chips %d, got %d", chips, player.Chips())
	}
	if player.Bet() != 0 {
		t.Errorf("Expected bet 0, got %d", player.Bet())
	}
	if player.Status() != Waiting {
		t.Errorf("Expected status Waiting, got %v", player.Status())
	}
}

func TestPlaceBet_Success(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	err := player.PlaceBet(100)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if player.Chips() != 900 {
		t.Errorf("Expected chips 900, got %d", player.Chips())
	}
	if player.Bet() != 100 {
		t.Errorf("Expected bet 100, got %d", player.Bet())
	}
	if player.Status() != Active {
		t.Errorf("Expected status Active, got %v", player.Status())
	}
}

func TestPlaceBet_InvalidAmount(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	tests := []struct {
		name   string
		amount int
	}{
		{"Zero amount", 0},
		{"Negative amount", -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := player.PlaceBet(tt.amount)
			if err != ErrInvalidBetAmount {
				t.Errorf("Expected ErrInvalidBetAmount, got %v", err)
			}
		})
	}
}

func TestPlaceBet_InsufficientChips(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 100)

	err := player.PlaceBet(200)
	if err != ErrInsufficientChips {
		t.Errorf("Expected ErrInsufficientChips, got %v", err)
	}

	// Verify state didn't change
	if player.Chips() != 100 {
		t.Errorf("Expected chips unchanged at 100, got %d", player.Chips())
	}
	if player.Bet() != 0 {
		t.Errorf("Expected bet unchanged at 0, got %d", player.Bet())
	}
}

func TestPlaceBet_Multiple(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	// First bet
	player.PlaceBet(100)
	// Second bet
	player.PlaceBet(50)

	if player.Chips() != 850 {
		t.Errorf("Expected chips 850, got %d", player.Chips())
	}
	if player.Bet() != 150 {
		t.Errorf("Expected bet 150, got %d", player.Bet())
	}
}

func TestFold(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	player.PlaceBet(100) // Active state
	player.Fold()

	if player.Status() != Folded {
		t.Errorf("Expected status Folded, got %v", player.Status())
	}
	// Chips and bet should remain unchanged
	if player.Chips() != 900 {
		t.Errorf("Expected chips 900, got %d", player.Chips())
	}
	if player.Bet() != 100 {
		t.Errorf("Expected bet 100, got %d", player.Bet())
	}
}

func TestAllIn(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 500)

	player.AllIn()

	if player.Status() != AllIn {
		t.Errorf("Expected status AllIn, got %v", player.Status())
	}
	if player.Chips() != 0 {
		t.Errorf("Expected chips 0, got %d", player.Chips())
	}
	if player.Bet() != 500 {
		t.Errorf("Expected bet 500, got %d", player.Bet())
	}
}

func TestAllIn_AfterBetting(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 500)

	player.PlaceBet(200) // Bet 200 first
	player.AllIn()       // All-in with remaining 300

	if player.Chips() != 0 {
		t.Errorf("Expected chips 0, got %d", player.Chips())
	}
	if player.Bet() != 500 { // 200 + 300
		t.Errorf("Expected bet 500, got %d", player.Bet())
	}
}

func TestResetBet(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	player.PlaceBet(200)
	player.ResetBet()

	if player.Bet() != 0 {
		t.Errorf("Expected bet 0, got %d", player.Bet())
	}
	// Chips should remain at 800 (bet was already deducted)
	if player.Chips() != 800 {
		t.Errorf("Expected chips 800, got %d", player.Chips())
	}
}

func TestAddChips(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 500)

	player.AddChips(300)

	if player.Chips() != 800 {
		t.Errorf("Expected chips 800, got %d", player.Chips())
	}
}

func TestSetHand(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	card1, _ := card.NewCard(card.Spades, card.Ace)
	card2, _ := card.NewCard(card.Hearts, card.King)
	hand := card.NewHand([]card.Card{card1, card2})

	player.SetHand(hand)

	playerHand := player.Hand()
	cards := playerHand.Cards()

	if len(cards) != 2 {
		t.Fatalf("Expected 2 cards, got %d", len(cards))
	}
	if cards[0] != card1 {
		t.Errorf("Expected first card %v, got %v", card1, cards[0])
	}
	if cards[1] != card2 {
		t.Errorf("Expected second card %v, got %v", card2, cards[1])
	}
}

func TestSetPosition(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	player.SetPosition(2)

	if player.Position() != 2 {
		t.Errorf("Expected position 2, got %d", player.Position())
	}
}

func TestSetStatus(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	player.SetStatus(Active)

	if player.Status() != Active {
		t.Errorf("Expected status Active, got %v", player.Status())
	}
}

func TestPlayerString(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("TestPlayer")
	player, _ := NewPlayer(id, nickname, 1000)

	str := player.String()
	expected := "TestPlayer"

	if str != expected {
		t.Errorf("Expected string %s, got %s", expected, str)
	}
}

// Test betting round scenario
func TestBettingRoundScenario(t *testing.T) {
	id := GeneratePlayerId()
	nickname, _ := NewNickname("Player1")
	player, _ := NewPlayer(id, nickname, 1000)

	// Pre-flop: place small blind
	if err := player.PlaceBet(10); err != nil {
		t.Fatalf("Small blind failed: %v", err)
	}

	// Call big blind
	if err := player.PlaceBet(10); err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	// Check state after pre-flop
	if player.Chips() != 980 {
		t.Errorf("Expected chips 980, got %d", player.Chips())
	}
	if player.Bet() != 20 {
		t.Errorf("Expected bet 20, got %d", player.Bet())
	}

	// New round: reset bet
	player.ResetBet()
	if player.Bet() != 0 {
		t.Errorf("Expected bet reset to 0, got %d", player.Bet())
	}

	// Flop: check (no bet)
	// Turn: bet 50
	if err := player.PlaceBet(50); err != nil {
		t.Fatalf("Turn bet failed: %v", err)
	}

	if player.Chips() != 930 {
		t.Errorf("Expected chips 930, got %d", player.Chips())
	}
}
