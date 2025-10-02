package identity

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// GetOrCreateUUID loads existing UUID from file or creates new one
func GetOrCreateUUID() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return uuid.New().String()
	}

	configDir := filepath.Join(homeDir, ".pokerhole")
	uuidFile := filepath.Join(configDir, "uuid")

	// Try to read existing UUID
	if data, err := os.ReadFile(uuidFile); err == nil {
		return string(data)
	}

	// Create new UUID
	newUUID := uuid.New().String()

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err == nil {
		// Save UUID to file
		_ = os.WriteFile(uuidFile, []byte(newUUID), 0644)
	}

	return newUUID
}

// GenerateNickname generates a random nickname
func GenerateNickname() string {
	adjectives := []string{
		"Lucky", "Brave", "Smart", "Quick", "Wild",
		"Cool", "Swift", "Bold", "Wise", "Epic",
		"Ace", "King", "Royal", "Grand", "Elite",
	}
	nouns := []string{
		"Shark", "Tiger", "Eagle", "Wolf", "Fox",
		"Dragon", "Phoenix", "Lion", "Bear", "Hawk",
		"Dealer", "Player", "Master", "Champion", "Legend",
	}

	adjIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	nounIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nouns))))
	numIdx, _ := rand.Int(rand.Reader, big.NewInt(1000))

	return fmt.Sprintf("%s%s%d",
		adjectives[adjIdx.Int64()],
		nouns[nounIdx.Int64()],
		numIdx.Int64())
}
