package ui

import (
	"fmt"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/charmbracelet/lipgloss"
)

// renderCardCompact renders a single card in compact inline format
func renderCardCompact(c card.Card) string {
	rank := getRankSymbol(c.Rank())
	suit := getSuitSymbol(c.Suit())
	suitColor := getSuitColor(c.Suit())

	// Card style: [♠A] format
	cardStyle := lipgloss.NewStyle().
		Foreground(suitColor).
		Bold(true)

	return cardStyle.Render(fmt.Sprintf("[%s%s]", suit, rank))
}

// renderCardBackCompact renders card back in compact format
func renderCardBackCompact() string {
	backStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7F8C8D")).
		Bold(true)

	return backStyle.Render("[??]")
}

// renderHandCompact renders 2 hole cards inline
func renderHandCompact(cards []card.Card, hideCards bool) string {
	if hideCards || len(cards) == 0 {
		return renderCardBackCompact() + " " + renderCardBackCompact()
	}

	if len(cards) == 1 {
		return renderCardCompact(cards[0]) + " " + renderCardBackCompact()
	}

	return renderCardCompact(cards[0]) + " " + renderCardCompact(cards[1])
}

// renderCommunityCardsCompact renders 5 community cards inline
func renderCommunityCardsCompact(cardStrings []string) string {
	cards := make([]string, 5)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5F6368"))

	for i := 0; i < 5; i++ {
		if i < len(cardStrings) && cardStrings[i] != "" {
			c := parseCardString(cardStrings[i])
			if c != nil {
				cards[i] = renderCardCompact(*c)
			} else {
				cards[i] = emptyStyle.Render("[--]")
			}
		} else {
			cards[i] = emptyStyle.Render("[--]")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cards[0], " ", cards[1], " ", cards[2], " ", cards[3], " ", cards[4])
}

// renderRankCardsCompact renders rank cards inline
func renderRankCardsCompact(cardStrings []string) string {
	if len(cardStrings) == 0 {
		return ""
	}

	var cards []string
	for _, cardStr := range cardStrings {
		if cardStr == "" {
			continue
		}
		c := parseCardString(cardStr)
		if c != nil {
			cards = append(cards, renderCardCompact(*c))
		}
	}

	if len(cards) == 0 {
		return ""
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cards...)
}
