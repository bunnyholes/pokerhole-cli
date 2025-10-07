package ui

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

// GradientText creates flowing gradient text effect
func GradientText(text string, offset int) string {
	// Premium color palette
	colors := []lipgloss.Color{
		lipgloss.Color("#FFD700"), // Gold
		lipgloss.Color("#FFA500"), // Orange
		lipgloss.Color("#FF69B4"), // Hot Pink
		lipgloss.Color("#9370DB"), // Medium Purple
		lipgloss.Color("#4169E1"), // Royal Blue
		lipgloss.Color("#00CED1"), // Dark Turquoise
		lipgloss.Color("#00FF00"), // Lime
		lipgloss.Color("#FFD700"), // Gold (loop)
	}

	result := ""
	for i, char := range text {
		colorIndex := (i + offset) % len(colors)
		style := lipgloss.NewStyle().Foreground(colors[colorIndex])
		result += style.Render(string(char))
	}
	return result
}

// WaveText creates wave-like color animation
func WaveText(text string, offset int) string {
	result := ""
	for i, char := range text {
		// Calculate wave position
		wave := math.Sin(float64(i+offset) * 0.5)

		// Map wave to color
		var color lipgloss.Color
		if wave > 0.5 {
			color = ColorAccentGold
		} else if wave > 0 {
			color = ColorAccentPurple
		} else if wave > -0.5 {
			color = ColorAccentBlue
		} else {
			color = ColorAccentGreen
		}

		style := lipgloss.NewStyle().Foreground(color).Bold(true)
		result += style.Render(string(char))
	}
	return result
}

// RainbowBorder creates animated rainbow border
func RainbowBorder(content string, offset int) string {
	colors := []lipgloss.Color{
		lipgloss.Color("#FF0000"), // Red
		lipgloss.Color("#FF7F00"), // Orange
		lipgloss.Color("#FFFF00"), // Yellow
		lipgloss.Color("#00FF00"), // Green
		lipgloss.Color("#0000FF"), // Blue
		lipgloss.Color("#4B0082"), // Indigo
		lipgloss.Color("#9400D3"), // Violet
	}

	borderColor := colors[offset%len(colors)]

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Align(lipgloss.Center)

	return style.Render(content)
}

// GlowText creates glowing text effect
func GlowText(text string, intensity bool) string {
	var style lipgloss.Style
	if intensity {
		// Bright glow
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Underline(true)
	} else {
		// Dim glow
		style = lipgloss.NewStyle().
			Foreground(ColorAccentGold).
			Bold(true)
	}
	return style.Render(text)
}

// PremiumTitle creates premium-looking title with effects
func PremiumTitle(text string, offset int) string {
	// Create gradient text
	gradientText := GradientText(text, offset)

	// Add decorative elements
	decorator := "═══"
	decoratorStyle := lipgloss.NewStyle().Foreground(ColorAccentGold)

	result := decoratorStyle.Render(decorator) + " " + gradientText + " " + decoratorStyle.Render(decorator)
	return result
}

// FlowingGradientLine creates flowing gradient horizontal line
func FlowingGradientLine(width int, offset int) string {
	chars := []string{"═", "─", "━", "═"}
	colors := []lipgloss.Color{
		ColorAccentGold,
		ColorAccentPurple,
		ColorAccentBlue,
		ColorAccentGreen,
	}

	result := ""
	for i := 0; i < width; i++ {
		charIndex := (i + offset) % len(chars)
		colorIndex := (i + offset) % len(colors)

		style := lipgloss.NewStyle().Foreground(colors[colorIndex])
		result += style.Render(chars[charIndex])
	}
	return result
}

// ShimmerText creates shimmering text effect
func ShimmerText(text string, phase int) string {
	// Create shimmer effect by highlighting different characters
	result := ""
	for i, char := range text {
		var style lipgloss.Style
		shimmerPos := (phase + i) % len(text)

		if shimmerPos < len(text)/3 {
			// Bright
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
		} else if shimmerPos < 2*len(text)/3 {
			// Medium
			style = lipgloss.NewStyle().
				Foreground(ColorAccentGold).
				Bold(true)
		} else {
			// Dim
			style = lipgloss.NewStyle().
				Foreground(ColorTextSecondary)
		}

		result += style.Render(string(char))
	}
	return result
}

// PulseText creates pulsing text effect
func PulseText(text string, pulse bool) string {
	var style lipgloss.Style
	if pulse {
		style = lipgloss.NewStyle().
			Foreground(ColorAccentGold).
			Bold(true).
			Underline(true)
	} else {
		style = lipgloss.NewStyle().
			Foreground(ColorAccentGold).
			Bold(true)
	}
	return style.Render(text)
}

// GradientBox creates box with gradient border
func GradientBox(content string, offset int, width int) string {
	// Side borders with cycling colors
	colors := []lipgloss.Color{
		ColorAccentGold,
		ColorAccentPurple,
		ColorAccentBlue,
	}
	borderColor := colors[offset%len(colors)]

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Align(lipgloss.Center).
		Padding(1, 2)

	return style.Render(content)
}
