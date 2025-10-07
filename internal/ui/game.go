package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bunnyholes/pokerhole/client/internal/core/application/service"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/card"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/player"
)

func (m Model) startOfflineSession() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.playerName)
	if name == "" {
		name = "Player"
	}

	game := service.NewOfflineGame(name)
	if err := game.Start(); err != nil {
		m = m.withStatus(statusError, fmt.Sprintf("게임 시작 실패: %v", err), 5*time.Second)
		return m, m.statusCommand(5 * time.Second)
	}

	m.game.offlineGame = game
	m.game.snapshot = game.GetGameState()
	m.screen = screenGame
	m.modal = modalNone
	m = m.withStatus(statusInfo, "오프라인 게임을 시작합니다.", 3*time.Second)

	cmd := tea.Batch(m.statusCommand(3*time.Second), animationTickCmd())
	return m, cmd
}

func (m Model) handleGameKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.game.offlineGame == nil {
		return m, nil
	}

	switch msg.String() {
	case "esc":
		m.screen = screenHome
		m = m.withStatus(statusInfo, "메뉴로 돌아갑니다.", 3*time.Second)
		return m, m.statusCommand(3 * time.Second)
	case "f":
		return m.performPlayerAction(vo.Fold, 0)
	case "c":
		return m.performPlayerAction(vo.Call, 0)
	case "k":
		return m.performPlayerAction(vo.Check, 0)
	case "a":
		return m.performPlayerAction(vo.AllIn, 0)
	case "r":
		amount := m.suggestRaiseAmount()
		return m.performPlayerAction(vo.Raise, amount)
	}

	return m, nil
}

func (m Model) performPlayerAction(action vo.PlayerAction, amount int) (tea.Model, tea.Cmd) {
	players := m.game.offlineGame.GetPlayers()
	if len(players) == 0 {
		return m, nil
	}

	me := players[0]
	snapshot := m.game.snapshot

	if snapshot.CurrentPlayer != 0 {
		m = m.withStatus(statusWarning, "지금은 AI 차례입니다.", 2*time.Second)
		return m, m.statusCommand(2 * time.Second)
	}

	if me.Status() == player.AllIn {
		m = m.withStatus(statusInfo, "올인 상태 - 자동으로 넘어갑니다.", 2*time.Second)
		return m, tea.Tick(450*time.Millisecond, func(time.Time) tea.Msg { return aiTurnMsg{} })
	}

	if action == vo.Raise {
		if amount <= snapshot.CurrentBet {
			amount = snapshot.CurrentBet + 50
		}
		maxAmount := me.Bet() + me.Chips()
		if amount >= maxAmount {
			action = vo.AllIn
			amount = maxAmount
		}
	}

	if err := m.game.offlineGame.PlayerAction(0, action, amount); err != nil {
		m = m.withStatus(statusError, fmt.Sprintf("액션 실패: %v", err), 4*time.Second)
		return m, m.statusCommand(4 * time.Second)
	}

	m.game.snapshot = m.currentSnapshot()
	m = m.withStatus(statusInfo, fmt.Sprintf("플레이어: %s", action.String()), 3*time.Second)

	updated, cmd := m.afterPlayerActed()
	return updated, tea.Batch(cmd, updated.statusCommand(3*time.Second))
}

