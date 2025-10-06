package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/bunnyholes/pokerhole/client/internal/ui/constants"
)

const (
	frameWidth  = constants.TerminalWidth
	frameHeight = constants.TerminalHeight
)

func (m Model) contentWidth() int {
	frameExtraWidth, _ := shellStyle.GetFrameSize()
	width := frameWidth - frameExtraWidth
	if width < 0 {
		width = 0
	}
	return width
}

func (m Model) contentHeight() int {
	_, frameExtraHeight := shellStyle.GetFrameSize()
	height := frameHeight - frameExtraHeight
	if height < 0 {
		height = 0
	}
	return height
}

func (m Model) renderStatusBar() string {
	width := m.contentWidth()

	if m.status.message == "" {
		base := statusBarStyle(statusNeutral)
		hintIcon := spinnerStyle().Render("●")
		hint := fmt.Sprintf("%s 도움말 [?]  |  정보 [H]  |  종료 [Ctrl+C]", hintIcon)
		return base.Width(width).Render(" " + hint)
	}

	style := statusBarStyle(m.status.level)
	return style.Width(width).Render(" " + m.status.message)
}

func (m Model) applyShell(body string) string {
	width := m.contentWidth()
	interiorHeight := m.contentHeight()

	status := m.renderStatusBar()
	statusHeight := lipgloss.Height(status)
	if statusHeight < 1 {
		statusHeight = 1
	}

	bodyHeight := interiorHeight - statusHeight
	if bodyHeight < 0 {
		bodyHeight = 0
	}

	rawLines := strings.Split(body, "\n")
	formatted := make([]string, 0, bodyHeight)
	base := lipgloss.NewStyle().Width(width)
	for _, line := range rawLines {
		formatted = append(formatted, base.Render(line))
	}

	if len(formatted) > bodyHeight {
		formatted = formatted[:bodyHeight]
	}

	for len(formatted) < bodyHeight {
		formatted = append(formatted, base.Render(""))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, formatted...)
	layout := lipgloss.JoinVertical(lipgloss.Left, content, status)

	frame := shellStyle.
		Width(frameWidth).
		Height(frameHeight).
		Render(layout)

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			frame,
			lipgloss.WithWhitespaceForeground(ColorTextSecondary),
			lipgloss.WithWhitespaceBackground(ColorBgPrimary),
		)
	}

	return frame
}
