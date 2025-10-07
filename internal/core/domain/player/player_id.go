// Package player provides player domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/player/vo/PlayerId.java
package player

import "github.com/google/uuid"

// PlayerId is a type-safe player identifier (Value Object)
type PlayerId struct {
	value uuid.UUID
}

// NewPlayerId creates a new PlayerId from UUID
func NewPlayerId(id uuid.UUID) PlayerId {
	return PlayerId{value: id}
}

// NewPlayerIdFromString creates a PlayerId from string
func NewPlayerIdFromString(idStr string) (PlayerId, error) {
	// TODO: parse UUID and validate
	return PlayerId{}, nil
}

// GeneratePlayerId generates a new random PlayerId
func GeneratePlayerId() PlayerId {
	return PlayerId{value: uuid.New()}
}

// Value returns the UUID value
func (p PlayerId) Value() uuid.UUID {
	return p.value
}

// String returns the string representation
func (p PlayerId) String() string {
	return p.value.String()
}

// Equals checks if two PlayerIds are equal
func (p PlayerId) Equals(other PlayerId) bool {
	return p.value == other.value
}
