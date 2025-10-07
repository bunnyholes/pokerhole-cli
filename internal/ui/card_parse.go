package ui

import (
	"strings"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
)

func parseHand(handStr string) []card.Card {
	handStr = strings.TrimSpace(handStr)
	if handStr == "" || handStr == "[]" {
		return nil
	}

	trimmed := strings.Trim(handStr, "[]")
	if trimmed == "" {
		return nil
	}

	pieces := strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
	var cards []card.Card
	for _, piece := range pieces {
		if c := parseCardString(strings.TrimSpace(piece)); c != nil {
			cards = append(cards, *c)
		}
	}

	return cards
}
