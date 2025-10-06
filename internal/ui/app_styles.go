package ui

import "github.com/charmbracelet/lipgloss"

var (
	appBackgroundStyle = lipgloss.NewStyle().
				Background(ColorBgPrimary).
				Foreground(ColorTextPrimary)

	shellStyle = appBackgroundStyle.Copy().
			Padding(1, 2)

	headerTitleStyle = lipgloss.NewStyle().
				Foreground(ColorAccentGold).
				Background(ColorBgSecondary).
				Bold(true).
				Padding(0, 1).
				MarginBottom(1)

	headerMetaStyle = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Background(ColorBgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNormal).
			Padding(1, 2)

	panelEmphasisStyle = panelStyle.Copy().
				BorderForeground(ColorAccentGold)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Padding(0, 0, 0, 2)

	menuItemSelectedStyle = menuItemStyle.Copy().
				Background(ColorBgPrimary).
				Foreground(ColorAccentGold).
				Bold(true).
				BorderLeft(true).
				BorderForeground(ColorAccentGold)

	menuItemDisabledStyle = menuItemStyle.Copy().
				Foreground(ColorTextMuted)

	menuDescStyle = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Padding(0, 0, 1, 4)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorAccentBlue).
			Bold(true)
)

func statusBarStyle(level statusLevel) lipgloss.Style {
	base := lipgloss.NewStyle().
		Background(ColorBgSecondary).
		Foreground(ColorTextSecondary).
		Padding(0, 1)

	switch level {
	case statusInfo:
		return base.Copy().Foreground(ColorAccentBlue)
	case statusSuccess:
		return base.Copy().Foreground(ColorAccentGreen)
	case statusWarning:
		return base.Copy().Foreground(ColorWarning)
	case statusError:
		return base.Copy().Foreground(ColorError)
	default:
		return base
	}
}

func spinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorAccentBlue)
}
