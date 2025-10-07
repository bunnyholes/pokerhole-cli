// Package player provides player domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/player/Player.java
package player

import (
	"errors"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

var (
	ErrInvalidBetAmount  = errors.New("invalid bet amount")
	ErrInsufficientChips = errors.New("insufficient chips")
)

// Player represents a poker player (Aggregate Root)
type Player struct {
	id       PlayerId
	nickname Nickname
	chips    int
	bet      int
	status   PlayerStatus
	hand     card.Hand
	position int
}

// NewPlayer creates a new Player
func NewPlayer(id PlayerId, nickname Nickname, chips int) (*Player, error) {
	// TODO: validate chips > 0
	return &Player{
		id:       id,
		nickname: nickname,
		chips:    chips,
		bet:      0,
		status:   Waiting,
		hand:     card.NewHand([]card.Card{}), // Initialize with empty hand
		position: 0,
	}, nil
}

// ID returns the player ID
func (p *Player) ID() PlayerId {
	return p.id
}

// Nickname returns the player nickname
func (p *Player) Nickname() Nickname {
	return p.nickname
}

// Chips returns the player's chip count
func (p *Player) Chips() int {
	return p.chips
}

// Bet returns the player's current bet
func (p *Player) Bet() int {
	return p.bet
}

// Status returns the player status
func (p *Player) Status() PlayerStatus {
	return p.status
}

// Hand returns the player's hand
func (p *Player) Hand() card.Hand {
	return p.hand
}

// Position returns the player's position at the table
func (p *Player) Position() int {
	return p.position
}

// PlaceBet places a bet (deducts from chips, adds to bet)
func (p *Player) PlaceBet(amount int) error {
	if amount <= 0 {
		return ErrInvalidBetAmount
	}
	if amount > p.chips {
		return ErrInsufficientChips
	}

	p.chips -= amount
	p.bet += amount
	p.status = Active

	return nil
}

// Fold marks the player as folded
func (p *Player) Fold() {
	p.status = Folded
}

// AllIn marks the player as all-in
func (p *Player) AllIn() {
	allInAmount := p.chips
	p.chips = 0
	p.bet += allInAmount
	p.status = AllIn
}

// ResetBet resets the player's bet (e.g., new round)
func (p *Player) ResetBet() {
	p.bet = 0
}

// AddChips adds chips to the player (e.g., winning pot)
func (p *Player) AddChips(amount int) {
	p.chips += amount
}

// SetHand sets the player's hand
func (p *Player) SetHand(hand card.Hand) {
	p.hand = hand
}

// SetPosition sets the player's position
func (p *Player) SetPosition(position int) {
	p.position = position
}

// SetStatus sets the player's status
func (p *Player) SetStatus(status PlayerStatus) {
	p.status = status
}

// String returns string representation
func (p *Player) String() string {
	return p.nickname.String()
}
