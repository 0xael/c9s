package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/elalwany/c9s/internal/config"
	"github.com/elalwany/c9s/internal/control"
	"github.com/elalwany/c9s/internal/data"
	. "github.com/elalwany/c9s/internal/ui"
	"github.com/elalwany/c9s/internal/views/dispatch"
	"github.com/elalwany/c9s/internal/views/help"
	"github.com/elalwany/c9s/internal/views/transcript"
	"github.com/elalwany/c9s/internal/widgets"
)

const (
	TabSessions = iota
	TabJobs
	TabDashboard
	TabHistory
	tabCount
)

var tabNames = []string{"Sessions", "Jobs", "Dashboard", "History"}

type viewMode int

const (
	modeList viewMode = iota
	modeDetail
	modeTranscript
)

type Model struct {
	cfg    config.Config
	width  int
	height int
	ready  bool

	activeTab int
	mode      viewMode

	snapshot   data.Snapshot
	sessCursor int
	jobsCursor int
	showDetail bool

	header    widgets.Header
	statusBar widgets.StatusBar
	confirm   widgets.Confirm
	dispatch  dispatch.Model
	help      help.Model
	transcript transcript.Model
	toast      string
	toastTimer int

	attachCmd string
}

func New(cfg config.Config) Model {
	return Model{
		cfg:       cfg,
		activeTab: TabSessions,
		header:    widgets.NewHeader(tabNames),
		statusBar: widgets.NewStatusBar(),
		confirm:   widgets.NewConfirm(),
		dispatch:  dispatch.New(),
		help:      help.New(),
		transcript: transcript.New(),
	}
}

func (m Model) AttachCmd() string { return m.attachCmd }

