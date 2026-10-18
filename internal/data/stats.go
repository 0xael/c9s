package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type statsCacheFile struct {
	DailyActivity    []DailyActivity    `json:"dailyActivity"`
	DailyModelTokens []DailyModelTokens `json:"dailyModelTokens"`
	ModelUsage       map[string]ModelTokens `json:"modelUsage"`
	TotalSessions    int                `json:"totalSessions"`
	TotalMessages    int                `json:"totalMessages"`
	FirstSessionDate string             `json:"firstSessionDate"`
	HourCounts       map[string]int     `json:"hourCounts"`
}

func LoadStatsCache(claudeHome string) (StatsCache, error) {
	path := filepath.Join(claudeHome, "stats-cache.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return StatsCache{}, err
	}

	var f statsCacheFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return StatsCache{}, err
	}

	return StatsCache{
		DailyActivity:    f.DailyActivity,
		DailyModelTokens: f.DailyModelTokens,
		ModelUsage:       f.ModelUsage,
		TotalSessions:    f.TotalSessions,
		TotalMessages:    f.TotalMessages,
		FirstSessionDate: f.FirstSessionDate,
		HourCounts:       f.HourCounts,
	}, nil
}
