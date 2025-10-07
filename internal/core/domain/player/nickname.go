// Package player provides player domain types
// Mirror of: pokerhole-server/src/main/java/dev/xiyo/pokerhole/core/domain/player/vo/Nickname.java
package player

import "errors"

// Nickname represents a player's nickname (Value Object)
type Nickname struct {
	value string
}

const (
	MinNicknameLength = 2
	MaxNicknameLength = 20
)

// NewNickname creates a new Nickname
func NewNickname(value string) (Nickname, error) {
	// TODO: validate length and characters
	if len(value) < MinNicknameLength || len(value) > MaxNicknameLength {
		return Nickname{}, errors.New("invalid nickname length")
	}
	return Nickname{value: value}, nil
}

// Value returns the string value
func (n Nickname) Value() string {
	return n.value
}

// String returns the string representation
func (n Nickname) String() string {
	return n.value
}

// Equals checks if two Nicknames are equal
func (n Nickname) Equals(other Nickname) bool {
	return n.value == other.value
}
