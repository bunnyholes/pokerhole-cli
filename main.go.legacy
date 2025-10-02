package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

var debugLog *log.Logger

// GameMode represents the current game mode
type GameMode int

const (
	ModeSplash     GameMode = iota // 초기 스플래시 화면
	ModeConnecting                 // 연결 시도 중 (백그라운드)
	ModeOnlineMenu                 // 온라인 모드 선택 화면
	ModeOnlineGame                 // 온라인 게임 화면
	ModeLocalGame                  // 로컬 게임 화면
)

// MatchType represents online match selection
type MatchType int

const (
	MatchRandom MatchType = iota // 랜덤 매칭
	MatchCode                    // 코드 매칭
)

// Connection result messages
type connectionSuccessMsg struct{}
type connectionFailureMsg struct{}
type splashTimeoutMsg struct{} // 5초 스플래시 타임아웃

// ClientInfo represents client identification
type ClientInfo struct {
	UUID     string `json:"uuid"`
	Nickname string `json:"nickname"`
}

// Model represents the application state
type model struct {
	spinner       spinner.Model
	mode          GameMode
	selectedMatch MatchType  // 온라인 모드 선택
	cursor        int        // 메뉴 커서 위치
	isOnline      bool       // 온라인 모드 가능 여부
	clientInfo    ClientInfo // 클라이언트 식별 정보
}

// Init is called when the program starts
func (m model) Init() tea.Cmd {
	// 스플래시 화면 5초 타이머 + 백그라운드 연결 시도
	return tea.Batch(
		m.spinner.Tick,
		tryConnect(m.clientInfo),
		splashTimer(),
	)
}

// nickhole generates a random nickname
// TODO: 추후 고도화 예정
func nickhole() string {
	adjectives := []string{
		"Lucky", "Brave", "Smart", "Quick", "Wild",
		"Cool", "Swift", "Bold", "Wise", "Epic",
		"Ace", "King", "Royal", "Grand", "Elite",
	}
	nouns := []string{
		"Shark", "Tiger", "Eagle", "Wolf", "Fox",
		"Dragon", "Phoenix", "Lion", "Bear", "Hawk",
		"Dealer", "Player", "Master", "Champion", "Legend",
	}

	// 랜덤 인덱스 생성
	adjIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	nounIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nouns))))
	numIdx, _ := rand.Int(rand.Reader, big.NewInt(1000))

	return fmt.Sprintf("%s%s%d",
		adjectives[adjIdx.Int64()],
		nouns[nounIdx.Int64()],
		numIdx.Int64())
}

// getOrCreateUUID loads existing UUID from file or creates new one
func getOrCreateUUID() string {
	// 홈 디렉토리의 .pokerhole 디렉토리에 UUID 저장
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 홈 디렉토리를 찾을 수 없으면 새 UUID 생성
		return uuid.New().String()
	}

	configDir := homeDir + "/.pokerhole"
	uuidFile := configDir + "/client_uuid.txt"

	// 기존 UUID 파일 읽기 시도
	if data, err := os.ReadFile(uuidFile); err == nil {
		existingUUID := string(data)
		if len(existingUUID) > 0 {
			fmt.Printf("기존 UUID 로드: %s\n", existingUUID)
			return existingUUID
		}
	}

	// 새 UUID 생성
	newUUID := uuid.New().String()
	fmt.Printf("새 UUID 생성: %s\n", newUUID)

	// 디렉토리 생성 (없으면)
	os.MkdirAll(configDir, 0755)

	// UUID 파일에 저장
	os.WriteFile(uuidFile, []byte(newUUID), 0644)

	return newUUID
}

// generateClientInfo creates unique client identification
func generateClientInfo() ClientInfo {
	return ClientInfo{
		UUID:     getOrCreateUUID(),
		Nickname: nickhole(),
	}
}

// splashTimer waits before transitioning (2 seconds for server connection)
func splashTimer() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return splashTimeoutMsg{}
	})
}