func (m Model) Init() tea.Cmd {
	return data.PollCmd(m.cfg.Claude.Home)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.header = m.header.SetWidth(msg.Width)
		m.statusBar = m.statusBar.SetWidth(msg.Width)
		m.confirm = m.confirm.SetWidth(msg.Width)
		m.dispatch = m.dispatch.SetWidth(msg.Width)
		m.help = m.help.SetSize(msg.Width, msg.Height)
		m.transcript = m.transcript.SetSize(msg.Width, m.contentHeight())
		return m, nil

	case data.SnapshotMsg:
		m.snapshot = msg.Snapshot
		m.statusBar = m.statusBar.Update(msg.Snapshot)
		m.header = m.header.SetActive(m.activeTab)
		m.clampCursors()
		return m, data.TickCmd(m.cfg.Claude.Home, m.cfg.Refresh.Fast)

	case ConfirmResultMsg:
		if msg.Confirmed {
			switch msg.Action {
			case ConfirmStop:
				cmds = append(cmds, stopCmd(msg.ID))
			case ConfirmRemove:
				cmds = append(cmds, removeCmd(msg.ID))
			}
		}
		return m, tea.Batch(cmds...)

	case DispatchResultMsg:
		if msg.Err != nil {
			m.toast = fmt.Sprintf("Error: %v", msg.Err)
		} else {
			m.toast = fmt.Sprintf("Launched: %s", msg.ID)
		}
		m.toastTimer = 5
		return m, data.PollCmd(m.cfg.Claude.Home)

	case ToastMsg:
		m.toast = msg.Text
		m.toastTimer = 5
		return m, nil

	case transcriptLoadedMsg:
		m.transcript = m.transcript.SetEntries(msg.entries)
		m.transcript = m.transcript.SetSize(m.width, m.contentHeight())
		m.mode = modeTranscript
		return m, nil

	case transcript.BackMsg:
		m.mode = modeList
		return m, nil

	case help.CloseMsg:
		return m, nil

	case tickToastMsg:
		if m.toastTimer > 0 {
			m.toastTimer--
			if m.toastTimer == 0 {
				m.toast = ""
			}
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickToastMsg{} })
		}
		return m, nil
	}

	if m.confirm.Active() {
		var cmd tea.Cmd
		m.confirm, cmd = m.confirm.Update(msg)
		return m, cmd
	}

	if m.dispatch.Active() {
		var cmd tea.Cmd
		m.dispatch, cmd = m.dispatch.Update(msg)
		return m, cmd
	}

	if m.help.Visible() {
		var cmd tea.Cmd
		m.help, cmd = m.help.Update(msg)
		return m, cmd
	}

	if m.mode == modeTranscript {
		var cmd tea.Cmd
		m.transcript, cmd = m.transcript.Update(msg)
		return m, cmd
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKey(keyMsg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, Keys.Help):
		m.help = m.help.Toggle()
		return m, nil
	case key.Matches(msg, Keys.Tab1):
		m.activeTab = TabSessions
		m.mode = modeList
		m.header = m.header.SetActive(m.activeTab)
		return m, nil
	case key.Matches(msg, Keys.Tab2):
		m.activeTab = TabJobs
		m.mode = modeList
		m.header = m.header.SetActive(m.activeTab)
		return m, nil
	case key.Matches(msg, Keys.Tab3):
		m.activeTab = TabDashboard
		m.mode = modeList
		m.header = m.header.SetActive(m.activeTab)
		return m, nil
	case key.Matches(msg, Keys.Tab4):
		m.activeTab = TabHistory
		m.mode = modeList
		m.header = m.header.SetActive(m.activeTab)
		return m, nil
	case key.Matches(msg, Keys.NextTab):
		m.activeTab = (m.activeTab + 1) % tabCount
		m.mode = modeList
		m.header = m.header.SetActive(m.activeTab)
		return m, nil
	case key.Matches(msg, Keys.PrevTab):
		m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
		m.mode = modeList
		m.header = m.header.SetActive(m.activeTab)
		return m, nil
	case key.Matches(msg, Keys.Refresh):
		return m, data.PollCmd(m.cfg.Claude.Home)
	case key.Matches(msg, Keys.Dispatch):
		m.dispatch = m.dispatch.Show()
		return m, nil
	case key.Matches(msg, Keys.Down):
		m.moveCursor(1)
		return m, nil
	case key.Matches(msg, Keys.Up):
		m.moveCursor(-1)
		return m, nil
	case key.Matches(msg, Keys.Top):
		m.setCursor(0)
		return m, nil
	case key.Matches(msg, Keys.Bottom):
		m.setCursor(m.listLen() - 1)
		return m, nil
	case key.Matches(msg, Keys.Enter):
		m.showDetail = !m.showDetail
		return m, nil
	case key.Matches(msg, Keys.Back):
		if m.showDetail {
			m.showDetail = false
		}
		return m, nil
	case key.Matches(msg, Keys.Stop):
		if s := m.selectedSession(); s != nil {
			id := s.ID
			if id == "" {
				id = fmt.Sprintf("%d", s.PID)
			}
			m.confirm = m.confirm.Show(ConfirmStop, id, s.Name)
		}
		return m, nil
	case key.Matches(msg, Keys.Remove):
		if m.activeTab == TabJobs {
			if j := m.selectedJob(); j != nil {
				m.confirm = m.confirm.Show(ConfirmRemove, j.ID, j.Name)
			}
		}
		return m, nil
	case key.Matches(msg, Keys.Respawn):
		if m.activeTab == TabJobs {
			if j := m.selectedJob(); j != nil {
				return m, respawnCmd(j.ID)
			}
		}
		return m, nil
	case key.Matches(msg, Keys.Attach):
		if s := m.selectedSession(); s != nil {
			id := s.ID
			if id == "" {
				id = fmt.Sprintf("%d", s.PID)
			}
			m.attachCmd = control.AttachCommand(id)
			return m, tea.Quit
		}
		return m, nil
	case key.Matches(msg, Keys.Transcript):
		if s := m.selectedSession(); s != nil {
			return m, m.loadTranscript(s.SessionID)
		}
		if m.activeTab == TabJobs {
			if j := m.selectedJob(); j != nil {
				return m, m.loadTranscript(j.SessionID)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading c9s..."
	}

	if m.mode == modeTranscript {
		return m.transcript.View()
	}

	header := m.header.View()
	statusBar := m.statusBar.View()

	headerH := lipgloss.Height(header)
	statusH := lipgloss.Height(statusBar)
	contentH := m.height - headerH - statusH
	if contentH < 1 {
		contentH = 1
	}

	content := m.contentView(contentH)

	view := lipgloss.JoinVertical(lipgloss.Left, header, content, statusBar)

	if m.help.Visible() {
		return m.help.View()
	}
	if m.confirm.Active() {
		overlay := m.confirm.View()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, overlay)
	}
	if m.dispatch.Active() {
		overlay := m.dispatch.View()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, overlay)
	}

	if m.toast != "" {
		toastLine := ToastStyle.Render(" " + m.toast + " ")
		view = view + "\n" + toastLine
	}

	return view
}

