package model

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour"
	"github.com/rigofekete/gator/internal/config"
	tuistyles "github.com/rigofekete/gator/internal/tui/styles"
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
		{name: "addfeed", label: "Add a new feed", args: []string{"Feed name", "Feed URL"}},
		{name: "agg", label: "Aggregate feeds", args: []string{"Poll interval (e.g., 30s)"}},
		{name: "browse", label: "Browse posts", args: []string{"Number of posts"}},
		{name: "users", label: "List all users"},
		{name: "feeds", label: "List all feeds"},
		{name: "following", label: "View your followed feeds"},
	}
}

type resultMsg struct {
	text string
	err  bool
}

type resultColorMsg struct{}

type loadingColorMsg struct{}

func parsePosts(raw string) string {
	lines := strings.Split(raw, "\n")
	var md strings.Builder
	for _, line := range lines {
		line = strings.TrimLeft(line, " \t")
		if strings.HasPrefix(line, "Found ") {
			md.WriteString(fmt.Sprintf("# %s\n\n", line))
		} else if strings.HasPrefix(line, "--- ") && strings.HasSuffix(line, " ---") {
			title := strings.TrimPrefix(line, "--- ")
			title = strings.TrimSuffix(title, " ---")
			md.WriteString(fmt.Sprintf("## %s\n\n", title))
		} else if strings.TrimSpace(line) != "" {
			md.WriteString(line + "\n")
		} else {
			md.WriteString("\n")
		}
	}
	return md.String()
}

type tuiModel struct {
	cursor   int
	cmdsList []tuiCmd
	view     string

	inputBuf  string
	selected  *tuiCmd
	collected []string

	resultMsg      string
	resultErr      bool
	resultColorIdx int

	width  int
	height int

	gatorCfg *config.Config

	loadingColorIdx int
	cancel          context.CancelFunc
	spinner         spinner.Model

	viewport viewport.Model
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
					m.view = "loading"
					m.loadingColorIdx = 0
					s := spinner.New()
					s.Spinner = spinner.Dot
					s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3c71a8"))
					m.spinner = s
					ctx, cancel := context.WithCancel(context.Background())
					m.cancel = cancel
					return m, tea.Batch(m.executeCmd(ctx, cmd, nil), m.spinner.Tick, tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg { return loadingColorMsg{} }))
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
					m.view = "loading"
					m.loadingColorIdx = 0
					s := spinner.New()
					s.Spinner = spinner.Dot
					s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3c71a8"))
					m.spinner = s
					ctx, cancel := context.WithCancel(context.Background())
					m.cancel = cancel
					return m, tea.Batch(m.executeCmd(ctx, m.selected, m.collected), m.spinner.Tick, tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg { return loadingColorMsg{} }))
				}
			case "backspace":
				if len(m.inputBuf) > 0 {
					m.inputBuf = m.inputBuf[:len(m.inputBuf)-1]
				}
			case "space":
				m.inputBuf += " "
			default:
				if len(msg.String()) == 1 {
					m.inputBuf += msg.String()
				}
			}
		case "result":
			m.view = "menu"
			m.selected = nil
		case "posts":
			switch msg.String() {
			case "q", "esc":
				m.view = "menu"
				m.selected = nil
				return m, nil
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case "loading":
			switch msg.String() {
			case "q", "ctrl+c":
				if m.cancel != nil {
					m.cancel()
				}
				return m, tea.Quit
			default:
				if m.cancel != nil {
					m.cancel()
				}
				m.view = "menu"
				m.selected = nil
				m.cancel = nil
				return m, nil
			}
		}
	case tea.PasteMsg:
		if m.view == "input" {
			m.inputBuf += msg.String()
		}
	case spinner.TickMsg:
		if m.view == "loading" {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case loadingColorMsg:
		if m.view == "loading" {
			m.loadingColorIdx = (m.loadingColorIdx + 1) % len(tuistyles.BlueGradient)
			return m, tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg { return loadingColorMsg{} })
		}
	case resultMsg:
		if m.view != "loading" {
			return m, nil
		}
		m.resultColorIdx = 0
		m.resultMsg = msg.text
		m.resultErr = msg.err
		m.cancel = nil

		if m.selected != nil && m.selected.name == "browse" && !msg.err {
			r, _ := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(m.width-2),
			)
			rendered, _ := r.Render(parsePosts(msg.text))
			m.viewport = viewport.New(
				viewport.WithWidth(m.width),
				viewport.WithHeight(m.height),
			)
			m.viewport.SetContent(rendered)
			m.view = "posts"
			return m, nil
		}

		m.view = "result"
		if cfg, err := config.Read(); err == nil {
			m.gatorCfg = &cfg
		}
		return m, tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg { return resultColorMsg{} })
	case resultColorMsg:
		if m.resultColorIdx < len(tuistyles.BlueGradient)-1 {
			m.resultColorIdx++
			return m, tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg { return resultColorMsg{} })
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.view == "posts" {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height)
		}
	}
	return m, nil
}

