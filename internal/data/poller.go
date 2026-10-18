package data

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SnapshotMsg struct {
	Snapshot Snapshot
}

func PollCmd(claudeHome string) tea.Cmd {
	return func() tea.Msg {
		return SnapshotMsg{Snapshot: poll(claudeHome)}
	}
}

func TickCmd(claudeHome string, interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return SnapshotMsg{Snapshot: poll(claudeHome)}
	})
}

type Poller struct {
	claudeHome string
	program    *tea.Program
	fastTick   time.Duration
	slowTick   time.Duration
	stopCh     chan struct{}
	lastSlow   time.Time
	lastSnap   Snapshot
}

func NewPoller(claudeHome string, p *tea.Program) *Poller {
	return &Poller{
		claudeHome: claudeHome,
		program:    p,
		fastTick:   2 * time.Second,
		slowTick:   5 * time.Second,
		stopCh:     make(chan struct{}),
	}
}

func (p *Poller) Start() {
	go p.run()
}

func (p *Poller) Stop() {
	close(p.stopCh)
}

func (p *Poller) run() {
	ticker := time.NewTicker(p.fastTick)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			snap := p.doPoll()
			p.program.Send(SnapshotMsg{Snapshot: snap})
		}
	}
}

func (p *Poller) doPoll() Snapshot {
	snap := Snapshot{FetchedAt: time.Now()}

	sessions, err := FetchAgents(p.claudeHome)
	if err != nil {
		snap.Error = err
		snap.Jobs = p.lastSnap.Jobs
		snap.Daemon = p.lastSnap.Daemon
		snap.Sessions = p.lastSnap.Sessions
		p.lastSnap = snap
		return snap
	}
	snap.Sessions = sessions

	now := time.Now()
	if now.Sub(p.lastSlow) >= p.slowTick {
		p.lastSlow = now
		jobs, _ := LoadJobs(p.claudeHome)
		snap.Jobs = jobs
		snap.Daemon = LoadDaemonInfo(p.claudeHome)
		for i := range snap.Sessions {
			snap.Sessions[i].Tokens = LoadSessionTokens(p.claudeHome, snap.Sessions[i].SessionID)
		}
	} else {
		snap.Jobs = p.lastSnap.Jobs
		snap.Daemon = p.lastSnap.Daemon
		mergeTokens(snap.Sessions, p.lastSnap.Sessions)
	}

	p.lastSnap = snap
	return snap
}

func mergeTokens(dst, src []Session) {
	idx := make(map[string]TokenUsage, len(src))
	for _, s := range src {
		if s.SessionID != "" {
			idx[s.SessionID] = s.Tokens
		}
	}
	for i := range dst {
		if t, ok := idx[dst[i].SessionID]; ok {
			dst[i].Tokens = t
		}
	}
}

func poll(claudeHome string) Snapshot {
	snap := Snapshot{FetchedAt: time.Now()}
	sessions, err := FetchAgents(claudeHome)
	if err != nil {
		snap.Error = err
		return snap
	}
	snap.Sessions = sessions
	snap.Jobs, _ = LoadJobs(claudeHome)
	snap.Daemon = LoadDaemonInfo(claudeHome)
	for i := range snap.Sessions {
		snap.Sessions[i].Tokens = LoadSessionTokens(claudeHome, snap.Sessions[i].SessionID)
	}
	return snap
}
