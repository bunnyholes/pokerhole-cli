// Package game provides game domain logic
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/game/GameId.java
package game

import "github.com/google/uuid"

// GameId is a type-safe game identifier (Value Object)
type GameId struct {
	value uuid.UUID
}

// NewGameId creates a new GameId from UUID
func NewGameId(id uuid.UUID) GameId {
	return GameId{value: id}
}

// NewGameIdFromString creates a GameId from string
func NewGameIdFromString(idStr string) (GameId, error) {
	// TODO: parse UUID and validate
	return GameId{}, nil
}

// GenerateGameId generates a new random GameId
func GenerateGameId() GameId {
	return GameId{value: uuid.New()}
}

// Value returns the UUID value
func (g GameId) Value() uuid.UUID {
	return g.value
}

// String returns the string representation
func (g GameId) String() string {
	return g.value.String()
}

// Equals checks if two GameIds are equal
func (g GameId) Equals(other GameId) bool {
	return g.value == other.value
}
