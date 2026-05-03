package styles

import (
	"charm.land/lipgloss/v2"
)

var BlueGradient = []string{
	"#1a3c6e", "#1f4d82", "#2863a0", "#3576b4",
	"#3c71a8", "#4a86c4", "#5b9bd5", "#6fb0e0",
	"#7ec8e3", "#91d4e8", "#a6d4e8", "#82b8d8",
	"#6fa3cc", "#5290bf", "#3a7ab0", "#276499",
}

var (
	TitleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5f5fff"))
	SelectedStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#ff5faf")).Foreground(lipgloss.Color("#eeeeee"))
	CursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5faf"))
	PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	HelpStyle        = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#626262"))
	ErrorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	InputStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#add8e6"))
	PromptStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#afffff"))
	AppStyle         = lipgloss.NewStyle().
				Padding(0, 1).
				Width(20).Height(2).
				Align(lipgloss.Center).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#3c71a8"))
)

const TitleASCII = `
░█▀▀░█▀█░▀█▀░█▀█░█▀▄░░░█▀▄░█▀▀░█▀▀░
░█░█░█▀█░░█░░█░█░█▀▄░░░█▀▄░▀▀█░▀▀█░
░▀▀▀░▀░▀░░▀░░▀▀▀░▀░▀░░░▀░▀░▀▀▀░▀▀▀░
`