func (m Model) afterPlayerActed() (Model, tea.Cmd) {
	if m.modal == modalShowdown {
		return m, nil
	}

	var cmds []tea.Cmd

	updated, cmd := m.evaluateRoundProgress()
	m = updated
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.game.snapshot = m.currentSnapshot()

	if m.modal == modalShowdown {
		if len(cmds) == 0 {
			return m, nil
		}
		return m, tea.Batch(cmds...)
	}

	if m.game.snapshot.CurrentPlayer == 1 && m.game.snapshot.Round != "SHOWDOWN" {
		cmds = append(cmds, tea.Tick(550*time.Millisecond, func(time.Time) tea.Msg { return aiTurnMsg{} }))
	}

	if len(cmds) == 0 {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m Model) performAITurn() (Model, tea.Cmd) {
	if m.game.offlineGame == nil {
		return m, nil
	}

	snapshot := m.game.snapshot
	if snapshot.CurrentPlayer != 1 || snapshot.Round == "SHOWDOWN" {
		return m, nil
	}

	players := m.game.offlineGame.GetPlayers()
	if len(players) < 2 {
		return m, nil
	}

	ai := players[1]
	var action vo.PlayerAction
	amount := snapshot.CurrentBet
	callAmount := snapshot.CurrentBet - ai.Bet()

	switch {
	case ai.Status() == player.AllIn:
		action = vo.Check
	case callAmount <= 0:
		action = vo.Check
	case callAmount < ai.Chips():
		action = vo.Call
	default:
		action = vo.AllIn
	}

	if err := m.game.offlineGame.PlayerAction(1, action, amount); err != nil {
		m = m.withStatus(statusError, fmt.Sprintf("AI 액션 실패: %v", err), 4*time.Second)
		return m, m.statusCommand(4 * time.Second)
	}

	m.game.snapshot = m.currentSnapshot()
	m = m.withStatus(statusInfo, fmt.Sprintf("AI: %s", action.String()), 3*time.Second)

	updated, cmd := m.evaluateRoundProgress()
	m = updated
	if cmd != nil {
		return m, tea.Batch(cmd, m.statusCommand(3*time.Second))
	}

	m.game.snapshot = m.currentSnapshot()
	if m.modal == modalShowdown {
		return m, m.statusCommand(3 * time.Second)
	}

	return m, m.statusCommand(3 * time.Second)
}

func (m Model) evaluateRoundProgress() (Model, tea.Cmd) {
	if m.game.offlineGame == nil {
		return m, nil
	}

	snapshot := m.game.snapshot
	players := m.game.offlineGame.GetPlayers()

	// Count active players and check if all have matching bets
	activePlayers := 0
	allBetsMatch := true
	maxBet := snapshot.CurrentBet

	for _, p := range players {
		status := p.Status()
		// Count non-folded players
		if status != player.Folded {
			activePlayers++
			// Check if active/waiting players have matching bets (AllIn players are exempt)
			if status == player.Active || status == player.Waiting {
				if p.Bet() != maxBet {
					allBetsMatch = false
				}
			}
		}
	}

	// Progress round if:
	// 1. Only one player left (others folded), OR
	// 2. All non-folded players have matching bets AND currentPlayer is back to 0
	//    (meaning betting round is complete - both players acted and turn cycled back)
	// Note: evaluateRoundProgress is only called AFTER a PlayerAction,
	// so currentPlayer==0 means we've completed a full cycle
	shouldProgress := activePlayers <= 1 || (allBetsMatch && snapshot.CurrentPlayer == 0)
	if !shouldProgress {
		return m, nil
	}

	if err := m.game.offlineGame.ProgressRound(); err != nil {
		m.status = statusState{message: fmt.Sprintf("라운드 진행 실패: %v", err), level: statusError, seq: m.status.seq + 1}
		return m, m.statusCommand(4 * time.Second)
	}

	m.game.snapshot = m.currentSnapshot()
	roundName := strings.ToUpper(m.game.snapshot.Round)
	m.status = statusState{message: fmt.Sprintf("라운드 진행: %s", roundName), level: statusInfo, seq: m.status.seq + 1}

	// Check if we reached showdown
	if m.game.snapshot.Round == "SHOWDOWN" {
		m.modal = modalShowdown
		return m, m.statusCommand(3 * time.Second)
	}

	return m, nil
}

func (m Model) suggestRaiseAmount() int {
	snapshot := m.game.snapshot
	players := m.game.offlineGame.GetPlayers()
	if len(players) == 0 {
		return snapshot.CurrentBet
	}

	me := players[0]
	minRaise := snapshot.CurrentBet + 50
	maxAmount := me.Bet() + me.Chips()
	return int(math.Min(float64(minRaise), float64(maxAmount)))
}

func (m Model) viewOfflineGame() string {
	if m.game.offlineGame == nil {
		return m.applyShell("오프라인 게임이 시작되지 않았습니다. [Enter]로 시작하세요.")
	}

	snapshot := m.game.snapshot

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		headerTitleStyle.Render("PokerHole - Offline Practice"),
		lipgloss.NewStyle().Foreground(ColorTextSecondary).PaddingLeft(2).Render(fmt.Sprintf("라운드: %s", snapshot.Round)),
	)

	community := m.renderCommunityArea(snapshot)
	players := m.renderPlayers(snapshot)
	pot := m.renderPotArea(snapshot)
	actions := m.renderActionBar()

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		community,
		"",
		players,
		"",
		pot,
		"",
		actions,
	)

	content := m.applyShell(body)

	if m.modal != modalNone {
		return m.renderModalOverlay(content)
	}

	return content
}

func (m Model) renderCommunityArea(snapshot service.GameStateSnapshot) string {
	label := headerMetaStyle.Render("Community")
	cards := renderCommunityCardsCompact(snapshot.CommunityCards)
	return panelStyle.Width(m.contentWidth()).Render(lipgloss.JoinVertical(lipgloss.Left, label, cards))
}

func (m Model) renderPlayers(snapshot service.GameStateSnapshot) string {
	width := m.contentWidth()

	var rows []string
	for idx, p := range snapshot.Players {
		hide := idx == 1 && snapshot.Round != "SHOWDOWN"
		cards := parseHand(p.Hand)
		line := renderPlayerRow(cards, p, hide, idx == snapshot.CurrentPlayer)
		rows = append(rows, line)
	}

	return panelStyle.Width(width).Render(strings.Join(rows, "\n"))
}

func renderPlayerRow(cards []card.Card, p service.PlayerSnapshot, hide bool, active bool) string {
	nameStyle := menuItemStyle
	if active {
		nameStyle = menuItemSelectedStyle
	}

	name := nameStyle.Render(p.Nickname)
	chips := statusBarStyle(statusNeutral).Render(fmt.Sprintf("칩 %d", p.Chips))
	bet := statusBarStyle(statusNeutral).Render(fmt.Sprintf("베팅 %d", p.Bet))

	cardView := renderHandCompact(cards, hide)

	parts := []string{name, chips, bet, cardView}
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func (m Model) renderPotArea(snapshot service.GameStateSnapshot) string {
	pot := statusBarStyle(statusInfo).Render(fmt.Sprintf("Pot %d", snapshot.Pot))
	bet := statusBarStyle(statusInfo).Render(fmt.Sprintf("현재 베팅 %d", snapshot.CurrentBet))
	return panelStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, pot, "  ", bet))
}

func (m Model) renderActionBar() string {
	actions := []struct {
		label string
		key   string
	}{
		{"폴드", "[F]"},
		{"콜", "[C]"},
		{"레이즈", "[R]"},
		{"체크", "[K]"},
		{"올인", "[A]"},
		{"메뉴", "[ESC]"},
	}

	var parts []string
	for _, action := range actions {
		parts = append(parts, helpKeyStyle.Render(action.key)+" "+action.label)
	}

	return panelStyle.Render(strings.Join(parts, "  "))
}