// tryConnect attempts to connect to the server with 1 second timeout
func tryConnect(clientInfo ClientInfo) tea.Cmd {
	return func() tea.Msg {
		debugLog.Printf("서버 연결 시도: localhost:7777")

		// 1초 타임아웃으로 서버 연결 시도 (빠른 데모용)
		conn, err := net.DialTimeout("tcp", "localhost:7777", 1*time.Second)
		if err != nil {
			debugLog.Printf("연결 실패: %v", err)
			return connectionFailureMsg{}
		}
		defer conn.Close()

		debugLog.Printf("연결 성공!")

		// 클라이언트 정보를 JSON으로 직렬화하여 전송
		data, err := json.Marshal(clientInfo)
		if err != nil {
			debugLog.Printf("JSON 직렬화 실패: %v", err)
			return connectionFailureMsg{}
		}

		debugLog.Printf("JSON 전송: %s", string(data))

		// 서버에 클라이언트 정보 전송
		_, err = conn.Write(append(data, '\n'))
		if err != nil {
			debugLog.Printf("전송 실패: %v", err)
			return connectionFailureMsg{}
		}

		debugLog.Printf("전송 완료 - connectionSuccessMsg 반환")
		return connectionSuccessMsg{}
	}
}

// Update handles incoming messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.mode == ModeOnlineMenu && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.mode == ModeOnlineMenu && m.cursor < 1 {
				m.cursor++
			}
		case "enter", " ":
			// 온라인 메뉴에서 선택
			if m.mode == ModeOnlineMenu {
				m.selectedMatch = MatchType(m.cursor)
				m.mode = ModeOnlineGame
			}
		}

	case splashTimeoutMsg:
		// 5초 타임아웃 -> 연결 상태에 따라 화면 전환
		debugLog.Printf("splashTimeoutMsg 수신 - isOnline: %v", m.isOnline)
		if m.isOnline {
			debugLog.Println("온라인 모드로 전환")
			m.mode = ModeOnlineMenu
		} else {
			debugLog.Println("로컬 모드로 전환")
			m.mode = ModeLocalGame
		}
		return m, nil

	case connectionSuccessMsg:
		// 서버 연결 성공 (백그라운드)
		debugLog.Println("connectionSuccessMsg 수신 - 온라인 모드 활성화")
		m.isOnline = true
		// 스플래시가 끝나면 자동 전환됨
		return m, nil

	case connectionFailureMsg:
		// 서버 연결 실패 (백그라운드)
		debugLog.Println("connectionFailureMsg 수신 - 로컬 모드로 설정")
		m.isOnline = false
		// 스플래시가 끝나면 자동 전환됨
		return m, nil

	default:
		// Update spinner animation
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	switch m.mode {
	case ModeSplash:
		return m.viewSplash()
	case ModeConnecting:
		return m.viewConnecting()
	case ModeOnlineMenu:
		return m.viewOnlineMenu()
	case ModeOnlineGame:
		return m.viewOnlineGame()
	case ModeLocalGame:
		return m.viewLocalGame()
	default:
		return "Unknown mode"
	}
}

// viewSplash shows the initial splash screen (5 seconds)
func (m model) viewSplash() string {
	var (
		primary   = lipgloss.Color("#FF6B9D")
		secondary = lipgloss.Color("#C9A0DC")
		accent    = lipgloss.Color("#67D5B5")
		gold      = lipgloss.Color("#FFD700")
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 3)

	titleStyle := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(accent).
		Italic(true).
		Align(lipgloss.Center)

	featureStyle := lipgloss.NewStyle().
		Foreground(secondary).
		PaddingLeft(2)

	title := titleStyle.Render("♠ ♥ ♦ ♣  POKER HOLDEM  ♣ ♦ ♥ ♠")
	subtitle := subtitleStyle.Render("Terminal Edition")

	// 원래 메인 화면 기능 소개
	features := featureStyle.Render(
		"* Real-time multiplayer\n" +
			"* Beautiful terminal UI\n" +
			"* Easy controls\n" +
			"* Sound effects",
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		subtitle,
		"",
		features,
		"",
	)

	return boxStyle.Render(content)
}

// viewConnecting shows the connection attempt screen
func (m model) viewConnecting() string {
	// Charm-inspired color palette
	var (
		primary = lipgloss.Color("#FF6B9D") // Pink
		accent  = lipgloss.Color("#67D5B5") // Teal
		gold    = lipgloss.Color("#FFD700") // Gold
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 3)

	titleStyle := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Align(lipgloss.Center)

	spinnerStyle := lipgloss.NewStyle().Foreground(accent)

	title := titleStyle.Render("♠ ♥ ♦ ♣  POKER HOLDEM  ♣ ♦ ♥ ♠")
	spinnerText := spinnerStyle.Render(m.spinner.View() + " Connecting to server...")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		spinnerText,
		"",
	)

	return boxStyle.Render(content)
}

