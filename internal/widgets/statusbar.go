package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/elalwany/c9s/internal/ui"
	"github.com/elalwany/c9s/internal/data"
)

type StatusBar struct {
	daemon      data.DaemonInfo
	sessions    int
	jobs        int
	totalTokens int64
	errMsg      string
	fetchedAt   time.Time
	width       int
}

func NewStatusBar() StatusBar {
	return StatusBar{}
}

func (s StatusBar) Update(snap data.Snapshot) StatusBar {
	s.daemon = snap.Daemon
	s.sessions = len(snap.Sessions)
	s.jobs = len(snap.Jobs)
	s.fetchedAt = snap.FetchedAt
	var total int64
	for _, sess := range snap.Sessions {
		total += sess.Tokens.Total()
	}
	s.totalTokens = total
	if snap.Error != nil {
		s.errMsg = snap.Error.Error()
	} else {
		s.errMsg = ""
	}
	return s
}

func (s StatusBar) SetWidth(w int) StatusBar {
	s.width = w
	return s
}

func (s StatusBar) View() string {
	var parts []string

	if s.daemon.Running {
		parts = append(parts, fmt.Sprintf("Daemon: v%s (pid %d)", s.daemon.Version, s.daemon.PID))
	} else {
		parts = append(parts, ui.StatusErr.Render("Daemon: not running"))
	}

	parts = append(parts, fmt.Sprintf("%d sessions", s.sessions))

	if s.jobs > 0 {
		parts = append(parts, fmt.Sprintf("%d bg jobs", s.jobs))
	}

	if s.totalTokens > 0 {
		parts = append(parts, fmt.Sprintf("Tokens: %s", ui.FmtTokens(s.totalTokens)))
	}

	if !s.fetchedAt.IsZero() {
		ago := time.Since(s.fetchedAt).Truncate(time.Second)
		parts = append(parts, fmt.Sprintf("Updated %s ago", ago))
	}

	left := strings.Join(parts, " | ")

	if s.errMsg != "" {
		errText := ui.StatusErr.Render("ERR: " + s.errMsg)
		gap := s.width - lipgloss.Width(left) - lipgloss.Width(errText)
		if gap < 1 {
			gap = 1
		}
		left = left + strings.Repeat(" ", gap) + errText
	}

	return ui.StatusBar.Width(s.width).Render(left)
}
