package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/bunnyholes/pokerhole/client/internal/ui/constants"
)

// Professional Design System for PokerHole
// Inspired by modern poker platforms like PokerStars and GGPoker

// Color Palette - Re-exported from constants package for convenience
var (
	// Background Colors
	ColorBgPrimary   = constants.ColorBgPrimary
	ColorBgSecondary = constants.ColorBgSecondary
	ColorBgTable     = constants.ColorBgTable
	ColorBgCard      = constants.ColorBgCard
	ColorBgCardBack  = constants.ColorBgCardBack

	// Text Colors
	ColorTextPrimary   = constants.ColorTextPrimary
	ColorTextSecondary = constants.ColorTextSecondary
	ColorTextMuted     = constants.ColorTextMuted

	// Accent Colors
	ColorAccentGold   = constants.ColorAccentGold
	ColorAccentGreen  = constants.ColorAccentGreen
	ColorAccentRed    = constants.ColorAccentRed
	ColorAccentBlue   = constants.ColorAccentBlue
	ColorAccentPurple = constants.ColorAccentPurple

	// Vintage Colors (for intro scene and retro elements)
	ColorVintageGold    = constants.ColorVintageGold
	ColorVintageGoldDim = constants.ColorVintageGoldDim

	// Card Suit Colors
	ColorSuitRed   = constants.ColorSuitRed
	ColorSuitBlack = constants.ColorSuitBlack

	// Status Colors
	ColorSuccess = constants.ColorSuccess
	ColorWarning = constants.ColorWarning
	ColorError   = constants.ColorError
	ColorInfo    = constants.ColorInfo

	// Border Colors
	ColorBorderSubtle = constants.ColorBorderSubtle
	ColorBorderNormal = constants.ColorBorderNormal
	ColorBorderStrong = constants.ColorBorderStrong
)

// Typography Styles
var (
	StyleH1 = lipgloss.NewStyle().
		Foreground(ColorTextPrimary).
		Bold(true).
		MarginBottom(1)

	StyleH2 = lipgloss.NewStyle().
		Foreground(ColorTextPrimary).
		Bold(true).
		MarginBottom(1)

	StyleH3 = lipgloss.NewStyle().
		Foreground(ColorTextSecondary).
		Bold(true)

	StyleBody = lipgloss.NewStyle().
			Foreground(ColorTextPrimary)

	StyleBodyMuted = lipgloss.NewStyle().
			Foreground(ColorTextSecondary)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Bold(true)
)

// Component Styles
var (
	// Card Styles
	StyleCard = lipgloss.NewStyle().
			Background(ColorBgCard).
			Foreground(ColorSuitBlack).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNormal).
			Padding(0, 1)

	StyleCardBack = lipgloss.NewStyle().
			Background(ColorBgCardBack).
			Foreground(ColorTextMuted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderSubtle).
			Padding(0, 1)

	// Panel Styles
	StylePanel = lipgloss.NewStyle().
			Background(ColorBgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNormal).
			Padding(1, 2)

	StylePanelHighlight = lipgloss.NewStyle().
				Background(ColorBgSecondary).
				Border(lipgloss.ThickBorder()).
				BorderForeground(ColorAccentGold).
				Padding(1, 2)

	// Button Styles
	StyleButton = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Background(ColorBgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNormal).
			Padding(0, 2).
			MarginRight(1)

	StyleButtonPrimary = lipgloss.NewStyle().
				Foreground(ColorBgPrimary).
				Background(ColorAccentGold).
				Bold(true).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorAccentGold).
				Padding(0, 3).
				MarginRight(1)

	StyleButtonActive = lipgloss.NewStyle().
				Foreground(ColorBgPrimary).
				Background(ColorAccentGreen).
				Bold(true).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorAccentGreen).
				Padding(0, 3).
				MarginRight(1)

	// Badge Styles
	StyleBadgeSuccess = lipgloss.NewStyle().
				Foreground(ColorBgPrimary).
				Background(ColorSuccess).
				Bold(true).
				Padding(0, 1)

	StyleBadgeWarning = lipgloss.NewStyle().
				Foreground(ColorBgPrimary).
				Background(ColorWarning).
				Bold(true).
				Padding(0, 1)

	StyleBadgeError = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Background(ColorError).
			Bold(true).
			Padding(0, 1)

	StyleBadgeInfo = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Background(ColorInfo).
			Bold(true).
			Padding(0, 1)

	// Chip Display
	StyleChips = lipgloss.NewStyle().
			Foreground(ColorAccentGold).
			Bold(true)

	// Bet Display
	StyleBet = lipgloss.NewStyle().
			Foreground(ColorAccentRed).
			Bold(true)

	// Pot Display
	StylePot = lipgloss.NewStyle().
			Foreground(ColorAccentGold).
			Bold(true).
			Background(ColorBgSecondary).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccentGold)
)

// Spacing Constants
const (
	SpacingXS = 1
	SpacingSM = 2
	SpacingMD = 3
	SpacingLG = 4
	SpacingXL = 6
)

// Border Styles
var (
	BorderRounded = lipgloss.RoundedBorder()
	BorderThick   = lipgloss.ThickBorder()
	BorderDouble  = lipgloss.DoubleBorder()
	BorderNormal  = lipgloss.NormalBorder()
)

// Helper Functions
func RenderTitle(text string) string {
	return StyleH1.Render(text)
}

func RenderSubtitle(text string) string {
	return StyleH2.Render(text)
}

func RenderLabel(text string) string {
	return StyleLabel.Render(text)
}

func RenderBadge(text string, variant string) string {
	switch variant {
	case "success":
		return StyleBadgeSuccess.Render(text)
	case "warning":
		return StyleBadgeWarning.Render(text)
	case "error":
		return StyleBadgeError.Render(text)
	case "info":
		return StyleBadgeInfo.Render(text)
	default:
		return StyleBadgeInfo.Render(text)
	}
}

func RenderChips(amount int) string {
	return StyleChips.Render("💰 " + lipgloss.NewStyle().Render(formatNumber(amount)))
}

func RenderBet(amount int) string {
	return StyleBet.Render("🎲 " + formatNumber(amount))
}

func RenderPot(amount int) string {
	return StylePot.Render("POT: " + formatNumber(amount))
}

func formatNumber(n int) string {
	return fmt.Sprintf("%d", n)
}
