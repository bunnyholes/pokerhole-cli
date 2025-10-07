//go:build ignore
// +build ignore

// Package deck provides deck adapter implementations
package deck

import (
	"github.com/bunnyholes/pokerhole/client/internal/adapter/out/websocket"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

// RemoteDeck is an online deck implementation (Adapter)
// Implements: card.DeckPort
// Requests cards from server instead of managing them locally
type RemoteDeck struct {
	client *websocket.Client
	gameID string
}

// NewRemoteDeck creates a new RemoteDeck
func NewRemoteDeck(client *websocket.Client, gameID string) *RemoteDeck {
	return &RemoteDeck{
		client: client,
		gameID: gameID,
	}
}

// Compile-time check: RemoteDeck implements card.DeckPort
var _ card.DeckPort = (*RemoteDeck)(nil)

// DrawCard requests a card from the server
func (d *RemoteDeck) DrawCard() (card.Card, error) {
	// TODO: implement
	// NOTE: In practice, server auto-sends cards (we just parse)
	// This is more of a placeholder
	return card.Card{}, nil
}

// Shuffle does nothing (server shuffles)
func (d *RemoteDeck) Shuffle(seed int64) error {
	// Server handles shuffling
	// Client doesn't need to do anything
	return nil
}

// RemainingCards returns 0 (client doesn't know)
func (d *RemoteDeck) RemainingCards() int {
	// Server tracks this, client doesn't need to know
	return 0
}

// Reset does nothing (server handles)
func (d *RemoteDeck) Reset() error {
	return nil
}
