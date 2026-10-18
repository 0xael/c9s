package data

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type costStateEnvelope struct {
	Type              string                    `json:"type"`
	SessionID         string                    `json:"sessionId"`
	TotalCostUSD      float64                   `json:"totalCostUSD"`
	TotalAPIDuration  int64                     `json:"totalAPIDuration"`
	TotalToolDuration int64                     `json:"totalToolDuration"`
	TotalDuration     int64                     `json:"totalDuration"`
	TotalLinesAdded   int                       `json:"totalLinesAdded"`
	TotalLinesRemoved int                       `json:"totalLinesRemoved"`
	StartTime         int64                     `json:"startTime"`
	ModelUsage        map[string]CostModelUsage `json:"modelUsage"`
}

func LoadAllSessionCosts(claudeHome string) ([]SessionCost, error) {
	projectsDir := filepath.Join(claudeHome, "projects")
	dirs, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, err
	}

	var costs []SessionCost
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		project := decodeProjectDir(d.Name())
		projDir := filepath.Join(projectsDir, d.Name())
		files, err := os.ReadDir(projDir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".jsonl") {
				continue
			}
			sc, ok := extractLastCostState(filepath.Join(projDir, f.Name()))
			if ok {
				sc.Project = project
				costs = append(costs, sc)
			}
		}
	}
	return costs, nil
}

func LoadSessionCost(claudeHome, sessionID string) (*SessionCost, error) {
	path, err := FindTranscriptPath(claudeHome, sessionID)
	if err != nil {
		return nil, err
	}
	sc, ok := extractLastCostState(path)
	if !ok {
		return nil, os.ErrNotExist
	}
	return &sc, nil
}

func extractLastCostState(path string) (SessionCost, bool) {
	f, err := os.Open(path)
	if err != nil {
		return SessionCost{}, false
	}
	defer f.Close()

	var last *costStateEnvelope
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) < 20 || !strings.Contains(string(line[:30]), "cost-state") {
			continue
		}
		var e costStateEnvelope
		if json.Unmarshal(line, &e) == nil && e.Type == "cost-state" {
			last = &e
		}
	}
	_ = sc.Err()

	if last == nil {
		return SessionCost{}, false
	}

	return SessionCost{
		SessionID:         last.SessionID,
		TotalCostUSD:      last.TotalCostUSD,
		TotalAPIDuration:  last.TotalAPIDuration,
		TotalToolDuration: last.TotalToolDuration,
		TotalDuration:     last.TotalDuration,
		TotalLinesAdded:   last.TotalLinesAdded,
		TotalLinesRemoved: last.TotalLinesRemoved,
		StartTime:         last.StartTime,
		ModelUsage:        last.ModelUsage,
	}, true
}

func decodeProjectDir(name string) string {
	if !strings.HasPrefix(name, "-") {
		return name
	}
	return strings.ReplaceAll(name[1:], "-", "/")
}
