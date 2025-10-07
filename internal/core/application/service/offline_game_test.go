package service

import (
	"testing"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

func TestNewOfflineGame(t *testing.T) {
	game := NewOfflineGame("TestPlayer")

	if game == nil {
		t.Fatal("NewOfflineGame returned nil")
	}

	// Check initial state
	players := game.GetPlayers()
	if len(players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(players))
	}

	// Check player names
	if players[0].Nickname().String() != "TestPlayer" {
		t.Errorf("Expected user player nickname 'TestPlayer', got '%s'", players[0].Nickname().String())
	}
	if players[1].Nickname().String() != "AI Player" {
		t.Errorf("Expected AI player nickname 'AI Player', got '%s'", players[1].Nickname().String())
	}

	// Check initial chips
	if players[0].Chips() != 1000 {
		t.Errorf("Expected player 0 to have 1000 chips, got %d", players[0].Chips())
	}
	if players[1].Chips() != 1000 {
		t.Errorf("Expected player 1 to have 1000 chips, got %d", players[1].Chips())
	}
}

func TestOfflineGame_Start(t *testing.T) {
	game := NewOfflineGame("TestPlayer")

	err := game.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Check game state
	state := game.GetGameState()

	// Check round is PreFlop
	if state.Round != "PRE_FLOP" {
		t.Errorf("Expected round PRE_FLOP, got %s", state.Round)
	}

	// Check pot (small blind 10 + big blind 20 = 30)
	if state.Pot != 30 {
		t.Errorf("Expected pot 30, got %d", state.Pot)
	}

	// Check current bet (big blind = 20)
	if state.CurrentBet != 20 {
		t.Errorf("Expected current bet 20, got %d", state.CurrentBet)
	}

	// Check players have hole cards
	players := game.GetPlayers()
	for i, p := range players {
		hand := p.Hand()
		cards := hand.Cards()
		if len(cards) != 2 {
			t.Errorf("Player %d: expected 2 hole cards, got %d", i, len(cards))
		}
	}

	// Check blinds were posted
	if players[0].Bet() != 10 {
		t.Errorf("Expected player 0 (small blind) bet 10, got %d", players[0].Bet())
	}
	if players[1].Bet() != 20 {
		t.Errorf("Expected player 1 (big blind) bet 20, got %d", players[1].Bet())
	}

	// Check chips deducted
	if players[0].Chips() != 990 {
		t.Errorf("Expected player 0 chips 990, got %d", players[0].Chips())
	}
	if players[1].Chips() != 980 {
		t.Errorf("Expected player 1 chips 980, got %d", players[1].Chips())
	}
}

func TestPlayerAction_Fold(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	players := game.GetPlayers()

	err := game.PlayerAction(0, vo.Fold, 0)
	if err != nil {
		t.Fatalf("Fold action failed: %v", err)
	}

	if players[0].Status() != player.Folded {
		t.Errorf("Expected player 0 status Folded, got %v", players[0].Status())
	}

	// Current player should move to player 1
	state := game.GetGameState()
	if state.CurrentPlayer != 1 {
		t.Errorf("Expected current player 1, got %d", state.CurrentPlayer)
	}
}

func TestPlayerAction_Call(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	players := game.GetPlayers()
	initialPot := game.GetGameState().Pot

	// Player 0 (small blind 10) calls big blind (20)
	// Needs to add 10 more (20 - 10 = 10)
	err := game.PlayerAction(0, vo.Call, 0)
	if err != nil {
		t.Fatalf("Call action failed: %v", err)
	}

	state := game.GetGameState()

	// Pot should increase by 10 (30 + 10 = 40)
	if state.Pot != initialPot+10 {
		t.Errorf("Expected pot %d, got %d", initialPot+10, state.Pot)
	}

	// Player 0 bet should be 20 now
	if players[0].Bet() != 20 {
		t.Errorf("Expected player 0 bet 20, got %d", players[0].Bet())
	}

	// Chips should be 980 (1000 - 20)
	if players[0].Chips() != 980 {
		t.Errorf("Expected player 0 chips 980, got %d", players[0].Chips())
	}
}

