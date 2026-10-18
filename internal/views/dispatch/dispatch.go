package dispatch

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/0xael/c9s/internal/ui"
	"github.com/0xael/c9s/internal/control"
)

type Model struct {
	input  textinput.Model
	active bool
	width  int
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter a prompt for the background agent..."
	ti.CharLimit = 1000
	ti.Width = 60
	return Model{input: ti}
}

func (m Model) Show() Model {
	m.active = true
	m.input.Focus()
	m.input.SetValue("")
	return m
}

func (m Model) Hide() Model {
	m.active = false
	m.input.Blur()
	return m
}

func (m Model) Active() bool {
	return m.active
}

func (m Model) SetWidth(w int) Model {
	m.width = w
	inputW := w - 10
	if inputW > 80 {
		inputW = 80
	}
	if inputW < 30 {
		inputW = 30
	}
	m.input.Width = inputW
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.active {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.Hide(), nil
		case "enter":
			prompt := m.input.Value()
			if prompt == "" {
				return m, nil
			}
			m = m.Hide()
			return m, launchCmd(prompt)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func launchCmd(prompt string) tea.Cmd {
	return func() tea.Msg {
		id, err := control.LaunchBackground(context.Background(), prompt, "")
		return ui.DispatchResultMsg{ID: id, Err: err}
	}
}

func (m Model) View() string {
	if !m.active {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.T.Primary).
		Render("Dispatch Background Agent")

	hint := lipgloss.NewStyle().
		Foreground(ui.T.DimText).
		Render("[enter: launch | esc: cancel]")

	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.T.Primary).
		Padding(1, 3).
		Width(m.input.Width + 10).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			m.input.View(),
			"",
			hint,
		))

	return lipgloss.Place(m.width, 10, lipgloss.Center, lipgloss.Center, box)
}