func (m tuiModel) executeCmd(ctx context.Context, cmd *tuiCmd, args []string) tea.Cmd {
	return func() tea.Msg {
		self, _ := os.Executable()
		execArgs := append([]string{"--exec", cmd.name}, args...)
		out, err := exec.CommandContext(ctx, self, execArgs...).CombinedOutput()
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
	case "loading":
		return m.loadingView()
	case "posts":
		return m.postsView()
	default:
		return m.menuView()
	}
}

func (m tuiModel) menuView() tea.View {
	var b strings.Builder
	b.WriteString(tuistyles.TitleStyle.Render(tuistyles.TitleASCII))
	b.WriteString("\n\n")
	var items []string
	for i, cmd := range m.cmdsList {
		if m.cursor == i {
			items = append(items, tuistyles.SelectedStyle.Render(" "+cmd.label+" "))
		} else {
			items = append(items, " "+cmd.label+" ")
		}
	}
	b.WriteString(lipgloss.NewStyle().PaddingLeft(1).Render(lipgloss.JoinVertical(lipgloss.Left, items...)))
	b.WriteString("\n\n")
	b.WriteString(tuistyles.UsernameStyle.Render(m.gatorCfg.CurrentUserName))
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String()))
}

func (m tuiModel) inputView() tea.View {
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(tuistyles.CursorStyle.Render(m.selected.args[len(m.collected)]) + "\n")

	currentArg := m.selected.args[len(m.collected)]
	borderStyle := tuistyles.InputBorderStyle
	if strings.Contains(currentArg, "URL") {
		borderStyle = tuistyles.WideInputBorderStyle
	}

	var inputText string
	if m.inputBuf == "" {
		inputText = borderStyle.Render(tuistyles.PlaceholderStyle.Render(currentArg))
	} else {
		inputText = borderStyle.Render(tuistyles.InputStyle.Render(m.inputBuf))
	}
	b.WriteString(fmt.Sprintf("%s\n\n", inputText))
	return tea.NewView(lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, b.String()))
}

func (m tuiModel) resultView() tea.View {
	var b strings.Builder
	blue := lipgloss.Color(tuistyles.BlueGradient[m.resultColorIdx])
	msgStyle := lipgloss.NewStyle().Foreground(blue).Bold(true)

	b.WriteString("\n\n")
	if m.resultErr {
		b.WriteString(tuistyles.ErrorStyle.Render(m.resultMsg))
	} else {
		b.WriteString(msgStyle.Render(m.resultMsg))
	}
	b.WriteString("\n\n")
	b.WriteString(tuistyles.HelpStyle.Render("press any key to continue"))
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String()))
}

func (m tuiModel) loadingView() tea.View {
	var b strings.Builder
	b.WriteString("\n\n")

	msg := "Processing..."
	if m.selected != nil && m.selected.name == "agg" {
		msg = "Aggregating feeds..."
	}

	blue := lipgloss.Color(tuistyles.BlueGradient[m.loadingColorIdx])
	msgStyle := lipgloss.NewStyle().Foreground(blue).Bold(true)
	b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), msgStyle.Render(msg)))
	b.WriteString(tuistyles.HelpStyle.Render("press any key to stop"))
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String()))
}

func (m tuiModel) postsView() tea.View {
	m.viewport.SetWidth(m.width)
	m.viewport.SetHeight(m.height)
	return tea.NewView(m.viewport.View())
}

func RunTUI(cfg *config.Config) {
	// Erase screen and move cursor top left
	fmt.Print("\033[2J\033[H")
	m := tuiModel{
		cmdsList: allCommands(),
		view:     "menu",
		gatorCfg: cfg,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print("\033[2J\033[H")
}
