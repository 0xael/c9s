package help

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elalwany/c9s/internal/ui"
)

type CloseMsg struct{}

type Model struct {
	visible       bool
	width, height int
}

func New() Model {
	return Model{}
}

func (m Model) Toggle() Model {
	m.visible = !m.visible
	return m
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) SetSize(w, h int) Model {
	m.width = w
	m.height = h
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ui.Keys.Help), key.Matches(msg, ui.Keys.Back), msg.String() == "q":
			m.visible = false
			return m, func() tea.Msg { return CloseMsg{} }
		}
	}
	return m, nil
}

type section struct {
	title    string
	bindings []binding
}

type binding struct {
	keys string
	desc string
}

func (m Model) View() string {
	if !m.visible || m.width == 0 {
		return ""
	}

	sections := []section{
		{
			title: "Navigation",
			bindings: []binding{
				{"1-4", "Switch tabs"},
				{"tab/shift+tab", "Next/prev tab"},
				{"j/k, ↑/↓", "Move up/down"},
				{"g/G", "Top/bottom"},
				{"ctrl+u/ctrl+d", "Page up/down"},
				{"enter", "Open detail"},
				{"esc", "Back / close"},
			},
		},
		{
			title: "Actions",
			bindings: []binding{
				{"d", "Dispatch new agent"},
				{"s", "Stop session"},
				{"x", "Remove job"},
				{"R", "Respawn job"},
				{"a", "Attach to session"},
				{"y", "Copy session ID"},
				{"r", "Force refresh"},
			},
		},
		{
			title: "Views",
			bindings: []binding{
				{"t", "Open transcript"},
				{"f", "Toggle follow mode"},
				{"/", "Filter"},
				{"?", "Toggle help"},
				{"q", "Quit"},
			},
		},
	}

	keyStyle := ui.HelpKey.Width(18)
	descStyle := ui.HelpDesc

	var content strings.Builder
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.T.Primary).
		Render("c9s Keyboard Shortcuts")
	content.WriteString(title)
	content.WriteString("\n\n")

	for i, sec := range sections {
		secTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(ui.T.Secondary).
			Render(sec.title)
		content.WriteString(secTitle)
		content.WriteString("\n")
		for _, b := range sec.bindings {
			line := fmt.Sprintf("  %s %s",
				keyStyle.Render(b.keys),
				descStyle.Render(b.desc),
			)
			content.WriteString(line)
			content.WriteString("\n")
		}
		if i < len(sections)-1 {
			content.WriteString("\n")
		}
	}

	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.T.Border).
		Padding(1, 3).
		Render(content.String())

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
