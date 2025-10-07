package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bunnyholes/pokerhole/client/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width  int
	height int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	width := m.width
	height := m.height
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 28
	}

	// 카드 보더 스타일
	cardBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	// 숨김 카드 보더 스타일 (커뮤니티 카드용)
	hiddenCardBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("220")) // 골드 색상

	// 스페이드 에이스 - 각 줄 개별 스타일링
	// 슈트 줄 - 좌측 정렬
	aceSuitLine := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Render("♠")

	// 랭크 줄 - 중앙 정렬
	aceRankLine := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Render("A")

	// 줄 결합
	aceContent := lipgloss.JoinVertical(
		lipgloss.Left,
		aceSuitLine,
		aceRankLine,
	)

	aceSpade := cardBorderStyle.Render(aceContent)

	// 하트 킹 - 각 줄 개별 스타일링
	// 슈트 줄 - 좌측 정렬
	kingSuitLine := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Render("♥")

	// 랭크 줄 - 중앙 정렬
	kingRankLine := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Render("K")

	// 줄 결합
	kingContent := lipgloss.JoinVertical(
		lipgloss.Left,
		kingSuitLine,
		kingRankLine,
	)

	kingHeart := cardBorderStyle.Render(kingContent)

	// 플레이어 카드를 나란히 배치
	playerCards := lipgloss.JoinHorizontal(
		lipgloss.Top,
		aceSpade,
		" ", // 카드 사이 간격
		kingHeart,
	)

	// 커뮤니티 카드 생성 (5장, 숨김 처리)
	// 숨김 카드 - 플레이어 카드와 동일한 레이아웃
	// 첫 번째 줄 - 좌측 정렬 (슈트 위치)
	hiddenSuitLine := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("220")). // 골드 색상
		Bold(true).
		Render("?")

	// 두 번째 줄 - 중앙 정렬 (랭크 위치)
	hiddenRankLine := lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("220")). // 골드 색상
		Bold(true).
		Render("?")

	// 줄 결합
	hiddenContent := lipgloss.JoinVertical(
		lipgloss.Left,
		hiddenSuitLine,
		hiddenRankLine,
	)

	// 5장의 숨김 카드 생성
	var communityCards []string
	for i := 0; i < 5; i++ {
		card := hiddenCardBorderStyle.Render(hiddenContent)
		communityCards = append(communityCards, card)
		if i < 4 {
			communityCards = append(communityCards, " ") // 카드 사이 간격
		}
	}

	// 커뮤니티 카드를 가로로 배치
	communityCardsRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		communityCards...,
	)

	// 플레이어 정보
	playerInfoStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextSecondary).
		MarginBottom(1).
		Width(11).
		Align(lipgloss.Center)

	playerInfo := playerInfoStyle.Render("Your Hand")

	// 커뮤니티 카드 정보
	communityInfoStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true).
		MarginBottom(1).
		Width(29). // 5장의 카드 너비에 맞춤
		Align(lipgloss.Center)

	communityInfo := communityInfoStyle.Render("Community Cards")

	// 커뮤니티 카드 섹션
	communitySection := lipgloss.JoinVertical(
		lipgloss.Center,
		communityInfo,
		communityCardsRow,
	)

	// 플레이어 카드 섹션
	playerSection := lipgloss.JoinVertical(
		lipgloss.Center,
		playerInfo,
		playerCards,
	)

	communityBlock := lipgloss.PlaceHorizontal(width, lipgloss.Center, communitySection)
	playerBlock := lipgloss.PlaceHorizontal(width, lipgloss.Center, playerSection)

	communityHeight := lipgloss.Height(communityBlock)
	playerHeight := lipgloss.Height(playerBlock)
	spaceAbovePlayer := height - playerHeight
	if spaceAbovePlayer < 0 {
		spaceAbovePlayer = 0
	}
	if spaceAbovePlayer < communityHeight {
		spaceAbovePlayer = communityHeight
	}

	topPadding := 0
	if spaceAbovePlayer > communityHeight {
		topPadding = (spaceAbovePlayer - communityHeight) / 2
	}
	gapAfterCommunity := spaceAbovePlayer - communityHeight - topPadding
	if gapAfterCommunity < 0 {
		gapAfterCommunity = 0
	}

	var layoutBuilder strings.Builder
	if topPadding > 0 {
		layoutBuilder.WriteString(strings.Repeat("\n", topPadding))
	}
	layoutBuilder.WriteString(communityBlock)
	if gapAfterCommunity > 0 {
		layoutBuilder.WriteString(strings.Repeat("\n", gapAfterCommunity))
	} else if len(communityBlock) > 0 && len(playerBlock) > 0 {
		layoutBuilder.WriteString("\n")
	}
	layoutBuilder.WriteString(playerBlock)

	return layoutBuilder.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
