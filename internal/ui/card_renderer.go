package ui

import (
	"fmt"

	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/charmbracelet/lipgloss"
)

// 카드 색상 정의
var (
	redSuitColor   = lipgloss.Color("#FF6B6B")
	blackSuitColor = lipgloss.Color("#2C3E50")
	cardBgColor    = lipgloss.Color("#FFFFFF")
	cardBorder     = lipgloss.Color("#3498DB") // 파란색 테두리
	cardBackColor  = lipgloss.Color("#34495E") // 뒷면 색상
)

// 카드 스타일 (3x3 디자인: 5칸 width, 3칸 height)
var cardBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(cardBorder).
	Background(cardBgColor).
	Padding(0, 1) // 좌우 1칸 패딩

var backCardBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(cardBorder).
	Background(cardBackColor).
	Padding(0, 1)

// 카드 뒷면 디자인 (3x3) - Lipgloss 스타일 사용
func renderCardBack() string {
	cardWidth := 5

	patternStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7F8C8D")).
		Bold(true)

	// 줄 스타일 - 뒷면 배경색 포함
	lineStyle := lipgloss.NewStyle().
		Background(cardBackColor).
		Width(cardWidth)

	// 상단 빈 줄
	topLine := lineStyle.Copy().
		Render(" ")

	// 중앙 패턴
	middleLine := lineStyle.Copy().
		Align(lipgloss.Center).
		Render(patternStyle.Render("░░░"))

	// 하단 빈 줄
	bottomLine := lineStyle.Copy().
		Render(" ")

	backContent := lipgloss.JoinVertical(
		lipgloss.Left,
		topLine,
		middleLine,
		bottomLine,
	)

	return backCardBorderStyle.Render(backContent)
}

// 카드 앞면 렌더링 (3x3 디자인) - Lipgloss 스타일 사용
func renderCard(c card.Card) string {
	rank := getRankSymbol(c.Rank())
	suit := getSuitSymbol(c.Suit())
	suitColor := getSuitColor(c.Suit())

	// 스타일 정의
	suitStyle := lipgloss.NewStyle().
		Foreground(suitColor).
		Bold(true)

	rankStyle := lipgloss.NewStyle().
		Foreground(suitColor).
		Bold(true)

	// 카드 너비 결정 (10은 2자리이므로 4칸, 나머지는 5칸)
	var cardWidth int
	if len(rank) == 2 {
		cardWidth = 4
	} else {
		cardWidth = 5
	}

	// 줄 스타일 - 배경색 포함
	lineStyle := lipgloss.NewStyle().
		Background(cardBgColor).
		Width(cardWidth)

	// 좌측 상단 suit (좌측 정렬)
	topLine := lineStyle.Copy().
		Align(lipgloss.Left).
		Render(suitStyle.Render(suit))

	// 중앙 rank (중앙 정렬)
	middleLine := lineStyle.Copy().
		Align(lipgloss.Center).
		Render(rankStyle.Render(rank))

	// 우측 하단 suit (우측 정렬)
	bottomLine := lineStyle.Copy().
		Align(lipgloss.Right).
		Render(suitStyle.Render(suit))

	// 3줄을 수직으로 결합
	cardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		topLine,
		middleLine,
		bottomLine,
	)

	return cardBorderStyle.Render(cardContent)
}

// 빈 카드 슬롯 (3x3 디자인) - Lipgloss 스타일 사용
func renderEmptyCardSlot() string {
	cardWidth := 5

	questionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BDC3C7")).
		Bold(true)

	emptyBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#95A5A6")).
		Padding(0, 1)

	// 줄 스타일 - 밝은 회색 배경
	lineStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#ECF0F1")).
		Width(cardWidth)

	// 상단 빈 줄
	topLine := lineStyle.Copy().
		Render(" ")

	// 중앙 물음표
	middleLine := lineStyle.Copy().
		Align(lipgloss.Center).
		Render(questionStyle.Render("?"))

	// 하단 빈 줄
	bottomLine := lineStyle.Copy().
		Render(" ")

	emptyContent := lipgloss.JoinVertical(
		lipgloss.Left,
		topLine,
		middleLine,
		bottomLine,
	)

	return emptyBorderStyle.Render(emptyContent)
}

// 랭크 심볼
func getRankSymbol(r card.Rank) string {
	switch r {
	case card.Ace:
		return "A"
	case card.King:
		return "K"
	case card.Queen:
		return "Q"
	case card.Jack:
		return "J"
	case card.Ten:
		return "10"
	default:
		return fmt.Sprintf("%d", r.Value())
	}
}

// 슈트 심볼 (유니코드)
func getSuitSymbol(s card.Suit) string {
	switch s {
	case card.Hearts:
		return "♥"
	case card.Diamonds:
		return "♦"
	case card.Clubs:
		return "♣"
	case card.Spades:
		return "♠"
	default:
		return "?"
	}
}

// 슈트 색상
func getSuitColor(s card.Suit) lipgloss.Color {
	switch s {
	case card.Hearts, card.Diamonds:
		return redSuitColor
	case card.Clubs, card.Spades:
		return blackSuitColor
	default:
		return lipgloss.Color("#95A5A6")
	}
}

