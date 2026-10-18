package data

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func LoadHistory(claudeHome string) ([]HistorySession, error) {
	path := filepath.Join(claudeHome, "history.jsonl")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	grouped := make(map[string]*HistorySession)
	var order []string

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 256*1024), 256*1024)
	for sc.Scan() {
		var e HistoryEntry
		if json.Unmarshal(sc.Bytes(), &e) != nil || e.SessionID == "" {
			continue
		}

		ts := time.UnixMilli(e.Timestamp)
		hs, ok := grouped[e.SessionID]
		if !ok {
			hs = &HistorySession{
				SessionID:   e.SessionID,
				Project:     e.Project,
				FirstPrompt: e.Display,
				FirstAt:     ts,
				LastAt:      ts,
			}
			grouped[e.SessionID] = hs
			order = append(order, e.SessionID)
		}

		hs.PromptCount++
		if ts.After(hs.LastAt) {
			hs.LastAt = ts
		}
		if ts.Before(hs.FirstAt) {
			hs.FirstAt = ts
			hs.FirstPrompt = e.Display
		}
		if len(hs.Prompts) < 20 {
			hs.Prompts = append(hs.Prompts, e)
		}
	}

	sessions := make([]HistorySession, 0, len(grouped))
	for _, sid := range order {
		sessions = append(sessions, *grouped[sid])
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastAt.After(sessions[j].LastAt)
	})

	return sessions, nil
}
