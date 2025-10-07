package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bunnyholes/pokerhole/client/internal/core/application/service"
	"github.com/bunnyholes/pokerhole/client/internal/core/domain/game/vo"
)

func TestNewModelStartsInIntro(t *testing.T) {
	m := NewModel(nil, false, "Tester")
	if m.screen != screenIntro {
		t.Fatalf("expected screenIntro, got %v", m.screen)
	}
}

func TestSkipIntroMovesToHome(t *testing.T) {
	m := NewModel(nil, false, "Tester")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(Model)
	if m.screen != screenHome {
		t.Fatalf("expected screenHome after keypress, got %v", m.screen)
	}
}

func TestStartOfflineGameFromHome(t *testing.T) {
	m := NewModel(nil, false, "Tester")
	m.screen = screenHome

	updated, cmd := m.handleHomeKey(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("expected command after starting offline session")
	}
	m = updated.(Model)

	if m.screen != screenGame {
		t.Fatalf("expected screenGame, got %v", m.screen)
	}
	if m.game.offlineGame == nil {
		t.Fatalf("offline game should be initialized")
	}
}

func TestHelpModalLifecycle(t *testing.T) {
	m := NewModel(nil, false, "Tester")
	m.screen = screenHome

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updated.(Model)
	if m.modal != modalHelp {
		t.Fatalf("expected help modal, got %v", m.modal)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if m.modal != modalNone {
		t.Fatalf("expected modalNone, got %v", m.modal)
	}
}

func TestShowdownRestart(t *testing.T) {
	m := NewModel(nil, false, "Tester")
	session, _ := m.startOfflineSession()
	m = session.(Model)

	// Drive game to showdown manually
	m.game.offlineGame.PlayerAction(0, vo.Call, 0)
	m.game.offlineGame.PlayerAction(1, vo.Check, 0)
	m.game.offlineGame.ProgressRound()
	m.game.offlineGame.PlayerAction(0, vo.Check, 0)
	m.game.offlineGame.PlayerAction(1, vo.Check, 0)
	m.game.offlineGame.ProgressRound()
	m.game.offlineGame.PlayerAction(0, vo.Check, 0)
	m.game.offlineGame.PlayerAction(1, vo.Check, 0)
	m.game.offlineGame.ProgressRound()
	m.game.offlineGame.PlayerAction(0, vo.Check, 0)
	m.game.offlineGame.PlayerAction(1, vo.Check, 0)
	m.game.offlineGame.ProgressRound()

	m.game.snapshot = m.currentSnapshot()
	m.modal = modalShowdown

	updated, cmd := m.handleModalKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if cmd == nil {
		t.Fatalf("expected command when restarting")
	}
	m = updated.(Model)

	if m.modal != modalNone {
		t.Fatalf("expected modalNone after restart, got %v", m.modal)
	}
	if m.game.snapshot.Round != "PRE_FLOP" {
		t.Fatalf("expected new round PRE_FLOP, got %s", m.game.snapshot.Round)
	}
}

func TestStatusCommandClearsMessage(t *testing.T) {
	m := NewModel(nil, false, "Tester")
	m = m.withStatus(statusInfo, "테스트", 10*time.Millisecond)
	if cmd := m.statusCommand(10 * time.Millisecond); cmd == nil {
		t.Fatalf("expected non-nil command")
	}
}

// ensure offline package referenced for build
var _ = service.NewOfflineGame
