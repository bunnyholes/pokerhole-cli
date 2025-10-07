package service

import (
	"fmt"
	"time"

	"github.com/bunnyholes/pokerhole/client/internal/adapter/out/deck"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

// OfflineGame represents a simple offline 2-player game
type OfflineGame struct {
	deck           card.DeckPort
	gameService    *GameService
	players        []*player.Player
	communityCards []card.Card
	round          vo.BettingRound
	pot            int
	currentBet     int
	currentPlayer  int
	gameState      game.GameState
}

// NewOfflineGame creates a new offline game with 2 players
func NewOfflineGame(userNickname string) *OfflineGame {
	// Create deck
	localDeck := deck.NewLocalDeck()
	seed := time.Now().UnixNano()
	localDeck.Shuffle(seed)

	// Create game service
	handEvaluator := game.NewHandEvaluator()
	gameService := NewGameService(localDeck, handEvaluator)

	// Create players
	userID := player.GeneratePlayerId()
	aiID := player.GeneratePlayerId()

	userNick, _ := player.NewNickname(userNickname)
	aiNick, _ := player.NewNickname("AI Player")

	userPlayer, _ := player.NewPlayer(userID, userNick, 1000)
	aiPlayer, _ := player.NewPlayer(aiID, aiNick, 1000)

	players := []*player.Player{userPlayer, aiPlayer}

	return &OfflineGame{
		deck:           localDeck,
		gameService:    gameService,
		players:        players,
		communityCards: make([]card.Card, 0),
		round:          vo.PreFlop,
		pot:            0,
		currentBet:     0,
		currentPlayer:  0,
		gameState:      game.Waiting,
	}
}

// Start starts the game
func (g *OfflineGame) Start() error {
	g.gameState = game.Playing

	// Deal hole cards
	err := g.gameService.DealHoleCards(g.players)
	if err != nil {
		return fmt.Errorf("failed to deal hole cards: %w", err)
	}

	// Set small blind and big blind
	g.players[0].PlaceBet(10)  // Small blind
	g.players[1].PlaceBet(20)  // Big blind
	g.pot = 30
	g.currentBet = 20

	return nil
}

// GetGameState returns the current game state as a snapshot
func (g *OfflineGame) GetGameState() GameStateSnapshot {
	snapshot := GameStateSnapshot{
		Round:          g.round.String(),
		Pot:            g.pot,
		CurrentBet:     g.currentBet,
		CommunityCards: formatCards(g.communityCards),
		Players:        formatPlayers(g.players),
		CurrentPlayer:  g.currentPlayer,
		WinnerIndex:    -1, // No winner by default
	}

	// If Showdown, evaluate hands and determine winner
	if g.round == vo.Showdown && len(g.communityCards) == 5 {
		snapshot = g.evaluateShowdown(snapshot)
	}

	return snapshot
}

// GameStateSnapshot represents a snapshot of the game state for UI
type GameStateSnapshot struct {
	Round          string
	Pot            int
	CurrentBet     int
	CommunityCards []string
	Players        []PlayerSnapshot
	CurrentPlayer  int
	WinnerIndex    int    // -1 if no winner yet, player index if there's a winner
	WinnerHandRank string // Hand ranking description (e.g., "One Pair", "Flush")
}

// PlayerSnapshot represents a player's state for UI
type PlayerSnapshot struct {
	Nickname  string
	Chips     int
	Bet       int
	Status    string
	Hand      string
	HandRank  string   // Best hand ranking (e.g., "One Pair - Queens")
	BestCards []string // Best 5 cards used for hand evaluation
}

// formatCards formats cards for display
func formatCards(cards []card.Card) []string {
	result := make([]string, len(cards))
	for i, c := range cards {
		result[i] = c.String()
	}
	return result
}

// evaluateShowdown evaluates all player hands and determines the winner
func (g *OfflineGame) evaluateShowdown(snapshot GameStateSnapshot) GameStateSnapshot {
	// Evaluate each player's hand
	type playerResult struct {
		index      int
		handResult vo.HandResult
	}

	results := make([]playerResult, 0)

	for i, p := range g.players {
		// Skip folded players
		if p.Status() == player.Folded {
			continue
		}

		playerCards := p.Hand().Cards()
		handResult, err := g.gameService.HandEvaluator.Evaluate(playerCards, g.communityCards)
		if err != nil {
			// If evaluation fails, skip this player
			continue
		}

		results = append(results, playerResult{
			index:      i,
			handResult: handResult,
		})

		// Update PlayerSnapshot with hand rank and rank cards
		snapshot.Players[i].HandRank = handResult.String()
		snapshot.Players[i].BestCards = formatCards(handResult.GetRankCards())
	}

	// Find winner (player with best hand)
	if len(results) > 0 {
		winnerIdx := 0
		for i := 1; i < len(results); i++ {
			if results[i].handResult.CompareTo(results[winnerIdx].handResult) > 0 {
				winnerIdx = i
			}
		}

		snapshot.WinnerIndex = results[winnerIdx].index
		snapshot.WinnerHandRank = results[winnerIdx].handResult.String()
	}

	return snapshot
}

// formatPlayers formats players for display
func formatPlayers(players []*player.Player) []PlayerSnapshot {
	result := make([]PlayerSnapshot, len(players))
	for i, p := range players {
		result[i] = PlayerSnapshot{
			Nickname:  p.Nickname().String(),
			Chips:     p.Chips(),
			Bet:       p.Bet(),
			Status:    p.Status().String(),
			Hand:      p.Hand().String(),
			HandRank:  "", // Will be filled during showdown
			BestCards: []string{}, // Will be filled during showdown
		}
	}
	return result
}

// PlayerAction executes a player action
func (g *OfflineGame) PlayerAction(playerIndex int, action vo.PlayerAction, amount int) error {
	if playerIndex >= len(g.players) {
		return fmt.Errorf("invalid player index: %d", playerIndex)
	}

	p := g.players[playerIndex]

	switch action {
	case vo.Fold:
		p.Fold()
	case vo.Call:
		callAmount := g.currentBet - p.Bet()
		if err := p.PlaceBet(callAmount); err != nil {
			return err
		}
		g.pot += callAmount
	case vo.Raise:
		raiseAmount := amount - p.Bet()
		if err := p.PlaceBet(raiseAmount); err != nil {
			return err
		}
		g.pot += raiseAmount
		g.currentBet = amount
	case vo.AllIn:
		allInAmount := p.Chips()
		p.AllIn()
		g.pot += allInAmount
		// Update currentBet if all-in amount exceeds it
		totalBet := p.Bet()
		if totalBet > g.currentBet {
			g.currentBet = totalBet
		}
	case vo.Check:
		// Do nothing
	}

	// Move to next player
	g.currentPlayer = (g.currentPlayer + 1) % len(g.players)

	return nil
}

// ProgressRound progresses to the next betting round
func (g *OfflineGame) ProgressRound() error {
	switch g.round {
	case vo.PreFlop:
		// Deal flop
		cards, err := g.gameService.DealFlop()
		if err != nil {
			return err
		}
		g.communityCards = append(g.communityCards, cards...)
		g.round = vo.Flop
		g.currentBet = 0
		g.currentPlayer = 0 // Reset to first player
		for _, p := range g.players {
			p.ResetBet()
		}

	case vo.Flop:
		// Deal turn
		turnCard, err := g.gameService.DealTurn()
		if err != nil {
			return err
		}
		g.communityCards = append(g.communityCards, turnCard)
		g.round = vo.Turn
		g.currentBet = 0
		g.currentPlayer = 0 // Reset to first player
		for _, p := range g.players {
			p.ResetBet()
		}

	case vo.Turn:
		// Deal river
		riverCard, err := g.gameService.DealRiver()
		if err != nil {
			return err
		}
		g.communityCards = append(g.communityCards, riverCard)
		g.round = vo.River
		g.currentBet = 0
		g.currentPlayer = 0 // Reset to first player
		for _, p := range g.players {
			p.ResetBet()
		}

	case vo.River:
		// Showdown
		g.round = vo.Showdown
		g.gameState = game.Finished

		// Determine winner and distribute pot
		return g.resolveShowdown()
	}

	return nil
}

// resolveShowdown determines the winner and distributes the pot
func (g *OfflineGame) resolveShowdown() error {
	winnerResolver := game.NewWinnerResolver(g.gameService.HandEvaluator)

	winners, err := winnerResolver.DetermineWinners(g.players, g.communityCards)
	if err != nil {
		return fmt.Errorf("failed to determine winners: %w", err)
	}

	if len(winners) == 0 {
		return fmt.Errorf("no winners found")
	}

	// Distribute pot evenly among winners
	potShare := g.pot / len(winners)
	remainder := g.pot % len(winners)

	for i, winner := range winners {
		share := potShare
		// Give remainder to first winner
		if i == 0 {
			share += remainder
		}
		winner.AddChips(share)
	}

	g.pot = 0
	return nil
}

// GetPlayers returns the players
func (g *OfflineGame) GetPlayers() []*player.Player {
	return g.players
}

// GetCommunityCards returns the community cards
func (g *OfflineGame) GetCommunityCards() []card.Card {
	return g.communityCards
}

// GetWinners returns the winners (only valid after showdown)
func (g *OfflineGame) GetWinners() ([]*player.Player, error) {
	if g.round != vo.Showdown {
		return nil, fmt.Errorf("game not in showdown state")
	}

	winnerResolver := game.NewWinnerResolver(g.gameService.HandEvaluator)
	return winnerResolver.DetermineWinners(g.players, g.communityCards)
}

// Restart resets the game for a new hand
func (g *OfflineGame) Restart() error {
	// Check if any player has run out of chips
	for _, p := range g.players {
		if p.Chips() <= 0 {
			g.gameState = game.Finished
			return fmt.Errorf("player %s has no chips left - game over", p.Nickname())
		}
	}

	// Reset deck
	if err := g.deck.Reset(); err != nil {
		return err
	}

	seed := time.Now().UnixNano()
	if err := g.deck.Shuffle(seed); err != nil {
		return err
	}

	// Reset game state
	g.communityCards = make([]card.Card, 0)
	g.round = vo.PreFlop
	g.pot = 0
	g.currentBet = 0
	g.currentPlayer = 0
	g.gameState = game.Waiting

	// Reset players
	for _, p := range g.players {
		p.ResetBet()
		p.SetStatus(player.Waiting)
		p.SetHand(card.NewHand([]card.Card{}))
	}

	// Start new game
	return g.Start()
}