func TestPlayerAction_Raise(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	players := game.GetPlayers()
	initialPot := game.GetGameState().Pot

	// Player 0 (bet=10) raises to 50
	// Needs to add 40 more (50 - 10 = 40)
	err := game.PlayerAction(0, vo.Raise, 50)
	if err != nil {
		t.Fatalf("Raise action failed: %v", err)
	}

	state := game.GetGameState()

	// Pot should increase by 40
	if state.Pot != initialPot+40 {
		t.Errorf("Expected pot %d, got %d", initialPot+40, state.Pot)
	}

	// Current bet should be 50
	if state.CurrentBet != 50 {
		t.Errorf("Expected current bet 50, got %d", state.CurrentBet)
	}

	// Player 0 bet should be 50
	if players[0].Bet() != 50 {
		t.Errorf("Expected player 0 bet 50, got %d", players[0].Bet())
	}
}

func TestPlayerAction_AllIn(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	players := game.GetPlayers()

	err := game.PlayerAction(0, vo.AllIn, 0)
	if err != nil {
		t.Fatalf("All-in action failed: %v", err)
	}

	if players[0].Status() != player.AllIn {
		t.Errorf("Expected player 0 status AllIn, got %v", players[0].Status())
	}

	if players[0].Chips() != 0 {
		t.Errorf("Expected player 0 chips 0, got %d", players[0].Chips())
	}
}

func TestPlayerAction_Check(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	players := game.GetPlayers()

	// Reset bets to allow check
	players[0].ResetBet()
	players[1].ResetBet()

	err := game.PlayerAction(0, vo.Check, 0)
	if err != nil {
		t.Fatalf("Check action failed: %v", err)
	}

	// Player should still be Active
	if players[0].Status() != player.Active && players[0].Status() != player.Waiting {
		t.Errorf("Expected player 0 status Active or Waiting, got %v", players[0].Status())
	}
}

func TestPlayerAction_InvalidPlayerIndex(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	err := game.PlayerAction(5, vo.Fold, 0)
	if err == nil {
		t.Error("Expected error for invalid player index, got nil")
	}
}

func TestProgressRound_PreFlopToFlop(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	initialCommunityCards := len(game.GetCommunityCards())
	if initialCommunityCards != 0 {
		t.Errorf("Expected 0 initial community cards, got %d", initialCommunityCards)
	}

	err := game.ProgressRound()
	if err != nil {
		t.Fatalf("ProgressRound failed: %v", err)
	}

	state := game.GetGameState()

	// Check round progressed to Flop
	if state.Round != "FLOP" {
		t.Errorf("Expected round FLOP, got %s", state.Round)
	}

	// Check 3 community cards were dealt
	communityCards := game.GetCommunityCards()
	if len(communityCards) != 3 {
		t.Errorf("Expected 3 community cards, got %d", len(communityCards))
	}

	// Check bets were reset
	if state.CurrentBet != 0 {
		t.Errorf("Expected current bet reset to 0, got %d", state.CurrentBet)
	}

	players := game.GetPlayers()
	if players[0].Bet() != 0 {
		t.Errorf("Expected player 0 bet reset to 0, got %d", players[0].Bet())
	}
}

func TestProgressRound_FlopToTurn(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	// Progress to Flop
	game.ProgressRound()

	// Progress to Turn
	err := game.ProgressRound()
	if err != nil {
		t.Fatalf("ProgressRound failed: %v", err)
	}

	state := game.GetGameState()

	// Check round progressed to Turn
	if state.Round != "TURN" {
		t.Errorf("Expected round TURN, got %s", state.Round)
	}

	// Check 4 community cards total (3 flop + 1 turn)
	communityCards := game.GetCommunityCards()
	if len(communityCards) != 4 {
		t.Errorf("Expected 4 community cards, got %d", len(communityCards))
	}
}

