package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/spinner"
	"github.com/charmbracelet/bubbles/v2/table"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type page int

const (
	pageSpinner page = iota
	pageProgress
	pageTextInput
	pageList
	pageTable
	pageViewport
	pageCards
	pageHelp
	totalPages
)

// List item implementation
type cardItem struct {
	title, desc string
}

func (i cardItem) Title() string       { return i.title }
func (i cardItem) Description() string { return i.desc }
func (i cardItem) FilterValue() string { return i.title }

// Key bindings
type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Quit},
	}
}

var keys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous page"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next page"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	currentPage page
	spinner     spinner.Model
	progress    progress.Model
	textInput   textinput.Model
	list        list.Model
	table       table.Model
	viewport    viewport.Model
	help        help.Model
	progressVal float64
}

func newModel() model {
	// Spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	// Progress
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
	)

	// TextInput
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = 50
	ti.SetWidth(50)

	// List
	items := []list.Item{
		cardItem{title: "A♠ Ace of Spades", desc: "The death card - highest card"},
		cardItem{title: "K♥ King of Hearts", desc: "The suicide king"},
		cardItem{title: "Q♦ Queen of Diamonds", desc: "Lady luck"},
		cardItem{title: "J♣ Jack of Clubs", desc: "One-eyed jack"},
		cardItem{title: "10♠ Ten of Spades", desc: "High value card"},
		cardItem{title: "9♥ Nine of Hearts", desc: "Nine of hearts"},
		cardItem{title: "8♦ Eight of Diamonds", desc: "Eight of diamonds"},
		cardItem{title: "7♣ Seven of Clubs", desc: "Seven of clubs"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 60, 20)
	l.Title = "Poker Cards"

	// Table
	columns := []table.Column{
		{Title: "Rank", Width: 6},
		{Title: "Hand", Width: 20},
		{Title: "Example", Width: 25},
	}
	rows := []table.Row{
		{"1", "Royal Flush", "A♠ K♠ Q♠ J♠ 10♠"},
		{"2", "Straight Flush", "9♥ 8♥ 7♥ 6♥ 5♥"},
		{"3", "Four of a Kind", "K♣ K♦ K♥ K♠ A♦"},
		{"4", "Full House", "Q♠ Q♦ Q♣ 7♥ 7♦"},
		{"5", "Flush", "A♦ J♦ 9♦ 6♦ 3♦"},
		{"6", "Straight", "10♠ 9♥ 8♦ 7♣ 6♠"},
		{"7", "Three of a Kind", "J♠ J♥ J♦ 5♣ 2♠"},
		{"8", "Two Pair", "A♠ A♥ K♦ K♣ 7♠"},
		{"9", "One Pair", "Q♠ Q♥ 9♦ 6♣ 3♠"},
		{"10", "High Card", "A♠ J♥ 8♦ 5♣ 2♠"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	// Viewport with content
	vp := viewport.New(
		viewport.WithWidth(60),
		viewport.WithHeight(15),
	)
	vp.SetContent(`Texas Hold'em Poker Rules

Each player is dealt two private cards (known as "hole cards") that belong to them alone.
Five community cards are dealt face-up on the "board".

Betting Rounds:
1. PRE-FLOP: After receiving 2 hole cards
2. FLOP: After 3 community cards dealt
3. TURN: After 4th community card
4. RIVER: After 5th community card
5. SHOWDOWN: Reveal hands, determine winner

Player Actions:
- FOLD: Quit the hand
- CHECK: Pass (only if currentBet == 0)
- CALL: Match current bet
- RAISE: Bet higher than current
- ALL_IN: Bet all chips

Hand Rankings (High to Low):
1. Royal Flush: A-K-Q-J-10 (same suit)
2. Straight Flush: 5 consecutive cards (same suit)
3. Four of a Kind: 4 cards same rank
4. Full House: 3 of a kind + pair
5. Flush: 5 cards same suit
6. Straight: 5 consecutive cards
7. Three of a Kind: 3 cards same rank
8. Two Pair: 2 pairs
9. One Pair: 1 pair
10. High Card: Highest card wins`)

	// Help
	h := help.New()

	return model{
		currentPage: pageSpinner,
		spinner:     sp,
		progress:    prog,
		textInput:   ti,
		list:        l,
		table:       t,
		viewport:    vp,
		help:        h,
		progressVal: 0.0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			if m.currentPage > 0 {
				m.currentPage--
			}
		case key.Matches(msg, keys.Right):
			if m.currentPage < totalPages-1 {
				m.currentPage++
			}
		case msg.String() >= "1" && msg.String() <= "8":
			pageNum := int(msg.String()[0] - '1')
			if pageNum < int(totalPages) {
				m.currentPage = page(pageNum)
			}
		}
	}

	// Update spinner
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// Update progress
	m.progressVal += 0.01
	if m.progressVal > 1.0 {
		m.progressVal = 0.0
	}

	// Update components based on current page
	switch m.currentPage {
	case pageTextInput:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	case pageList:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case pageTable:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	case pageViewport:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Padding(1, 0)

	switch m.currentPage {
	case pageSpinner:
		title := titleStyle.Render("PAGE 1: SPINNER COMPONENT")
		spinnerView := fmt.Sprintf("%s Loading...", m.spinner.View())
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", spinnerView)

	case pageProgress:
		title := titleStyle.Render("PAGE 2: PROGRESS COMPONENT")
		progressView := fmt.Sprintf("%s %.0f%%", m.progress.ViewAs(m.progressVal), m.progressVal*100)
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", progressView)

	case pageTextInput:
		title := titleStyle.Render("PAGE 3: TEXTINPUT COMPONENT")
		inputView := m.textInput.View()
		echoText := m.textInput.Value()
		if echoText == "" {
			echoText = "(your text will appear here)"
		}
		echo := fmt.Sprintf("Echo: %s", echoText)
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", inputView, "", echo)

	case pageList:
		title := titleStyle.Render("PAGE 4: LIST COMPONENT")
		listView := m.list.View()
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", listView)

	case pageTable:
		title := titleStyle.Render("PAGE 5: TABLE COMPONENT")
		tableView := m.table.View()
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", tableView)

	case pageViewport:
		title := titleStyle.Render("PAGE 6: VIEWPORT COMPONENT (Scrollable)")
		viewportView := m.viewport.View()
		scrollInfo := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", viewportView, "", scrollInfo)

	case pageCards:
		title := titleStyle.Render("PAGE 7: POKER CARDS")

		// Use lipgloss built-in API only - no manual pixel rendering
		cardStyle := lipgloss.NewStyle().
			Width(12).
			Height(7).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

		aceSpade := cardStyle.
			Background(lipgloss.Color("#5a3c8c")).
			Render("A♠")

		kingHeart := cardStyle.
			Background(lipgloss.Color("#d32f2f")).
			Render("K♥")

		queenDiamond := cardStyle.
			Background(lipgloss.Color("#ffa726")).
			Render("Q♦")

		jackClub := cardStyle.
			Background(lipgloss.Color("#43a047")).
			Render("J♣")

		// Empty blue card (no text)
		emptyCard := cardStyle.
			Background(lipgloss.Color("#2196f3")).
			Render("")

		cardRow := lipgloss.JoinHorizontal(lipgloss.Top, aceSpade, "  ", kingHeart, "  ", queenDiamond, "  ", jackClub, "  ", emptyCard)

		content = lipgloss.JoinVertical(lipgloss.Left, title, "", cardRow)

	case pageHelp:
		title := titleStyle.Render("PAGE 8: HELP COMPONENT")
		helpView := m.help.View(keys)
		content = lipgloss.JoinVertical(lipgloss.Left, title, "", helpView)
	}

	// Navigation
	pages := []string{
		"1: Spinner",
		"2: Progress",
		"3: TextInput",
		"4: List",
		"5: Table",
		"6: Viewport",
		"7: Cards",
		"8: Help",
	}

	var navItems []string
	for i, p := range pages {
		if page(i) == m.currentPage {
			navItems = append(navItems, lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFD700")).
				Render(p))
		} else {
			navItems = append(navItems, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render(p))
		}
	}

	nav := strings.Join(navItems, " | ")
	helpText := "← → / h l: Navigate | 1-8: Jump to page | q: Quit"

	return lipgloss.JoinVertical(lipgloss.Left, content, "", nav, helpText)
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
	}
}
