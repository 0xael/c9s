package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Primary      lipgloss.Color
	Secondary    lipgloss.Color
	Accent       lipgloss.Color
	Muted        lipgloss.Color
	Success      lipgloss.Color
	Warning      lipgloss.Color
	Error        lipgloss.Color
	Border       lipgloss.Color
	ActiveBorder lipgloss.Color
	Text         lipgloss.Color
	DimText      lipgloss.Color
	BgDark       lipgloss.Color
}

var DefaultTheme = Theme{
	Primary:      lipgloss.Color("#00D7FF"),
	Secondary:    lipgloss.Color("#87AFFF"),
	Accent:       lipgloss.Color("#FF5FAF"),
	Muted:        lipgloss.Color("#6C7086"),
	Success:      lipgloss.Color("#A6E3A1"),
	Warning:      lipgloss.Color("#F9E2AF"),
	Error:        lipgloss.Color("#F38BA8"),
	Border:       lipgloss.Color("#45475A"),
	ActiveBorder: lipgloss.Color("#00D7FF"),
	Text:         lipgloss.Color("#CDD6F4"),
	DimText:      lipgloss.Color("#6C7086"),
	BgDark:       lipgloss.Color("#1E1E2E"),
}

var T = DefaultTheme

var (
	TabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(T.Primary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(T.Primary).
			Padding(0, 2)

	TabInactive = lipgloss.NewStyle().
			Foreground(T.DimText).
			Padding(0, 2)

	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(T.Secondary).
			Padding(0, 1)

	TableRow = lipgloss.NewStyle().
			Foreground(T.Text).
			Padding(0, 1)

	TableRowSelected = lipgloss.NewStyle().
				Bold(true).
				Foreground(T.Primary).
				Background(lipgloss.Color("#313244")).
				Padding(0, 1)

	StatusBar = lipgloss.NewStyle().
			Foreground(T.DimText).
			Background(lipgloss.Color("#181825")).
			Padding(0, 1)

	DetailTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(T.Primary).
			MarginBottom(1)

	DetailLabel = lipgloss.NewStyle().
			Foreground(T.Muted).
			Width(14)

	DetailValue = lipgloss.NewStyle().
			Foreground(T.Text)

	HelpKey = lipgloss.NewStyle().
		Bold(true).
		Foreground(T.Accent)

	HelpDesc = lipgloss.NewStyle().
			Foreground(T.DimText)

	ToastStyle = lipgloss.NewStyle().
			Foreground(T.Text).
			Background(lipgloss.Color("#313244")).
			Padding(0, 2).
			Bold(true)

	ConfirmBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(T.Warning).
			Padding(1, 3).
			Width(50)

	StatusBusy = lipgloss.NewStyle().Foreground(T.Success)
	StatusIdle = lipgloss.NewStyle().Foreground(T.Muted)
	StatusDone = lipgloss.NewStyle().Foreground(T.Secondary)
	StatusErr  = lipgloss.NewStyle().Foreground(T.Error)

	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(T.Primary)

	PaneBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(T.Border)

	PaneBorderActive = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(T.ActiveBorder)
)

func StatusStyle(status, state string) lipgloss.Style {
	switch {
	case status == "busy":
		return StatusBusy
	case state == "done":
		return StatusDone
	case state == "error":
		return StatusErr
	case status == "idle":
		return StatusIdle
	default:
		return StatusIdle
	}
}

func FmtTokens(n int64) string {
	if n == 0 {
		return "-"
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}
