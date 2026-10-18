package web

import "github.com/0xael/c9s/internal/data"

type WebSession struct {
	PID       int            `json:"pid"`
	SessionID string         `json:"sessionId"`
	CWD       string         `json:"cwd"`
	Kind      string         `json:"kind"`
	Status    string         `json:"status"`
	State     string         `json:"state"`
	Name      string         `json:"name"`
	StartedAt int64          `json:"startedAt"`
	ID        string         `json:"id"`
	Version   string         `json:"version"`
	Model     string         `json:"model"`
	Tokens    data.TokenUsage `json:"tokens"`
}

func toWebSession(s data.Session) WebSession {
	return WebSession{
		PID:       s.PID,
		SessionID: s.SessionID,
		CWD:       s.CWD,
		Kind:      s.Kind,
		Status:    s.Status,
		State:     s.State,
		Name:      s.Name,
		StartedAt: s.StartedAt,
		ID:        s.ID,
		Version:   s.Version,
		Model:     s.Model,
		Tokens:    s.Tokens,
	}
}

type OverviewResponse struct {
	TotalSessions    int     `json:"totalSessions"`
	TotalMessages    int     `json:"totalMessages"`
	TotalCostUSD     float64 `json:"totalCostUSD"`
	TotalTokens      int64   `json:"totalTokens"`
	FirstSessionDate string  `json:"firstSessionDate"`
	LiveSessions     int     `json:"liveSessions"`
	TotalLinesAdded  int     `json:"totalLinesAdded"`
	TotalLinesRemoved int    `json:"totalLinesRemoved"`
}

type HistoryRow struct {
	SessionID   string   `json:"sessionId"`
	Project     string   `json:"project"`
	ProjectShort string  `json:"projectShort"`
	FirstPrompt string   `json:"firstPrompt"`
	PromptCount int      `json:"promptCount"`
	FirstAt     int64    `json:"firstAt"`
	LastAt      int64    `json:"lastAt"`
	DurationMs  int64    `json:"durationMs"`
	CostUSD     float64  `json:"costUSD"`
	TotalTokens int64    `json:"totalTokens"`
	LinesAdded  int      `json:"linesAdded"`
	LinesRemoved int     `json:"linesRemoved"`
	Models      []string `json:"models"`
}

type SessionDetail struct {
	SessionID string            `json:"sessionId"`
	Project   string            `json:"project"`
	Cost      *data.SessionCost `json:"cost"`
	Prompts   []data.HistoryEntry `json:"prompts"`
}

type DailyCost struct {
	Date        string             `json:"date"`
	CostByModel map[string]float64 `json:"costByModel"`
}

type CostResponse struct {
	Daily        []DailyCost        `json:"daily"`
	TotalByModel map[string]float64 `json:"totalByModel"`
}

type ModelInfo struct {
	Name         string  `json:"name"`
	InputTokens  int64   `json:"inputTokens"`
	OutputTokens int64   `json:"outputTokens"`
	CacheRead    int64   `json:"cacheRead"`
	CacheWrite   int64   `json:"cacheWrite"`
	TotalTokens  int64   `json:"totalTokens"`
	CostUSD      float64 `json:"costUSD"`
}