// 카드 리스트를 가로로 배치
func renderCardsHorizontal(cards []card.Card, hideCards bool, maxCards int) string {
	renderedCards := []string{}

	for i := 0; i < maxCards; i++ {
		if i < len(cards) {
			if hideCards {
				renderedCards = append(renderedCards, renderCardBack())
			} else {
				renderedCards = append(renderedCards, renderCard(cards[i]))
			}
		} else {
			renderedCards = append(renderedCards, renderEmptyCardSlot())
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedCards...)
}

// 커뮤니티 카드 렌더링 (5장)
func renderCommunityCards(cardStrings []string) string {
	cards := make([]string, 5)

	for i := 0; i < 5; i++ {
		if i < len(cardStrings) && cardStrings[i] != "" {
			// 문자열에서 카드 파싱
			c := parseCardString(cardStrings[i])
			if c != nil {
				cards[i] = renderCard(*c)
			} else {
				cards[i] = renderEmptyCardSlot()
			}
		} else {
			cards[i] = renderEmptyCardSlot()
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
}

// 카드 문자열 파싱 (예: "♠A" 또는 "A♠")
func parseCardString(s string) *card.Card {
	if s == "" {
		return nil
	}

	// 룬(rune) 배열로 변환하여 유니코드 문자 올바르게 처리
	runes := []rune(s)
	if len(runes) < 2 {
		return nil
	}

	// 첫 번째 룬이 수트인지 확인 (♠, ♥, ♦, ♣)
	firstRune := runes[0]
	var suitStr, rankStr string

	if firstRune == '♥' || firstRune == '♦' || firstRune == '♣' || firstRune == '♠' {
		// 수트가 앞에 있는 형식: "♠A"
		suitStr = string(firstRune)
		rankStr = string(runes[1:])
	} else {
		// 랭크가 앞에 있는 형식: "A♠"
		suitStr = string(runes[len(runes)-1])
		rankStr = string(runes[:len(runes)-1])
	}

	var rank card.Rank
	switch rankStr {
	case "A":
		rank = card.Ace
	case "K":
		rank = card.King
	case "Q":
		rank = card.Queen
	case "J":
		rank = card.Jack
	case "10":
		rank = card.Ten
	case "9":
		rank = card.Nine
	case "8":
		rank = card.Eight
	case "7":
		rank = card.Seven
	case "6":
		rank = card.Six
	case "5":
		rank = card.Five
	case "4":
		rank = card.Four
	case "3":
		rank = card.Three
	case "2":
		rank = card.Two
	default:
		return nil
	}

	var suit card.Suit
	switch suitStr {
	case "♥":
		suit = card.Hearts
	case "♦":
		suit = card.Diamonds
	case "♣":
		suit = card.Clubs
	case "♠":
		suit = card.Spades
	default:
		return nil
	}

	c, err := card.NewCard(suit, rank)
	if err != nil {
		return nil
	}
	return &c
}

// renderHandCards renders a hand (2 cards) horizontally
func renderHandCards(handStr string) string {
	cards := parseHand(handStr)
	if len(cards) == 0 {
		return renderEmptyCardSlot() + " " + renderEmptyCardSlot()
	}

	renderedCards := []string{}
	for _, c := range cards {
		renderedCards = append(renderedCards, renderCard(c))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedCards...)
}

// renderCommunityCardsLarge renders community cards (up to 5) horizontally
func renderCommunityCardsLarge(cardStrings []string) string {
	cards := make([]string, 5)

	for i := 0; i < 5; i++ {
		if i < len(cardStrings) && cardStrings[i] != "" {
			c := parseCardString(cardStrings[i])
			if c != nil {
				cards[i] = renderCard(*c)
			} else {
				cards[i] = renderEmptyCardSlot()
			}
		} else {
			cards[i] = renderEmptyCardSlot()
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
}

// renderRankCards renders only the rank cards (no empty slots)
func renderRankCards(cardStrings []string) string {
	if len(cardStrings) == 0 {
		return ""
	}

	renderedCards := []string{}
	for _, cardStr := range cardStrings {
		if cardStr == "" {
			continue
		}
		c := parseCardString(cardStr)
		if c != nil {
			renderedCards = append(renderedCards, renderCard(*c))
		}
	}

	if len(renderedCards) == 0 {
		return ""
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedCards...)
}

// 플레이어 정보 한 줄 (미니)
func renderPlayerBox(name string, chips int, bet int, status string, cards []card.Card, hideCards bool, isActive bool) string {
	// 카드 렌더링
	cardDisplay := renderCardsHorizontal(cards, hideCards, 2)

	// 이름 + 정보 (한 줄)
	chipsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F39C12"))
	betStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C"))

	nameStyle := lipgloss.NewStyle().Bold(true)
	if isActive {
		nameStyle = nameStyle.Foreground(lipgloss.Color("#2ECC71"))
		name = "▶" + name
	}

	// Build info string with status
	infoStr := fmt.Sprintf("%s %s %s",
		nameStyle.Render(name),
		chipsStyle.Render(fmt.Sprintf("칩%d", chips)),
		betStyle.Render(fmt.Sprintf("베팅%d", bet)),
	)

	// Add all-in indicator
	if status == "ALL_IN" {
		allInStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			Bold(true)
		infoStr += " " + allInStyle.Render("[올인]")
	}

	// 전체 조합 (가로로)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		cardDisplay,
		lipgloss.NewStyle().Padding(1, 0, 0, 1).Render(infoStr),
	)
}
