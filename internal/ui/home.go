package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	homeSectionLabelStyle = lipgloss.NewStyle().
				Foreground(ColorTextSecondary).
				Bold(true).
				MarginBottom(1)

	homeMenuEntryStyle = lipgloss.NewStyle().
				Padding(0, 1)

	homeMenuEntryActiveStyle = homeMenuEntryStyle.Copy().
					Background(ColorBgSecondary).
					Foreground(ColorAccentGold).
					Bold(true)

	homeMenuEntryDisabledStyle = homeMenuEntryStyle.Copy().
					Foreground(ColorTextMuted)

	homeMenuTitleStyle = lipgloss.NewStyle().
				Foreground(ColorTextPrimary).
				Bold(true)

	homeMenuDescriptionStyle = lipgloss.NewStyle().
					Foreground(ColorTextSecondary).
					MarginLeft(2)

	homeDetailPanelStyle = panelStyle.Copy()

	homeDetailHeadingStyle = lipgloss.NewStyle().
				Foreground(ColorAccentGold).
				Bold(true).
				MarginBottom(1)

	homeDetailBodyStyle = lipgloss.NewStyle().
				Foreground(ColorTextSecondary)

	homeHintStyle = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			MarginTop(1)
)

func (m Model) buildHomeMenu() []menuItem {
	items := []menuItem{
		{
			title:       "오프라인 연습전",
			description: "AI와 함께 텍사스 홀덤의 흐름을 연습합니다.",
			action:      homeActionOffline,
		},
		{
			title:       "온라인 매치",
			description: "실제 서버에 접속하여 다른 플레이어와 겨룹니다.",
			action:      homeActionOnlineMatch,
		},
		{
			title:       "게임 종료",
			description: "포커홀 클라이언트를 종료합니다.",
			action:      homeActionQuit,
		},
	}

	if !m.online {
		items[1].disabled = true
		if m.client == nil {
			items[1].disabledMsg = "서버 연결을 찾을 수 없습니다."
		} else {
			items[1].disabledMsg = "서버가 온라인 모드를 아직 지원하지 않습니다."
		}
	}

	return items
}

func (m Model) handleHomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp, tea.KeyShiftTab:
		if m.home.selected > 0 {
			m.home.selected--
		} else {
			m.home.selected = len(m.home.items) - 1
		}
		return m, nil

	case tea.KeyDown, tea.KeyTab:
		if m.home.selected < len(m.home.items)-1 {
			m.home.selected++
		} else {
			m.home.selected = 0
		}
		return m, nil

	case tea.KeyEnter:
		if len(m.home.items) == 0 {
			return m, nil
		}
		item := m.home.items[m.home.selected]
		if item.disabled {
			m = m.withStatus(statusWarning, item.disabledMsg, 4*time.Second)
			return m, m.statusCommand(4 * time.Second)
		}
		updated, cmd := m.executeMenuAction(item)
		return updated, cmd
	}

	if msg.Type == tea.KeyRunes {
		switch strings.ToLower(string(msg.Runes)) {
		case "1":
			return m.activateMenuItem(0)
		case "2":
			return m.activateMenuItem(1)
		case "3":
			return m.activateMenuItem(2)
		}
	}

	return m, nil
}

func (m Model) activateMenuItem(idx int) (tea.Model, tea.Cmd) {
	if idx < 0 || idx >= len(m.home.items) {
		return m, nil
	}
	m.home.selected = idx
	item := m.home.items[idx]
	if item.disabled {
		m = m.withStatus(statusWarning, item.disabledMsg, 4*time.Second)
		return m, m.statusCommand(4 * time.Second)
	}
	return m.executeMenuAction(item)
}

func (m Model) executeMenuAction(item menuItem) (tea.Model, tea.Cmd) {
	switch item.action {
	case homeActionOffline:
		return m.startOfflineSession()
	case homeActionOnlineMatch:
		m = m.withStatus(statusInfo, "온라인 매치는 준비 중입니다.", 4*time.Second)
		return m, m.statusCommand(4 * time.Second)
	case homeActionQuit:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m Model) viewHome() string {
	width := m.contentWidth()

	title := headerTitleStyle.Copy().
		Width(width).
		Align(lipgloss.Center).
		Render("POKERHOLE · TEXAS HOLD'EM TRAINER")

	status := m.homeStatusLine(width)

	gap := 4
	if width < 68 {
		gap = 2
	}

	leftWidth := (width - gap) / 2
	if leftWidth < 28 {
		leftWidth = 28
	}

	rightWidth := width - gap - leftWidth
	if rightWidth < 28 {
		rightWidth = 28
		if leftWidth+gap+rightWidth > width {
			leftWidth = width - gap - rightWidth
		}
	}

	left := m.renderMenuColumn(leftWidth)
	right := m.renderMenuDetail(rightWidth)
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, right)

	layout := lipgloss.JoinVertical(lipgloss.Left,
		title,
		status,
		"",
		body,
	)

	return m.applyShell(layout)
}

