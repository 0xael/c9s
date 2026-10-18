package data

import (
	"encoding/json"
	"time"
)

type TokenUsage struct {
	InputTokens              int64
	OutputTokens             int64
	CacheCreationInputTokens int64
	CacheReadInputTokens     int64
	ThinkingTokens           int64
}

func (t TokenUsage) Total() int64 {
	return t.InputTokens + t.OutputTokens + t.CacheCreationInputTokens + t.CacheReadInputTokens
}

type Session struct {
	PID       int    `json:"pid"`
	SessionID string `json:"sessionId"`
	CWD       string `json:"cwd"`
	Kind      string `json:"kind"`
	Status    string `json:"status"`
	State     string `json:"state"`
	Name      string `json:"name"`
	StartedAt int64  `json:"startedAt"`
	ID        string `json:"id"`
	Version    string `json:"-"`
	Model      string `json:"-"`
	SocketPath string `json:"-"`
	Tokens     TokenUsage `json:"-"`
}

type Job struct {
	ID         string          `json:"-"`
	State      string          `json:"state"`
	Detail     string          `json:"detail"`
	Tempo      string          `json:"tempo"`
	Tokens     int             `json:"tokens"`
	Intent     string          `json:"intent"`
	Name       string          `json:"name"`
	SessionID  string          `json:"sessionId"`
	CWD        string          `json:"cwd"`
	CLIVersion string          `json:"cliVersion"`
	CreatedAt  string          `json:"createdAt"`
	UpdatedAt  string          `json:"updatedAt"`
	Output     *JobOutput      `json:"output"`
	InFlight   InFlight        `json:"inFlight"`
	Timeline   []TimelineEntry `json:"-"`
	Pinned     bool            `json:"-"`
}

type JobOutput struct {
	Result string `json:"result"`
}

type InFlight struct {
	Tasks  int `json:"tasks"`
	Queued int `json:"queued"`
}

type TimelineEntry struct {
	At     string `json:"at"`
	State  string `json:"state"`
	Detail string `json:"detail"`
	Text   string `json:"text"`
}

type DaemonInfo struct {
	PID       int
	Version   string
	StartedAt int64
	Running   bool
	Workers   int
}

type Snapshot struct {
	Sessions  []Session
	Jobs      []Job
	Daemon    DaemonInfo
	FetchedAt time.Time
	Error     error
}

type TranscriptEntry struct {
	Type      string          `json:"type"`
	Timestamp string          `json:"timestamp"`
	Message   json.RawMessage `json:"message"`
}
