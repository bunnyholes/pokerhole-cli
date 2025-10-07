package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Gothic ornamental characters
const (
	// Corner decorations
	CornerTopLeft     = "╔"
	CornerTopRight    = "╗"
	CornerBottomLeft  = "╚"
	CornerBottomRight = "╝"

	// Ornamental symbols
	OrnamentDiamond  = "◆"
	OrnamentStar     = "✦"
	OrnamentFleur    = "❈"
	OrnamentCross    = "✤"
	OrnamentSnowflake = "❉"
	OrnamentClub     = "♣"
	OrnamentSpade    = "♠"
	OrnamentHeart    = "♥"
	OrnamentDiamond2 = "♦"
)

// GothicFrame creates ornate gothic-style frame
func GothicFrame(content string, width int, height int) string {
	var lines []string

	// Top ornamental border
	topBorder := GothicTopBorder(width)
	lines = append(lines, topBorder)

	// Content lines (centered)
	contentLines := strings.Split(content, "\n")
	for _, line := range contentLines {
		framedLine := GothicSideBorders(line, width)
		lines = append(lines, framedLine)
	}

	// Fill remaining height with empty framed lines
	for len(lines) < height-1 {
		framedLine := GothicSideBorders("", width)
		lines = append(lines, framedLine)
	}

	// Bottom ornamental border
	bottomBorder := GothicBottomBorder(width)
	lines = append(lines, bottomBorder)

	return strings.Join(lines, "\n")
}

// GothicTopBorder creates ornate top border
func GothicTopBorder(width int) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)

	// ╔═══❈═══╗ pattern
	innerWidth := width - 2
	if innerWidth < 10 {
		innerWidth = 10
	}

	// Create pattern with ornaments
	pattern := strings.Repeat("═", (innerWidth-3)/2) + "❈" + strings.Repeat("═", (innerWidth-3)/2)

	// Ensure exact width
	for len(pattern) < innerWidth {
		pattern += "═"
	}
	if len(pattern) > innerWidth {
		pattern = pattern[:innerWidth]
	}

	border := CornerTopLeft + pattern + CornerTopRight
	return goldStyle.Render(border)
}

// GothicBottomBorder creates ornate bottom border
func GothicBottomBorder(width int) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)

	innerWidth := width - 2
	if innerWidth < 10 {
		innerWidth = 10
	}

	pattern := strings.Repeat("═", (innerWidth-3)/2) + "❈" + strings.Repeat("═", (innerWidth-3)/2)

	for len(pattern) < innerWidth {
		pattern += "═"
	}
	if len(pattern) > innerWidth {
		pattern = pattern[:innerWidth]
	}

	border := CornerBottomLeft + pattern + CornerBottomRight
	return goldStyle.Render(border)
}

// GothicSideBorders adds side borders with ornaments
func GothicSideBorders(content string, width int) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)

	innerWidth := width - 2

	// Pad or trim content to fit
	paddedContent := content
	contentLen := lipgloss.Width(content)

	if contentLen < innerWidth {
		// Center the content
		leftPad := (innerWidth - contentLen) / 2
		rightPad := innerWidth - contentLen - leftPad
		paddedContent = strings.Repeat(" ", leftPad) + content + strings.Repeat(" ", rightPad)
	} else if contentLen > innerWidth {
		paddedContent = content[:innerWidth]
	}

	return goldStyle.Render("║") + paddedContent + goldStyle.Render("║")
}

// OrnamentalDivider creates decorative divider line
func OrnamentalDivider(width int) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold)
	dimGoldStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)

	// Pattern: ─❈─✦─◆─✦─❈─
	ornaments := []string{"❈", "✦", "◆", "✦"}

	var parts []string
	remaining := width
	ornamentIndex := 0

	for remaining > 0 {
		if remaining >= 3 {
			parts = append(parts, dimGoldStyle.Render("─"))
			parts = append(parts, goldStyle.Render(ornaments[ornamentIndex%len(ornaments)]))
			parts = append(parts, dimGoldStyle.Render("─"))
			remaining -= 3
			ornamentIndex++
		} else {
			parts = append(parts, dimGoldStyle.Render(strings.Repeat("─", remaining)))
			remaining = 0
		}
	}

	return strings.Join(parts, "")
}