func (m Model) contentView(height int) string {
	style := lipgloss.NewStyle().Width(m.width).Height(height)

	switch m.activeTab {
	case TabSessions:
		return style.Render(m.sessionsView(height))
	case TabJobs:
		return style.Render(m.jobsView(height))
	case TabDashboard:
		return style.Render("\n  Dashboard (coming in Phase 2)")
	case TabHistory:
		return style.Render("\n  History (coming in Phase 2)")
	}
	return ""
}

func (m Model) sessionsView(height int) string {
	sessions := m.snapshot.Sessions
	if len(sessions) == 0 {
		return "\n  No sessions found. Is Claude Code running?"
	}

	if m.showDetail && m.sessCursor < len(sessions) {
		return m.sessionsSplitView(sessions, height)
	}

	return m.sessionsTableView(sessions, height)
}

func (m Model) sessionsTableView(sessions []data.Session, height int) string {
	header := TableHeader.Render(fmt.Sprintf("  %-22s %-8s %-4s %-22s %-10s %10s",
		"NAME", "STATUS", "KIND", "CWD", "AGE", "TOKENS"))

	var rows []string
	for i, s := range sessions {
		kind := "int"
		if s.Kind == "background" {
			kind = "bg"
		}
		status := s.Status
		if s.State != "" {
			status = s.State
		}
		cwd := truncPath(s.CWD, 22)
		age := shortDuration(time.Since(time.UnixMilli(s.StartedAt)))
		stStyle := StatusStyle(s.Status, s.State)

		row := fmt.Sprintf("  %-22s %s %-4s %-22s %-10s %10s",
			trunc(s.Name, 22),
			stStyle.Width(8).Render(status),
			kind, cwd, age, FmtTokens(s.Tokens.Total()))

		if i == m.sessCursor {
			rows = append(rows, TableRowSelected.Width(m.width).Render(row))
		} else {
			rows = append(rows, TableRow.Render(row))
		}
	}

	return header + "\n" + strings.Join(rows, "\n")
}

