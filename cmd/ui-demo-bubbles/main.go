package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/bunnyholes/pokerhole/client/internal/core/application/service"
)

func main() {
	// Create offline game
	game := service.NewOfflineGame("테스트플레이어")
	game.Start()

	gameState := game.GetGameState()

	// Initialize components
	columns := []table.Column{
		{Title: "", Width: 3},
		{Title: "플레이어", Width: 20},
		{Title: "칩", Width: 15},
		{Title: "베팅", Width: 8},
		{Title: "핸드", Width: 20},
		{Title: "상태", Width: 10},
	}
	playerTable := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(5),
	)

	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("86")).
		BorderBottom(true).
		Bold(true)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	playerTable.SetStyles(tableStyle)

	chipProgress := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	// Color scheme
	titleColor := lipgloss.Color("205")
	borderColor := lipgloss.Color("86")
	potColor := lipgloss.Color("220")
	cardColor := lipgloss.Color("213")

	// Title - full width
	termWidth := 120 // Default terminal width
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(titleColor).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(borderColor).
		Padding(0, 2).
		Align(lipgloss.Center).
		Width(termWidth)
	title := titleStyle.Render("포커홀")

	// Game info boxes
	roundStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("141")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("141")).
		Padding(0, 2)
	roundBox := roundStyle.Render(fmt.Sprintf("■ 라운드: %s", gameState.Round))

	// Pot with progress bar
	maxPot := 2000
	potPercent := float64(gameState.Pot) / float64(maxPot)
	if potPercent > 1.0 {
		potPercent = 1.0
	}
	potBar := chipProgress.ViewAs(potPercent)
	potStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(potColor).
		Padding(0, 1)
	potBox := potStyle.Render(fmt.Sprintf("※ 팟: %d 칩\n%s", gameState.Pot, potBar))

	betStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("208")).
		Padding(0, 2)
	betBox := betStyle.Render(fmt.Sprintf("● 현재 베팅: %d", gameState.CurrentBet))

	gameInfo := lipgloss.JoinHorizontal(lipgloss.Top, roundBox, "  ", potBox, "  ", betBox)

	// Community cards - full width
	cardBoxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(cardColor).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(termWidth)
	communityCards := cardBoxStyle.Render("커뮤니티 카드\n\n[ 대기중... ]")

	// Player table
	var rows []table.Row
	maxChips := 1000

	for i, p := range gameState.Players {
		indicator := "  "
		if i == gameState.CurrentPlayer {
			indicator = "▶"
		}

		playerName := p.Nickname
		if i == 0 {
			playerName = "[당신] " + playerName
		} else {
			playerName = "[컴퓨터] " + playerName
		}

		chipPercent := float64(p.Chips) / float64(maxChips)
		if chipPercent > 1.0 {
			chipPercent = 1.0
		}
		chipsWithBar := fmt.Sprintf("%d\n%s", p.Chips, chipProgress.ViewAs(chipPercent))

		handStr := "[ 숨김 ]"
		if i == 0 && p.Hand != "" {
			handStr = p.Hand
		}

		rows = append(rows, table.Row{
			indicator,
			playerName,
			chipsWithBar,
			fmt.Sprintf("%d", p.Bet),
			handStr,
			p.Status,
		})
	}

	playerTable.SetRows(rows)

	tableTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Padding(0, 0, 1, 0).
		Render("━━━ 플레이어 정보 ━━━")

	tableBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)
	tableContent := tableBox.Render(tableTitle + "\n" + playerTable.View())

	// Action guide - full width
	actionGuide := "■ 조작키: [F] 폴드  |  [C] 콜  |  [R] 레이즈  |  [A] 올인  |  [K] 체크"
	actionBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("141")).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(termWidth).
		Render(actionGuide)

	// Print game screen only
	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	fmt.Println(gameInfo)
	fmt.Println()
	fmt.Println(communityCards)
	fmt.Println()
	fmt.Println(tableContent)
	fmt.Println()
	fmt.Println(actionBox)
	fmt.Println()
}
