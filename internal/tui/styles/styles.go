package styles

import (
	"charm.land/lipgloss/v2"
)

var BlueGradient = []string{
	"#1a3c6e", "#214c7d", "#285c8c", "#2f6c9b",
	"#367caa", "#3d8cb9", "#449cc8", "#4bacc7",
	"#5cbcd6", "#6dcde5", "#7ed8ea", "#8fe3ef",
	"#a0eef4", "#b0eef9", "#c0eef9", "#d0eef9",
}

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5f5fff"))
	SelectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ff5faf")).Foreground(lipgloss.Color("#eeeeee"))
	// TODO: rename and restyle
	UsernameStyle    = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(lipgloss.Color("#ba99e3"))
	CursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5faf"))
	PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	HelpStyle        = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#626262"))
	ErrorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	InputStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#add8e6"))
	PromptStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#afffff"))
	InputBorderStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Width(20).Height(2).
				Align(lipgloss.Center).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#3c71a8"))
	WideInputBorderStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Width(40).Height(2).
				Align(lipgloss.Center).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#3c71a8"))
)

const TitleASCII = `
░█▀▀░█▀█░▀█▀░█▀█░█▀▄░░░█▀▄░█▀▀░█▀▀░
░█░█░█▀█░░█░░█░█░█▀▄░░░█▀▄░▀▀█░▀▀█░
░▀▀▀░▀░▀░░▀░░▀▀▀░▀░▀░░░▀░▀░▀▀▀░▀▀▀░
`