// viewOnlineMenu shows the online mode menu
func (m model) viewOnlineMenu() string {
	var (
		primary   = lipgloss.Color("#FF6B9D")
		secondary = lipgloss.Color("#C9A0DC")
		accent    = lipgloss.Color("#67D5B5")
		gold      = lipgloss.Color("#FFD700")
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 3)

	titleStyle := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Align(lipgloss.Center)

	menuStyle := lipgloss.NewStyle().
		Foreground(secondary).
		PaddingLeft(2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		PaddingLeft(2)

	title := titleStyle.Render("♠ ♥ ♦ ♣  ONLINE MODE  ♣ ♦ ♥ ♠")

	// 메뉴 아이템 렌더링
	var menu string
	choices := []string{"Random Match", "Code Match"}
	for i, choice := range choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			menu += selectedStyle.Render(cursor+" "+choice) + "\n"
		} else {
			menu += menuStyle.Render(cursor+" "+choice) + "\n"
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		menu,
		"",
	)

	return boxStyle.Render(content)
}

// viewOnlineGame shows the online game screen (dummy)
func (m model) viewOnlineGame() string {
	var (
		primary = lipgloss.Color("#FF6B9D")
		gold    = lipgloss.Color("#FFD700")
		accent  = lipgloss.Color("#67D5B5")
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 3)

	titleStyle := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Align(lipgloss.Center)

	infoStyle := lipgloss.NewStyle().
		Foreground(accent).
		Align(lipgloss.Center)

	matchType := "Random Match"
	if m.selectedMatch == MatchCode {
		matchType = "Code Match"
	}

	title := titleStyle.Render("♠ ♥ ♦ ♣  ONLINE GAME  ♣ ♦ ♥ ♠")
	info := infoStyle.Render("Mode: " + matchType)
	dummy := infoStyle.Render("\n[Game Screen - Coming Soon]\n")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		info,
		dummy,
		"",
	)

	return boxStyle.Render(content)
}

// viewLocalGame shows the local game screen (dummy)
func (m model) viewLocalGame() string {
	var (
		primary = lipgloss.Color("#FF6B9D")
		gold    = lipgloss.Color("#FFD700")
		accent  = lipgloss.Color("#67D5B5")
		subtle  = lipgloss.Color("#6B6B83")
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 3)

	titleStyle := lipgloss.NewStyle().
		Foreground(gold).
		Bold(true).
		Align(lipgloss.Center)

	infoStyle := lipgloss.NewStyle().
		Foreground(accent).
		Align(lipgloss.Center)

	warningStyle := lipgloss.NewStyle().
		Foreground(subtle).
		Italic(true).
		Align(lipgloss.Center)

	title := titleStyle.Render("♠ ♥ ♦ ♣  LOCAL MODE  ♣ ♦ ♥ ♠")
	warning := warningStyle.Render("(Server connection failed)")
	dummy := infoStyle.Render("\n[Local Game Screen - Coming Soon]\n")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		warning,
		"",
		dummy,
		"",
	)

	return boxStyle.Render(content)
}

func main() {
	// 디버그 로그 파일 초기화
	homeDir, _ := os.UserHomeDir()
	logFile, err := os.OpenFile(homeDir+"/.pokerhole/client.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("로그 파일 생성 실패: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	debugLog = log.New(logFile, "[CLIENT] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	debugLog.Println("==================== 클라이언트 시작 ====================")

	// Initialize spinner with dots style
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#67D5B5"))

	// 클라이언트 고유 정보 생성
	clientInfo := generateClientInfo()
	debugLog.Printf("Client UUID: %s", clientInfo.UUID)
	debugLog.Printf("Nickname: %s", clientInfo.Nickname)

	p := tea.NewProgram(
		model{
			spinner:    s,
			mode:       ModeSplash, // 초기 스플래시 화면
			cursor:     0,
			isOnline:   false,
			clientInfo: clientInfo,
		},
		tea.WithAltScreen(), // Full screen mode
	)

	if _, err := p.Run(); err != nil {
		debugLog.Printf("에러 발생: %v", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	debugLog.Println("==================== 클라이언트 종료 ====================")
}
