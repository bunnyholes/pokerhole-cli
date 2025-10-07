package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bunnyholes/pokerhole/client/internal/ui"
)

type menuItem struct {
	title       string
	description string
	disabled    bool
	disabledMsg string
}

type model struct {
	items    []menuItem
	selected int
	width    int
	height   int
	online   bool           // 서버 연결 상태
	spinner  spinner.Model // 접속 시도중 스피너
}

func initialModel() model {
	// 스피너 설정
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ui.ColorWarning)

	m := model{
		online:   false, // 서버 연결 없음 (데모용)
		selected: 0,
		spinner:  s,
	}

	// 메뉴 아이템 구성
	items := []menuItem{
		{
			title:       "랜덤 매치",
			description: "랜덤한 플레이어와 매칭됩니다. (오프라인 모드 포함)",
		},
		{
			title:       "코드 매칭",
			description: "같은 코드를 입력한 사람들끼리 매칭이 됩니다. (온라인 전용)",
			disabled:    true, // 오프라인에서는 비활성화
			disabledMsg: "온라인 연결이 필요합니다.",
		},
		{
			title:       "게임 종료",
			description: "포커홀 클라이언트를 종료합니다.",
		},
	}

	// 온라인 상태에 따라 코드 매칭 활성화
	if m.online {
		items[1].disabled = false
		items[1].disabledMsg = ""
	}

	m.items = items
	return m
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyUp:
			// 이전 활성화된 항목으로 이동
			prev := m.selected - 1
			if prev < 0 {
				prev = len(m.items) - 1
			}
			// 비활성화된 항목 건너뛰기
			for m.items[prev].disabled {
				prev--
				if prev < 0 {
					prev = len(m.items) - 1
				}
				// 무한 루프 방지
				if prev == m.selected {
					break
				}
			}
			m.selected = prev

		case tea.KeyDown, tea.KeyTab:
			// 다음 활성화된 항목으로 이동
			next := m.selected + 1
			if next >= len(m.items) {
				next = 0
			}
			// 비활성화된 항목 건너뛰기
			for m.items[next].disabled {
				next++
				if next >= len(m.items) {
					next = 0
				}
				// 무한 루프 방지
				if next == m.selected {
					break
				}
			}
			m.selected = next

		case tea.KeyEnter:
			item := m.items[m.selected]
			if item.disabled {
				// 비활성화된 항목은 무시
				return m, nil
			}
			// 실제로는 각 항목에 맞는 동작 수행
			// 데모이므로 종료만 구현
			if m.selected == 2 { // 게임 종료 (3번째 메뉴)
				return m, tea.Quit
			}
		}

		if msg.Type == tea.KeyRunes {
			switch strings.ToLower(string(msg.Runes)) {
			case "1":
				if len(m.items) > 0 && !m.items[0].disabled {
					m.selected = 0
				}
			case "2":
				if len(m.items) > 1 && !m.items[1].disabled {
					m.selected = 1
				}
			case "3":
				if len(m.items) > 2 && !m.items[2].disabled {
					m.selected = 2
				}
			case "q":
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	width := 80

	// 타이틀 (상하 구분선)
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true).
		Align(lipgloss.Center).
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderForeground(ui.ColorAccentGold)
	title := titleStyle.Render("POKERHOLE · TEXAS HOLD'EM TRAINER")

	// 상태 라인
	statusStyle := lipgloss.NewStyle().
		Foreground(ui.ColorBgPrimary).
		Background(ui.ColorWarning).
		Bold(true).
		Padding(0, 2)

	statusNoteStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextSecondary).
		PaddingLeft(2)

	// 접속 시도중 스피너
	statusLeft := lipgloss.JoinHorizontal(
		lipgloss.Left,
		statusStyle.Render("OFFLINE"),
		lipgloss.NewStyle().PaddingLeft(1).Render(m.spinner.View()),
		statusNoteStyle.Render("서버 연결 불가"),
	)

	playerStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true)

	playerLabelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextSecondary)

	statusRight := lipgloss.JoinHorizontal(
		lipgloss.Left,
		playerLabelStyle.Render("플레이어"),
		lipgloss.NewStyle().Render(" "),
		playerStyle.Render("GUEST"),
	)

	statusLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(width/2).Render(statusLeft),
		lipgloss.PlaceHorizontal(width-width/2, lipgloss.Right, statusRight),
	)

	// 메뉴만 표시
	menuPanel := m.renderMenuColumn(width)

	// 전체 레이아웃
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		statusLine,
		"",
		menuPanel,
	)

	// 외곽 테두리 없이 출력
	return content
}

func (m model) renderMenuColumn(width int) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true).
		Align(lipgloss.Center).
		Width(width).
		MarginBottom(2)

	entryStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(width)


	disabledStyle := entryStyle.Copy().
		Foreground(ui.ColorTextMuted)

	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextSecondary)

	var rows []string

	indicatorStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true)

	for i, item := range m.items {
		// 기본 인디케이터
		indicator := "  " // 공백 2칸

		// 선택된 항목은 인디케이터 변경
		if i == m.selected && !item.disabled {
			indicator = "► "
		}

		// 박스 스타일 적용
		if i == m.selected && !item.disabled {
			// 선택된 항목 - 배경색이 있는 풀사이즈 박스
			// 스타일을 중첩하지 않고 직접 텍스트 구성
			content := indicator + item.title + "\n" + "  " + item.description

			boxStyle := lipgloss.NewStyle().
				Width(width).
				Padding(1, 2).
				Background(lipgloss.Color("238")).
				Foreground(ui.ColorAccentGold).
				Bold(true).
				Align(lipgloss.Left)

			rows = append(rows, boxStyle.Render(content))
		} else {
			// 일반 항목 - 배경색 없음
			// 타이틀 스타일 설정
			ts := titleStyle
			if item.disabled {
				ts = ts.Copy().Foreground(ui.ColorTextMuted)
			}

			// 텍스트 구성
			titleLine := indicatorStyle.Render(indicator) + ts.Render(item.title)
			descLine := "  " + item.description

			// 비활성화된 항목 설명도 흐릿하게
			if item.disabled {
				descLine = descStyle.Copy().Foreground(ui.ColorTextMuted).Render(item.description)
				descLine = "  " + descLine
			}

			content := titleLine + "\n" + descLine

			style := entryStyle
			if item.disabled {
				style = disabledStyle
			}
			rows = append(rows, style.Render(content))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("게임 모드"),
		content,
	)
}

func (m model) renderDetailColumn(width int) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(ui.ColorAccentGold).
		Width(width).
		MarginBottom(1)

	contentStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(width)

	headingStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextPrimary).
		Bold(true).
		MarginBottom(1)

	bodyStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextSecondary)

	hintStyle := lipgloss.NewStyle().
		Foreground(ui.ColorTextSecondary).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccentGold).
		Bold(true)

	selected := m.items[m.selected]

	heading := headingStyle.Render(selected.title)
	description := bodyStyle.Render(selected.description)

	sections := []string{heading, description}

	if selected.disabled && selected.disabledMsg != "" {
		warning := bodyStyle.Copy().
			Foreground(ui.ColorWarning).
			Render(selected.disabledMsg)
		sections = append(sections, warning)
	}

	helpText := strings.Join([]string{
		keyStyle.Render("[Enter]") + bodyStyle.Render(" 실행"),
		keyStyle.Render("[↑/↓]") + bodyStyle.Render(" 이동"),
		keyStyle.Render("[q]") + bodyStyle.Render(" 종료"),
	}, "  ")
	sections = append(sections, hintStyle.Render(helpText))

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("선택된 항목"),
		contentStyle.Render(content),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