func (m Model) sessionsSplitView(sessions []data.Session, height int) string {
	borderW := 2
	leftInner := m.width*35/100 - borderW
	rightInner := m.width - leftInner - borderW*2 - 3
	if leftInner < 10 {
		leftInner = 10
	}
	if rightInner < 20 {
		rightInner = 20
	}
	paneH := height - 2
	if paneH < 3 {
		paneH = 3
	}

	var leftLines []string
	for i, s := range sessions {
		status := s.Status
		if s.State != "" {
			status = s.State
		}
		stStyle := StatusStyle(s.Status, s.State)
		nameW := leftInner - 10
		if nameW < 5 {
			nameW = 5
		}
		line := fmt.Sprintf(" %-*s %s", nameW, trunc(s.Name, nameW), stStyle.Render(status))
		if i == m.sessCursor {
			leftLines = append(leftLines, TableRowSelected.Render(line))
		} else {
			leftLines = append(leftLines, TableRow.Render(line))
		}
	}
	left := PaneBorder.Width(leftInner).Height(paneH).Render(strings.Join(leftLines, "\n"))

	s := sessions[m.sessCursor]
	var detail []string
	detail = append(detail, DetailTitle.Render(s.Name))
	detail = append(detail, "")
	detail = append(detail, detailRow("PID", fmt.Sprintf("%d", s.PID)))
	detail = append(detail, detailRow("Session ID", s.SessionID))
	detail = append(detail, detailRow("CWD", s.CWD))
	detail = append(detail, detailRow("Kind", s.Kind))
	detail = append(detail, detailRow("Status", s.Status))
	if s.State != "" {
		detail = append(detail, detailRow("State", s.State))
	}
	if s.Version != "" {
		detail = append(detail, detailRow("Version", s.Version))
	}
	if s.Model != "" {
		detail = append(detail, detailRow("Model", s.Model))
	}
	age := humanize.Time(time.UnixMilli(s.StartedAt))
	detail = append(detail, detailRow("Started", age))
	if s.Tokens.Total() > 0 {
		detail = append(detail, "")
		detail = append(detail, DetailTitle.Render("Token Usage"))
		detail = append(detail, detailRow("Input", FmtTokens(s.Tokens.InputTokens)))
		detail = append(detail, detailRow("Output", FmtTokens(s.Tokens.OutputTokens)))
		detail = append(detail, detailRow("Cache Write", FmtTokens(s.Tokens.CacheCreationInputTokens)))
		detail = append(detail, detailRow("Cache Read", FmtTokens(s.Tokens.CacheReadInputTokens)))
		if s.Tokens.ThinkingTokens > 0 {
			detail = append(detail, detailRow("Thinking", FmtTokens(s.Tokens.ThinkingTokens)))
		}
		detail = append(detail, detailRow("Total", FmtTokens(s.Tokens.Total())))
	}
	detail = append(detail, "")
	detail = append(detail, HelpKey.Render("[t]")+" transcript  "+HelpKey.Render("[s]")+" stop  "+HelpKey.Render("[a]")+" attach")

	right := PaneBorderActive.Width(rightInner).Height(paneH).Render(strings.Join(detail, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)
}

func (m Model) jobsView(height int) string {
	jobs := m.snapshot.Jobs
	if len(jobs) == 0 {
		return "\n  No background jobs found."
	}

	if m.showDetail && m.jobsCursor < len(jobs) {
		return m.jobsSplitView(jobs, height)
	}

	return m.jobsTableView(jobs, height)
}

func (m Model) jobsTableView(jobs []data.Job, height int) string {
	header := TableHeader.Render(fmt.Sprintf("  %-10s %-22s %-10s %-8s %-30s",
		"ID", "NAME", "STATE", "TOKENS", "DETAIL"))

	var rows []string
	for i, j := range jobs {
		stStyle := StatusStyle("", j.State)
		row := fmt.Sprintf("  %-10s %-22s %s %-8d %-30s",
			trunc(j.ID, 10),
			trunc(j.Name, 22),
			stStyle.Width(10).Render(j.State),
			j.Tokens,
			trunc(j.Detail, 30))

		if i == m.jobsCursor {
			rows = append(rows, TableRowSelected.Width(m.width).Render(row))
		} else {
			rows = append(rows, TableRow.Render(row))
		}
	}

	return header + "\n" + strings.Join(rows, "\n")
}

func (m Model) jobsSplitView(jobs []data.Job, height int) string {
	borderW := 2
	leftInner := m.width*35/100 - borderW
	rightInner := m.width - leftInner - borderW*2 - 3
	if leftInner < 10 {
		leftInner = 10
	}
	if rightInner < 20 {
		rightInner = 20
	}
	paneH := height - 2
	if paneH < 3 {
		paneH = 3
	}

	var leftLines []string
	for i, j := range jobs {
		stStyle := StatusStyle("", j.State)
		nameW := leftInner - 10
		if nameW < 5 {
			nameW = 5
		}
		line := fmt.Sprintf(" %-*s %s", nameW, trunc(j.Name, nameW), stStyle.Render(j.State))
		if i == m.jobsCursor {
			leftLines = append(leftLines, TableRowSelected.Render(line))
		} else {
			leftLines = append(leftLines, TableRow.Render(line))
		}
	}
	left := PaneBorder.Width(leftInner).Height(paneH).Render(strings.Join(leftLines, "\n"))

	j := jobs[m.jobsCursor]
	var detail []string
	detail = append(detail, DetailTitle.Render(j.Name))
	detail = append(detail, "")
	detail = append(detail, detailRow("ID", j.ID))
	detail = append(detail, detailRow("State", j.State))
	detail = append(detail, detailRow("Tokens", fmt.Sprintf("%d", j.Tokens)))
	detail = append(detail, detailRow("Intent", j.Intent))
	detail = append(detail, detailRow("CWD", j.CWD))
	if j.CreatedAt != "" {
		detail = append(detail, detailRow("Created", j.CreatedAt))
	}
	if j.Output != nil && j.Output.Result != "" {
		result := trunc(j.Output.Result, 60)
		detail = append(detail, detailRow("Result", result))
	}

	if len(j.Timeline) > 0 {
		detail = append(detail, "")
		detail = append(detail, DetailTitle.Render("Timeline"))
		for _, e := range j.Timeline {
			ts := e.At
			if len(ts) > 19 {
				ts = ts[11:19]
			}
			stStyle := StatusStyle("", e.State)
			detail = append(detail, fmt.Sprintf("  %s  %s  %s",
				HelpDesc.Render(ts),
				stStyle.Width(8).Render(e.State),
				trunc(e.Detail, 40)))
		}
	}

	detail = append(detail, "")
	detail = append(detail, HelpKey.Render("[t]")+" transcript  "+HelpKey.Render("[x]")+" remove  "+HelpKey.Render("[R]")+" respawn")

	right := PaneBorderActive.Width(rightInner).Height(paneH).Render(strings.Join(detail, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)
}

// Helpers

func (m *Model) moveCursor(delta int) {
	switch m.activeTab {
	case TabSessions:
		m.sessCursor += delta
		m.clampCursors()
	case TabJobs:
		m.jobsCursor += delta
		m.clampCursors()
	}
}

func (m *Model) setCursor(pos int) {
	switch m.activeTab {
	case TabSessions:
		m.sessCursor = pos
	case TabJobs:
		m.jobsCursor = pos
	}
	m.clampCursors()
}

func (m *Model) clampCursors() {
	if n := len(m.snapshot.Sessions); n > 0 {
		if m.sessCursor >= n {
			m.sessCursor = n - 1
		}
		if m.sessCursor < 0 {
			m.sessCursor = 0
		}
	} else {
		m.sessCursor = 0
	}
	if n := len(m.snapshot.Jobs); n > 0 {
		if m.jobsCursor >= n {
			m.jobsCursor = n - 1
		}
		if m.jobsCursor < 0 {
			m.jobsCursor = 0
		}
	} else {
		m.jobsCursor = 0
	}
}

func (m Model) listLen() int {
	switch m.activeTab {
	case TabSessions:
		return len(m.snapshot.Sessions)
	case TabJobs:
		return len(m.snapshot.Jobs)
	}
	return 0
}

func (m Model) selectedSession() *data.Session {
	if m.activeTab != TabSessions || len(m.snapshot.Sessions) == 0 {
		return nil
	}
	if m.sessCursor >= len(m.snapshot.Sessions) {
		return nil
	}
	s := m.snapshot.Sessions[m.sessCursor]
	return &s
}

func (m Model) selectedJob() *data.Job {
	if m.activeTab != TabJobs || len(m.snapshot.Jobs) == 0 {
		return nil
	}
	if m.jobsCursor >= len(m.snapshot.Jobs) {
		return nil
	}
	j := m.snapshot.Jobs[m.jobsCursor]
	return &j
}

func (m Model) contentHeight() int {
	h := m.height - 3
	if h < 1 {
		h = 1
	}
	return h
}

func (m Model) loadTranscript(sessionID string) tea.Cmd {
	claudeHome := m.cfg.Claude.Home
	maxLines := m.cfg.UI.TranscriptInitialLines
	return func() tea.Msg {
		path, err := data.FindTranscriptPath(claudeHome, sessionID)
		if err != nil {
			return ToastMsg{Text: fmt.Sprintf("Transcript not found: %v", err)}
		}
		entries, err := data.LoadTranscript(path, maxLines)
		if err != nil {
			return ToastMsg{Text: fmt.Sprintf("Error loading transcript: %v", err)}
		}
		return transcriptLoadedMsg{sessionID: sessionID, entries: entries}
	}
}

type transcriptLoadedMsg struct {
	sessionID string
	entries   []data.TranscriptEntry
}

type tickToastMsg struct{}

func stopCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := control.StopSession(context.Background(), id)
		if err != nil {
			return ToastMsg{Text: fmt.Sprintf("Stop failed: %v", err)}
		}
		return ToastMsg{Text: fmt.Sprintf("Stopped: %s", id)}
	}
}

func removeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := control.RemoveJob(context.Background(), id)
		if err != nil {
			return ToastMsg{Text: fmt.Sprintf("Remove failed: %v", err)}
		}
		return ToastMsg{Text: fmt.Sprintf("Removed: %s", id)}
	}
}

func respawnCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := control.RespawnJob(context.Background(), id)
		if err != nil {
			return ToastMsg{Text: fmt.Sprintf("Respawn failed: %v", err)}
		}
		return ToastMsg{Text: fmt.Sprintf("Respawned: %s", id)}
	}
}

func detailRow(label, value string) string {
	return fmt.Sprintf("  %s %s", DetailLabel.Render(label+":"), DetailValue.Render(value))
}

func trunc(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func truncPath(path string, max int) string {
	if len(path) <= max {
		return path
	}
	base := filepath.Base(path)
	if base != "." && len(base) <= max-2 {
		return "~/" + base
	}
	return trunc(path, max)
}

func shortDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

