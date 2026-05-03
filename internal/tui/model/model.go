package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/rigofekete/gator/internal/tui/styles"
)

type tuiCmd struct {
	name  string
	label string
	args  []string
}

func allCommands() []tuiCmd {
	return []tuiCmd{
		{name: "register", label: "Register a new user", args: []string{"Username"}},
		{name: "login", label: "Login user", args: []string{"Username"}},
	}
}

type resultMsg struct {
	text string
	err  bool
}

type resultColorMsg struct{}

type tuiModel struct {
	cursor   int
	cmdsList []tuiCmd
	view     string

	inputBuf  string
	selected  *tuiCmd
	collected []string

	resultMsg string
	resultErr bool

	resultColorIdx int

	width  int
	height int
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch m.view {
		case "menu":
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.cmdsList)-1 {
					m.cursor++
				}
			case "enter":
				cmd := &m.cmdsList[m.cursor]
				m.selected = cmd
				m.collected = make([]string, 0)
				m.inputBuf = ""
				if len(cmd.args) == 0 {
					return m, m.executeCmd(cmd, nil)
				}
				m.view = "input"
			}
		case "input":
			switch msg.String() {
			case "esc":
				m.view = "menu"
				m.selected = nil
			case "enter":
				m.collected = append(m.collected, m.inputBuf)
				m.inputBuf = ""
				if len(m.collected) >= len(m.selected.args) {
					return m, m.executeCmd(m.selected, m.collected)
				}
			case "backspace":
				if len(m.inputBuf) > 0 {
					m.inputBuf = m.inputBuf[:len(m.inputBuf)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.inputBuf += msg.String()
				}
			}
		case "result":
			m.view = "menu"
			m.selected = nil
		}
	case resultMsg:
		m.resultMsg = msg.text
		m.resultErr = msg.err
		m.view = "result"
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return resultColorMsg{} })
	case resultColorMsg:
		m.resultColorIdx = (m.resultColorIdx + 1) % len(styles.BlueGradient)
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return resultColorMsg{} })
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m tuiModel) executeCmd(cmd *tuiCmd, args []string) tea.Cmd {
	return func() tea.Msg {
		self, _ := os.Executable()
		execArgs := append([]string{"--exec", cmd.name}, args...)
		out, err := exec.Command(self, execArgs...).CombinedOutput()
		return resultMsg{
			text: strings.TrimSpace(string(out)),
			err:  err != nil,
		}
	}
}

func (m tuiModel) View() tea.View {
	switch m.view {
	case "input":
		return m.inputView()
	case "result":
		return m.resultView()
	default:
		return m.menuView()
	}
}

func (m tuiModel) menuView() tea.View {
	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render(styles.TitleASCII))
	b.WriteString("\n\n")
	var items []string
	for i, cmd := range m.cmdsList {
		if m.cursor == i {
			items = append(items, styles.SelectedStyle.Render(" "+cmd.label))
		} else {
			items = append(items, " "+cmd.label)
		}
	}
	b.WriteString(lipgloss.NewStyle().PaddingLeft(1).Render(lipgloss.JoinVertical(lipgloss.Left, items...)))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("enter select - q quit"))
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String()))
}

func (m tuiModel) inputView() tea.View {
	var b strings.Builder
	b.WriteString("\n\n")
	// TODO: rename and restyle this cursorStyle
	b.WriteString(styles.CursorStyle.Render(m.selected.label) + "\n")
	var inputText string
	if m.inputBuf == "" {
		inputText = styles.AppStyle.Render(styles.PlaceholderStyle.Render("username"))
	} else {
		inputText = styles.AppStyle.Render(styles.InputStyle.Render(m.inputBuf))
	}
	b.WriteString(fmt.Sprintf("%s\n\n", inputText))
	return tea.NewView(lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, b.String()))
}

func (m tuiModel) resultView() tea.View {
	var b strings.Builder
	blue := lipgloss.Color(styles.BlueGradient[m.resultColorIdx])
	msgStyle := lipgloss.NewStyle().Foreground(blue).Bold(true)

	b.WriteString("\n\n")
	if m.resultErr {
		b.WriteString(styles.ErrorStyle.Render(m.resultMsg))
	} else {
		b.WriteString(msgStyle.Render(m.resultMsg))
	}
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("press any key to continue"))
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String()))
}

func RunTUI() {
	fmt.Print("\033[2J\033[H")
	m := tuiModel{
		cmdsList: allCommands(),
		view:     "menu",
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
	// Erase screen and move cursor top left
	fmt.Print("\033[2J\033[H")
}
