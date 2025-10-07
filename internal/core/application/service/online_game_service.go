//go:build ignore
// +build ignore

// Package service provides application services
package service

import (
	"github.com/bunnyholes/pokerhole/client/internal/adapter/out/websocket"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

// OnlineGameService handles online game mode (Application Service)
type OnlineGameService struct {
	client      *websocket.Client
	gameService *GameService
}

// NewOnlineGameService creates a new OnlineGameService
func NewOnlineGameService(client *websocket.Client, gameService *GameService) *OnlineGameService {
	return &OnlineGameService{
		client:      client,
		gameService: gameService,
	}
}

// JoinGame joins an online game
func (s *OnlineGameService) JoinGame(playerID player.PlayerId) error {
	// TODO: implement
	// 1. Send JOIN message to server
	// 2. Wait for confirmation
	return nil
}

// ExecuteAction sends action to server (online mode)
func (s *OnlineGameService) ExecuteAction(gameID game.GameId, playerID player.PlayerId, action game.PlayerAction, amount int) error {
	// TODO: implement
	// 1. Validate action locally (UX)
	// 2. Send action to server
	// 3. Wait for server confirmation
	return nil
}

// SyncGameState syncs game state from server
func (s *OnlineGameService) SyncGameState(gameID game.GameId) error {
	// TODO: implement
	// 1. Request game state from server
	// 2. Update local state
	return nil
}
