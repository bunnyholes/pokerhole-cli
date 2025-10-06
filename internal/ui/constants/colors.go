package constants

import "github.com/charmbracelet/lipgloss"

// Color Palette - Dark Poker Theme
var (
	// Background Colors
	ColorBgPrimary   = lipgloss.Color("#0F1419") // Deep dark background
	ColorBgSecondary = lipgloss.Color("#1C252C") // Slightly lighter
	ColorBgTable     = lipgloss.Color("#1A5F3E") // Poker table green
	ColorBgCard      = lipgloss.Color("#FFFFFF") // Card background
	ColorBgCardBack  = lipgloss.Color("#2C3E50") // Card back

	// Text Colors
	ColorTextPrimary   = lipgloss.Color("#E8EAED") // High contrast text
	ColorTextSecondary = lipgloss.Color("#9AA0A6") // Muted text
	ColorTextMuted     = lipgloss.Color("#5F6368") // Very muted

	// Accent Colors
	ColorAccentGold   = lipgloss.Color("#FFB900") // Gold for chips
	ColorAccentGreen  = lipgloss.Color("#10B981") // Success/active
	ColorAccentRed    = lipgloss.Color("#EF4444") // Danger/fold
	ColorAccentBlue   = lipgloss.Color("#3B82F6") // Info/highlight
	ColorAccentPurple = lipgloss.Color("#8B5CF6") // Special

	// Vintage Colors (for intro scene and retro elements)
	ColorVintageGold    = lipgloss.Color("#D4AF37") // Vintage gold
	ColorVintageGoldDim = lipgloss.Color("#B8960F") // Dimmed vintage gold

	// Card Suit Colors
	ColorSuitRed   = lipgloss.Color("#DC2626") // Hearts, Diamonds
	ColorSuitBlack = lipgloss.Color("#1F2937") // Spades, Clubs

	// Status Colors
	ColorSuccess = lipgloss.Color("#10B981")
	ColorWarning = lipgloss.Color("#F59E0B")
	ColorError   = lipgloss.Color("#EF4444")
	ColorInfo    = lipgloss.Color("#3B82F6")

	// Border Colors
	ColorBorderSubtle = lipgloss.Color("#374151")
	ColorBorderNormal = lipgloss.Color("#4B5563")
	ColorBorderStrong = lipgloss.Color("#FFB900")
)