// VintageCardArt creates ASCII art of a playing card
func VintageCardArt(suit string, rank string, faceDown bool) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold)
	dimGoldStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)
	redStyle := lipgloss.NewStyle().Foreground(ColorCasinoRed)
	creamStyle := lipgloss.NewStyle().Foreground(ColorCream)

	if faceDown {
		// Ornate card back
		return strings.Join([]string{
			goldStyle.Render("┌─────┐"),
			goldStyle.Render("│") + dimGoldStyle.Render("▓▓▓") + goldStyle.Render("│"),
			goldStyle.Render("│") + dimGoldStyle.Render("▓❈▓") + goldStyle.Render("│"),
			goldStyle.Render("│") + dimGoldStyle.Render("▓▓▓") + goldStyle.Render("│"),
			goldStyle.Render("└─────┘"),
		}, "\n")
	}

	// Face up card
	suitStyle := creamStyle
	if suit == "♥" || suit == "♦" {
		suitStyle = redStyle
	}

	// Format rank (right-align for 10)
	displayRank := rank
	if len(rank) == 1 {
		displayRank = " " + rank
	}

	return strings.Join([]string{
		goldStyle.Render("┌─────┐"),
		goldStyle.Render("│") + suitStyle.Render(suit) + creamStyle.Render(displayRank) + "  " + goldStyle.Render("│"),
		goldStyle.Render("│  ") + suitStyle.Render(suit) + "  " + goldStyle.Render("│"),
		goldStyle.Render("│  ") + creamStyle.Render(displayRank) + suitStyle.Render(suit) + goldStyle.Render("│"),
		goldStyle.Render("└─────┘"),
	}, "\n")
}

// ChipStackArt creates ASCII art of poker chips
func ChipStackArt() string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold)
	dimStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)

	return strings.Join([]string{
		"   " + goldStyle.Render("╱") + dimStyle.Render("▀▀") + goldStyle.Render("╲"),
		"  " + goldStyle.Render("│") + dimStyle.Render("▓▓▓") + goldStyle.Render("│"),
		"  " + goldStyle.Render("│") + dimStyle.Render("▓▓▓") + goldStyle.Render("│"),
		"  " + goldStyle.Render("╲") + dimStyle.Render("▄▄") + goldStyle.Render("╱"),
	}, "\n")
}

// PokerTableTopView creates ASCII art of poker table from above
func PokerTableTopView() string {
	greenStyle := lipgloss.NewStyle().Foreground(ColorCasinoGreen).Bold(true)
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold)

	return strings.Join([]string{
		"        " + goldStyle.Render("╔═══════════════════════════════╗"),
		"      " + goldStyle.Render("╔═╝") + greenStyle.Render("░░░░░░░░░░░░░░░░░░░░░░░░░") + goldStyle.Render("╚═╗"),
		"     " + goldStyle.Render("╔╝") + greenStyle.Render("░░░░░░░░░░░░░░░░░░░░░░░░░░░") + goldStyle.Render("╚╗"),
		"    " + goldStyle.Render("╔╝") + greenStyle.Render("░░░") + goldStyle.Render("T E X A S   H O L D ' E M") + greenStyle.Render("░░░") + goldStyle.Render("╚╗"),
		"    " + goldStyle.Render("║") + greenStyle.Render("░░░░░░░░░░░░░░░░░░░░░░░░░░░░░") + goldStyle.Render("║"),
		"    " + goldStyle.Render("║") + greenStyle.Render("░░░░░░░░") + goldStyle.Render("[ FELT ]") + greenStyle.Render("░░░░░░░░") + goldStyle.Render("║"),
		"     " + goldStyle.Render("╚╗") + greenStyle.Render("░░░░░░░░░░░░░░░░░░░░░░░░░░░") + goldStyle.Render("╔╝"),
		"      " + goldStyle.Render("╚═╗") + greenStyle.Render("░░░░░░░░░░░░░░░░░░░░░░░░░") + goldStyle.Render("╔═╝"),
		"        " + goldStyle.Render("╚═══════════════════════════════╝"),
	}, "\n")
}