func TestProgressRound_TurnToRiver(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	// Progress through rounds
	game.ProgressRound() // Flop
	game.ProgressRound() // Turn

	// Progress to River
	err := game.ProgressRound()
	if err != nil {
		t.Fatalf("ProgressRound failed: %v", err)
	}

	state := game.GetGameState()

	// Check round progressed to River
	if state.Round != "RIVER" {
		t.Errorf("Expected round RIVER, got %s", state.Round)
	}

	// Check 5 community cards total (3 flop + 1 turn + 1 river)
	communityCards := game.GetCommunityCards()
	if len(communityCards) != 5 {
		t.Errorf("Expected 5 community cards, got %d", len(communityCards))
	}
}

func TestProgressRound_RiverToShowdown(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	// Progress through all rounds
	game.ProgressRound() // Flop
	game.ProgressRound() // Turn
	game.ProgressRound() // River

	// Progress to Showdown
	err := game.ProgressRound()
	if err != nil {
		t.Fatalf("ProgressRound failed: %v", err)
	}

	state := game.GetGameState()

	// Check round progressed to Showdown
	if state.Round != "SHOWDOWN" {
		t.Errorf("Expected round SHOWDOWN, got %s", state.Round)
	}
}

func TestGetGameState_Snapshot(t *testing.T) {
	game := NewOfflineGame("TestPlayer")
	game.Start()

	state := game.GetGameState()

	// Verify all fields are populated correctly
	if state.Round == "" {
		t.Error("GameState Round is empty")
	}
	if state.Pot <= 0 {
		t.Error("GameState Pot should be positive")
	}
	if len(state.Players) != 2 {
		t.Errorf("Expected 2 players in GameState, got %d", len(state.Players))
	}

	// Check player snapshots
	for i, p := range state.Players {
		if p.Nickname == "" {
			t.Errorf("Player %d nickname is empty", i)
		}
		if p.Hand == "" {
			t.Errorf("Player %d hand is empty", i)
		}
		if p.Status == "" {
			t.Errorf("Player %d status is empty", i)
		}
	}
}

func TestCompleteGameFlow(t *testing.T) {
	// Test a complete game flow from start to showdown
	game := NewOfflineGame("TestPlayer")

	// Start game
	err := game.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	state := game.GetGameState()
	if state.Round != "PRE_FLOP" {
		t.Errorf("Expected PRE_FLOP, got %s", state.Round)
	}

	// Pre-flop: Player 0 calls
	game.PlayerAction(0, vo.Call, 0)
	// Player 1 checks
	game.PlayerAction(1, vo.Check, 0)

	// Progress to Flop
	game.ProgressRound()
	state = game.GetGameState()
	if state.Round != "FLOP" {
		t.Errorf("Expected FLOP, got %s", state.Round)
	}
	if len(game.GetCommunityCards()) != 3 {
		t.Errorf("Expected 3 flop cards, got %d", len(game.GetCommunityCards()))
	}

	// Flop: Both check
	game.PlayerAction(0, vo.Check, 0)
	game.PlayerAction(1, vo.Check, 0)

	// Progress to Turn
	game.ProgressRound()
	state = game.GetGameState()
	if state.Round != "TURN" {
		t.Errorf("Expected TURN, got %s", state.Round)
	}
	if len(game.GetCommunityCards()) != 4 {
		t.Errorf("Expected 4 cards (flop+turn), got %d", len(game.GetCommunityCards()))
	}

	// Turn: Player 0 bets, Player 1 calls
	game.PlayerAction(0, vo.Raise, 50)
	game.PlayerAction(1, vo.Call, 0)

	// Progress to River
	game.ProgressRound()
	state = game.GetGameState()
	if state.Round != "RIVER" {
		t.Errorf("Expected RIVER, got %s", state.Round)
	}
	if len(game.GetCommunityCards()) != 5 {
		t.Errorf("Expected 5 cards (flop+turn+river), got %d", len(game.GetCommunityCards()))
	}

	// River: Both check
	game.PlayerAction(0, vo.Check, 0)
	game.PlayerAction(1, vo.Check, 0)

	// Progress to Showdown
	game.ProgressRound()
	state = game.GetGameState()
	if state.Round != "SHOWDOWN" {
		t.Errorf("Expected SHOWDOWN, got %s", state.Round)
	}
}
