package transcript

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elalwany/c9s/internal/ui"
	"github.com/elalwany/c9s/internal/data"
)

type BackMsg struct{}

type Model struct {
	entries       []data.TranscriptEntry
	lines         []styledLine
	offset        int
	follow        bool
	width, height int
	sessionID     string
}

func New() Model {
	return Model{follow: true}
}

func (m Model) SetSize(w, h int) Model {
	m.width = w
	m.height = h
	if len(m.entries) > 0 {
		m.lines = parseEntries(m.entries, w)
		if m.follow {
			m.scrollToBottom()
		}
	}
	return m
}

func (m Model) SetEntries(entries []data.TranscriptEntry) Model {
	m.entries = entries
	m.lines = parseEntries(entries, m.width)
	if m.follow {
		m.scrollToBottom()
	}
	return m
}

func (m *Model) scrollToBottom() {
	max := len(m.lines) - m.height
	if max < 0 {
		max = 0
	}
	m.offset = max
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ui.Keys.Back), key.Matches(msg, ui.Keys.Quit):
			return m, func() tea.Msg { return BackMsg{} }
		case key.Matches(msg, ui.Keys.Follow):
			m.follow = !m.follow
			if m.follow {
				m.scrollToBottom()
			}
		case key.Matches(msg, ui.Keys.Down):
			m.follow = false
			m.offset++
		case key.Matches(msg, ui.Keys.Up):
			m.follow = false
			m.offset--
		case key.Matches(msg, ui.Keys.PageDown):
			m.follow = false
			m.offset += m.height / 2
		case key.Matches(msg, ui.Keys.PageUp):
			m.follow = false
			m.offset -= m.height / 2
		case key.Matches(msg, ui.Keys.Top):
			m.follow = false
			m.offset = 0
		case key.Matches(msg, ui.Keys.Bottom):
			m.follow = true
			m.scrollToBottom()
		}
		m.clampOffset()
	}
	return m, nil
}

func (m *Model) clampOffset() {
	max := len(m.lines) - m.height
	if max < 0 {
		max = 0
	}
	if m.offset > max {
		m.offset = max
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if len(m.lines) == 0 {
		empty := lipgloss.NewStyle().Foreground(ui.T.DimText).Render("  No conversation messages found in transcript.")
		statusLine := lipgloss.NewStyle().
			Foreground(ui.T.DimText).
			Background(lipgloss.Color("#181825")).
			Width(m.width).
			Render(" Transcript [esc/q:back]")
		return lipgloss.JoinVertical(lipgloss.Left, "\n\n"+empty, statusLine)
	}

	visibleLines := m.height - 1
	end := m.offset + visibleLines
	if end > len(m.lines) {
		end = len(m.lines)
	}

	var sb strings.Builder
	start := m.offset
	if start < 0 {
		start = 0
	}
	for i := start; i < end; i++ {
		sb.WriteString(m.lines[i].text)
		if i < end-1 {
			sb.WriteByte('\n')
		}
	}

	for i := end - start; i < visibleLines; i++ {
		sb.WriteByte('\n')
	}

	content := sb.String()

	followIndicator := " "
	if m.follow {
		followIndicator = lipgloss.NewStyle().Foreground(ui.T.Success).Render("FOLLOW")
	}
	pos := fmt.Sprintf(" %d/%d ", m.offset+1, len(m.lines))
	statusLine := lipgloss.NewStyle().
		Foreground(ui.T.DimText).
		Background(lipgloss.Color("#181825")).
		Width(m.width).
		Render(fmt.Sprintf(" Transcript [esc/q:back f:follow j/k:scroll]  %s  %s", followIndicator, pos))

	return lipgloss.JoinVertical(lipgloss.Left, content, statusLine)
}