func (m Model) homeStatusLine(width int) string {
	leftWidth := width / 2
	rightWidth := width - leftWidth

	left := lipgloss.NewStyle().Width(leftWidth).Render(m.renderConnectionStatus())
	right := lipgloss.PlaceHorizontal(rightWidth, lipgloss.Right, m.renderPlayerSummary())

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) renderConnectionStatus() string {
	badge := lipgloss.NewStyle().
		Foreground(ColorBgPrimary).
		Bold(true).
		Padding(0, 2)

	label := "OFFLINE"
	note := "오프라인 연습 모드"

	if m.online {
		badge = badge.Background(ColorAccentGreen)
		label = "ONLINE"
		note = "서버 연결 완료"
	} else {
		badge = badge.Background(ColorWarning)
	}

	noteStyle := lipgloss.NewStyle().
		Foreground(ColorTextSecondary).
		PaddingLeft(2)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		badge.Render(label),
		noteStyle.Render(note),
	)
}

func (m Model) renderPlayerSummary() string {
	name := strings.TrimSpace(m.playerName)
	if name == "" {
		name = "GUEST"
	} else {
		name = strings.ToUpper(name)
	}

	labelStyle := lipgloss.NewStyle().Foreground(ColorTextSecondary)
	nameStyle := lipgloss.NewStyle().Foreground(ColorAccentGold).Bold(true)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("플레이어"),
		lipgloss.NewStyle().Render(" "),
		nameStyle.Render(name),
	)
}

func (m Model) renderMenuColumn(width int) string {
	if width < 28 {
		width = 28
	}

	frameWidth, _ := panelStyle.GetFrameSize()
	innerWidth := width - frameWidth
	if innerWidth < 1 {
		innerWidth = width
	}

	var rows []string

	for i, item := range m.home.items {
		base := lipgloss.JoinVertical(
			lipgloss.Left,
			homeMenuTitleStyle.Copy().Width(innerWidth).Render(fmt.Sprintf("%d. %s", i+1, item.title)),
			homeMenuDescriptionStyle.Copy().Width(innerWidth).Render(item.description),
		)

		style := homeMenuEntryStyle.Copy()
		if item.disabled {
			style = homeMenuEntryDisabledStyle.Copy()
		}
		if i == m.home.selected && !item.disabled {
			style = homeMenuEntryActiveStyle.Copy()
		}

		rows = append(rows, style.Width(innerWidth).Render(base))
	}

	if len(rows) == 0 {
		rows = append(rows, homeMenuEntryStyle.Copy().Width(innerWidth).Render("등록된 메뉴가 없습니다."))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		homeSectionLabelStyle.Copy().Width(width).Render("게임 모드"),
		panelStyle.Copy().Width(width).Render(content),
	)
}

func (m Model) renderMenuDetail(width int) string {
	if width < 28 {
		width = 28
	}

	frameWidth, _ := homeDetailPanelStyle.GetFrameSize()
	innerWidth := width - frameWidth
	if innerWidth < 1 {
		innerWidth = width
	}

	panel := homeDetailPanelStyle.Copy().Width(width)

	if len(m.home.items) == 0 {
		empty := homeDetailBodyStyle.Copy().Width(innerWidth).Render("선택 가능한 메뉴가 없습니다.")
		return panel.Render(empty)
	}

	selected := m.home.items[m.home.selected]

	heading := homeDetailHeadingStyle.Copy().Width(innerWidth).Render(selected.title)
	description := homeDetailBodyStyle.Copy().Width(innerWidth).Render(selected.description)

	sections := []string{heading, description}

	if selected.disabled && selected.disabledMsg != "" {
		warning := homeDetailBodyStyle.Copy().
			Foreground(ColorWarning).
			Width(innerWidth).
			Render(selected.disabledMsg)
		sections = append(sections, warning)
	}

	helpText := strings.Join([]string{
		helpKeyStyle.Render("[Enter]") + homeDetailBodyStyle.Render(" 실행"),
		helpKeyStyle.Render("[↑/↓]") + homeDetailBodyStyle.Render(" 이동"),
		helpKeyStyle.Render("[?]") + homeDetailBodyStyle.Render(" 도움말"),
		helpKeyStyle.Render("[Esc]") + homeDetailBodyStyle.Render(" 뒤로"),
	}, "  ")
	sections = append(sections, homeHintStyle.Copy().Width(innerWidth).Render(helpText))

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return panel.Render(content)
}