// OrnateTitle creates ornate title with decorations
func OrnateTitle(text string, tick int) string {
	// Title with ornamental frame
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)

	// Cycling ornaments
	ornaments := []string{"❈", "✦", "◆", "❉"}
	leftOrn := ornaments[tick%len(ornaments)]
	rightOrn := ornaments[(tick+2)%len(ornaments)]

	// Pulsing title
	titleText := PulseGold(text, tick)

	leftDeco := goldStyle.Render(leftOrn) + dimStyle.Render("═══") + goldStyle.Render("❈") + dimStyle.Render("═══") + goldStyle.Render(leftOrn)
	rightDeco := goldStyle.Render(rightOrn) + dimStyle.Render("═══") + goldStyle.Render("❈") + dimStyle.Render("═══") + goldStyle.Render(rightOrn)

	return leftDeco + "  " + titleText + "  " + rightDeco
}

// FeltBackground creates felt texture pattern
func FeltBackground(width int) string {
	greenStyle := lipgloss.NewStyle().Foreground(ColorCasinoGreen)
	darkGreenStyle := lipgloss.NewStyle().Foreground(ColorCasinoGreenDark)

	// Alternating pattern for felt texture
	pattern := ""
	for i := 0; i < width; i++ {
		if i%4 == 0 {
			pattern += greenStyle.Render("░")
		} else if i%4 == 2 {
			pattern += darkGreenStyle.Render("▒")
		} else {
			pattern += greenStyle.Render("░")
		}
	}

	return pattern
}

// VintagePokerLogo creates ASCII art logo
func VintagePokerLogo() string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)

	return strings.Join([]string{
		"",
		"        " + goldStyle.Render("╔═══════════════════════════════════╗"),
		"        " + goldStyle.Render("║") + "   " + dimStyle.Render("♠ ♥ ♦ ♣") + "  " + goldStyle.Render("POKERHOLE") + "  " + dimStyle.Render("♣ ♦ ♥ ♠") + "   " + goldStyle.Render("║"),
		"        " + goldStyle.Render("║") + "         " + dimStyle.Render("Est. 2025") + " " + goldStyle.Render("◆") + " " + dimStyle.Render("Texas Hold'em") + "        " + goldStyle.Render("║"),
		"        " + goldStyle.Render("╚═══════════════════════════════════╝"),
		"",
	}, "\n")
}

// CornerOrnament creates corner decorative element
func CornerOrnament(position string) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold)

	switch position {
	case "top-left":
		return goldStyle.Render("╔═❈")
	case "top-right":
		return goldStyle.Render("❈═╗")
	case "bottom-left":
		return goldStyle.Render("╚═❈")
	case "bottom-right":
		return goldStyle.Render("❈═╝")
	default:
		return ""
	}
}

// OrnamentalSeparator creates separator with suit symbols
func OrnamentalSeparator(width int) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold)
	dimStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)

	// Pattern: ─♠─♥─♦─♣─
	suits := []string{"♠", "♥", "♦", "♣"}

	var parts []string
	remaining := width
	suitIndex := 0

	for remaining > 0 {
		if remaining >= 3 {
			parts = append(parts, dimStyle.Render("─"))
			parts = append(parts, goldStyle.Render(suits[suitIndex%len(suits)]))
			parts = append(parts, dimStyle.Render("─"))
			remaining -= 3
			suitIndex++
		} else {
			parts = append(parts, dimStyle.Render(strings.Repeat("─", remaining)))
			remaining = 0
		}
	}

	return strings.Join(parts, "")
}

// VintageMoneyBadge creates ornate money display
func VintageMoneyBadge(label string, amount int, tick int) string {
	goldStyle := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(ColorGoldDim)

	amountText := MoneyGlow(fmt.Sprintf("%d", amount), tick)

	leftBracket := goldStyle.Render("〔")
	rightBracket := goldStyle.Render("〕")
	labelText := dimStyle.Render(label + ":")

	return leftBracket + " " + labelText + " " + amountText + " " + rightBracket
}
