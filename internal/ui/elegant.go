package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Classic Casino Color Palette
var (
	// Casino felt green
	ColorCasinoGreen     = lipgloss.Color("#0B6623")
	ColorCasinoGreenDark = lipgloss.Color("#064018")

	// Elegant gold
	ColorGold       = lipgloss.Color("#D4AF37")
	ColorGoldBright = lipgloss.Color("#FFD700")
	ColorGoldDim    = lipgloss.Color("#B8960F")

	// Refined whites
	ColorCream = lipgloss.Color("#F5F5DC")
	ColorIvory = lipgloss.Color("#FFFFF0")

	// Casino red (for hearts/diamonds)
	ColorCasinoRed = lipgloss.Color("#DC143C")
)

// GoldGlow creates subtle gold glow effect (elegant pulse)
// intensity: 0.0 (dim) to 1.0 (bright)
func GoldGlow(text string, intensity float64) string {
	// Clamp intensity
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}

	// Interpolate between dim gold and bright gold
	var color lipgloss.Color
	if intensity > 0.7 {
		color = ColorGoldBright
	} else if intensity > 0.4 {
		color = ColorGold
	} else {
		color = ColorGoldDim
	}

	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	return style.Render(text)
}

// SubtleFade creates gentle fade in/out effect
// phase: animation phase (0-100)
func SubtleFade(text string, phase int) string {
	// Calculate opacity using sine wave for smooth fade
	normalizedPhase := float64(phase%100) / 100.0
	opacity := (math.Sin(normalizedPhase*2*math.Pi) + 1) / 2 // 0.0 to 1.0

	var style lipgloss.Style
	if opacity > 0.6 {
		style = lipgloss.NewStyle().Foreground(ColorCream).Bold(true)
	} else if opacity > 0.3 {
		style = lipgloss.NewStyle().Foreground(ColorCream)
	} else {
		style = lipgloss.NewStyle().Foreground(ColorTextSecondary)
	}

	return style.Render(text)
}

// PulseGold creates slow, elegant gold pulse
// tick: animation tick counter
func PulseGold(text string, tick int) string {
	// Slow pulse cycle (0-120 ticks for full cycle)
	cycle := tick % 120
	intensity := (math.Sin(float64(cycle)/120.0*2*math.Pi) + 1) / 2

	return GoldGlow(text, intensity)
}

// ElegantBorder creates simple gold border line
func ElegantBorder(width int) string {
	borderChar := "─"
	line := strings.Repeat(borderChar, width)

	style := lipgloss.NewStyle().Foreground(ColorGold)
	return style.Render(line)
}

// DoubleElegantBorder creates refined double border
func DoubleElegantBorder(width int) string {
	borderChar := "═"
	line := strings.Repeat(borderChar, width)

	style := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
	return style.Render(line)
}

// GoldAccent adds gold accent to specific characters (first letter, brackets, etc)
func GoldAccent(text string) string {
	result := ""
	inBracket := false

	for i, char := range text {
		charStr := string(char)

		// First letter is gold
		if i == 0 {
			style := lipgloss.NewStyle().Foreground(ColorGold).Bold(true)
			result += style.Render(charStr)
		} else if charStr == "[" || charStr == "]" {
			// Brackets are gold
			style := lipgloss.NewStyle().Foreground(ColorGold)
			result += style.Render(charStr)
			if charStr == "[" {
				inBracket = true
			} else {
				inBracket = false
			}
		} else if inBracket {
			// Content inside brackets is gold
			style := lipgloss.NewStyle().Foreground(ColorGoldBright).Bold(true)
			result += style.Render(charStr)
		} else {
			// Regular text is cream
			style := lipgloss.NewStyle().Foreground(ColorCream)
			result += style.Render(charStr)
		}
	}

	return result
}

// SoftHighlight creates subtle highlight for selected items
func SoftHighlight(text string, isActive bool) string {
	if isActive {
		style := lipgloss.NewStyle().
			Foreground(ColorGoldBright).
			Bold(true)
		return style.Render(text)
	}

	style := lipgloss.NewStyle().Foreground(ColorCream)
	return style.Render(text)
}

// SpacedTitle creates elegant spaced-out title
func SpacedTitle(text string) string {
	spaced := ""
	for i, char := range text {
		if i > 0 && char != ' ' {
			spaced += " "
		}
		spaced += string(char)
	}
	return spaced
}

// ElegantBox creates a refined box with gold border
func ElegantBox(content string, width int, isHighlighted bool) string {
	var borderColor lipgloss.Color
	if isHighlighted {
		borderColor = ColorGoldBright
	} else {
		borderColor = ColorGold
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return style.Render(content)
}

// SoftGlow creates very subtle glow effect (for active elements)
func SoftGlow(text string, isActive bool, tick int) string {
	if !isActive {
		style := lipgloss.NewStyle().Foreground(ColorCream)
		return style.Render(text)
	}

	// Very subtle pulse for active elements
	cycle := tick % 60
	intensity := 0.5 + 0.3*(math.Sin(float64(cycle)/60.0*2*math.Pi)+1)/2 // 0.5 to 0.8

	var color lipgloss.Color
	if intensity > 0.7 {
		color = ColorGoldBright
	} else {
		color = ColorGold
	}

	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	return style.Render(text)
}

// CardSuitColor returns appropriate color for card suit
func CardSuitColor(suit string) lipgloss.Color {
	if suit == "♥" || suit == "♦" {
		return ColorCasinoRed
	}
	return ColorCream
}

// MoneyGlow creates subtle glow for money/chips display
func MoneyGlow(text string, tick int) string {
	// Very slow, subtle pulse
	cycle := tick % 90
	intensity := 0.6 + 0.2*(math.Sin(float64(cycle)/90.0*2*math.Pi)+1)/2 // 0.6 to 0.8

	var color lipgloss.Color
	if intensity > 0.75 {
		color = ColorGold
	} else {
		color = ColorGoldDim
	}

	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(text)
}
