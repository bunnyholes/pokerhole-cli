package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) handleModalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.modal {
	case modalHelp, modalAbout:
		switch msg.String() {
		case "esc", "q":
			m.modal = modalNone
			return m, nil
		}
	case modalShowdown:
		switch msg.String() {
		case "n", "N":
			return m.restartAfterShowdown()
		case "esc", "q":
			m.modal = modalNone
			m.screen = screenHome
			m.game.offlineGame = nil
			m = m.withStatus(statusInfo, "메뉴로 돌아갑니다.", 3*time.Second)
			return m, m.statusCommand(3 * time.Second)
		}
	}
	return m, nil
}

func (m Model) restartAfterShowdown() (tea.Model, tea.Cmd) {
	if m.game.offlineGame == nil {
		m = m.withStatus(statusError, "재시작할 게임이 없습니다.", 4*time.Second)
		return m, m.statusCommand(4 * time.Second)
	}

	if err := m.game.offlineGame.Restart(); err != nil {
		m = m.withStatus(statusError, err.Error(), 4*time.Second)
		return m, m.statusCommand(4 * time.Second)
	}

	m.game.snapshot = m.currentSnapshot()
	m.modal = modalNone
	m.screen = screenGame
	m = m.withStatus(statusSuccess, "새 게임이 시작되었습니다.", 3*time.Second)
	return m, m.statusCommand(3 * time.Second)
}

func (m Model) renderModalOverlay(base string) string {
	overlay := m.renderModal()
	width := m.contentWidth()

	overlayPlaced := lipgloss.Place(width, 0, lipgloss.Center, lipgloss.Center, overlay)
	combined := lipgloss.JoinVertical(lipgloss.Left, base, "", overlayPlaced)
	return combined
}

func (m Model) renderModal() string {
	switch m.modal {
	case modalHelp:
		return m.renderHelpModal()
	case modalAbout:
		return m.renderAboutModal()
	case modalShowdown:
		return m.renderShowdownModal()
	default:
		return ""
	}
}

func (m Model) renderHelpModal() string {
	width := minInt(68, m.contentWidth()-4)
	header := headerTitleStyle.Width(width).Render("포커 조작 도움말")
	content := []string{
		"플레이 액션:",
		"  [F] 폴드  |  [C] 콜",
		"  [R] 레이즈 | [K] 체크",
		"  [A] 올인",
		"",
		"일반 조작:",
		"  [ESC] 메뉴로 돌아가기",
		"  [H] 정보  |  [?] 도움말",
	}

	body := strings.Join(content, "\n")
	return panelEmphasisStyle.Width(width).Render(header + "\n\n" + body + "\n\n" + menuDescStyle.Render("[ESC] 또는 [Q] 닫기"))
}

func (m Model) renderAboutModal() string {
	width := minInt(68, m.contentWidth()-4)
	header := headerTitleStyle.Width(width).Render("PokerHole")
	lines := []string{
		"버전: v1.0.0",
		"클라이언트: Go, Bubble Tea",
		"스타일: Lip Gloss",
		"게임 아키텍처: 이벤트 소싱",
		"",
		"Made with ♠ ♥ ♦ ♣",
	}

	body := strings.Join(lines, "\n")
	footer := menuDescStyle.Render("[ESC] 또는 [Q] 닫기")
	return panelEmphasisStyle.Width(width).Render(header + "\n\n" + body + "\n\n" + footer)
}

func (m Model) renderShowdownModal() string {
	if m.game.offlineGame == nil {
		return ""
	}

	width := minInt(74, m.contentWidth()-4)
	snapshot := m.game.snapshot

	var lines []string
	lines = append(lines, headerTitleStyle.Align(lipgloss.Center).Width(width).Render("SHOWDOWN"))

	if snapshot.WinnerIndex >= 0 && snapshot.WinnerIndex < len(snapshot.Players) {
		winner := snapshot.Players[snapshot.WinnerIndex]
		message := fmt.Sprintf("승자: %s", winner.Nickname)
		lines = append(lines, statusBarStyle(statusSuccess).Width(width).Render(" "+message))
		if winner.HandRank != "" {
			lines = append(lines, menuDescStyle.Width(width).Render("핸드: "+winner.HandRank))
		}
	} else {
		lines = append(lines, statusBarStyle(statusInfo).Width(width).Render(" 비겼습니다"))
	}

	lines = append(lines, "")
	lines = append(lines, m.renderShowdownPlayers(width))
	lines = append(lines, "")
	lines = append(lines, menuDescStyle.Width(width).Render("[N] 새 게임  •  [ESC] 메뉴로"))

	content := strings.Join(lines, "\n")
	return panelEmphasisStyle.Width(width).Render(content)
}

func (m Model) renderShowdownPlayers(width int) string {
	if m.game.offlineGame == nil {
		return ""
	}

	snapshot := m.game.snapshot
	var rows []string

	for _, p := range snapshot.Players {
		cards := parseHand(p.Hand)
		cardView := renderHandCompact(cards, false)
		rank := p.HandRank
		if rank == "" {
			rank = "핸드 정보 없음"
		}
		row := lipgloss.JoinHorizontal(lipgloss.Left,
			menuItemStyle.Render(p.Nickname), "  ",
			cardView, "  ",
			menuDescStyle.Render(rank),
		)
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
