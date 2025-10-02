package state

import (
	"sync"

	"github.com/bunnyholes/pokerhole/client/internal/network"
)

// GameState represents the current game state
type GameState struct {
	GameID         string
	Round          string
	Pot            int
	CurrentBet     int
	CommunityCards []string
	Players        []PlayerState
	CurrentPlayer  string
	ValidActions   []string
	mu             sync.RWMutex
}

// PlayerState represents a player's state
type PlayerState struct {
	ID       string
	Nickname string
	Chips    int
	Bet      int
	Status   string
	Position int
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	return &GameState{
		Players:        make([]PlayerState, 0),
		CommunityCards: make([]string, 0),
		ValidActions:   make([]string, 0),
	}
}

// Update updates the game state from a server message
func (s *GameState) Update(msg network.ServerMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload := msg.Payload
	if payload == nil {
		return
	}

	// Update fields from payload
	if gameID, ok := payload["gameId"].(string); ok {
		s.GameID = gameID
	}
	if round, ok := payload["round"].(string); ok {
		s.Round = round
	}
	if pot, ok := payload["pot"].(float64); ok {
		s.Pot = int(pot)
	}
	if currentBet, ok := payload["currentBet"].(float64); ok {
		s.CurrentBet = int(currentBet)
	}
	if currentPlayer, ok := payload["currentPlayer"].(string); ok {
		s.CurrentPlayer = currentPlayer
	}

	// Update community cards
	if cards, ok := payload["communityCards"].([]interface{}); ok {
		s.CommunityCards = make([]string, 0, len(cards))
		for _, card := range cards {
			if cardStr, ok := card.(string); ok {
				s.CommunityCards = append(s.CommunityCards, cardStr)
			}
		}
	}

	// Update players
	if players, ok := payload["players"].([]interface{}); ok {
		s.Players = make([]PlayerState, 0, len(players))
		for _, player := range players {
			if playerMap, ok := player.(map[string]interface{}); ok {
				ps := PlayerState{}
				if id, ok := playerMap["id"].(string); ok {
					ps.ID = id
				}
				if nickname, ok := playerMap["nickname"].(string); ok {
					ps.Nickname = nickname
				}
				if chips, ok := playerMap["chips"].(float64); ok {
					ps.Chips = int(chips)
				}
				if bet, ok := playerMap["bet"].(float64); ok {
					ps.Bet = int(bet)
				}
				if status, ok := playerMap["status"].(string); ok {
					ps.Status = status
				}
				if position, ok := playerMap["position"].(float64); ok {
					ps.Position = int(position)
				}
				s.Players = append(s.Players, ps)
			}
		}
	}

	// Update valid actions
	if actions, ok := payload["validActions"].([]interface{}); ok {
		s.ValidActions = make([]string, 0, len(actions))
		for _, action := range actions {
			if actionStr, ok := action.(string); ok {
				s.ValidActions = append(s.ValidActions, actionStr)
			}
		}
	}
}

// GetSnapshot returns a thread-safe copy of the game state
func (s *GameState) GetSnapshot() GameState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := GameState{
		GameID:         s.GameID,
		Round:          s.Round,
		Pot:            s.Pot,
		CurrentBet:     s.CurrentBet,
		CurrentPlayer:  s.CurrentPlayer,
		CommunityCards: make([]string, len(s.CommunityCards)),
		Players:        make([]PlayerState, len(s.Players)),
		ValidActions:   make([]string, len(s.ValidActions)),
	}

	copy(snapshot.CommunityCards, s.CommunityCards)
	copy(snapshot.Players, s.Players)
	copy(snapshot.ValidActions, s.ValidActions)

	return snapshot
}
