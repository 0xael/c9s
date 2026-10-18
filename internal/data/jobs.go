package data

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadJobs(claudeHome string) ([]Job, error) {
	jobsDir := filepath.Join(claudeHome, "jobs")
	entries, err := os.ReadDir(jobsDir)
	if err != nil {
		return nil, nil
	}

	pins := loadPins(claudeHome)

	var jobs []Job
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		id := e.Name()
		statePath := filepath.Join(jobsDir, id, "state.json")
		raw, err := os.ReadFile(statePath)
		if err != nil {
			continue
		}
		var j Job
		if err := json.Unmarshal(raw, &j); err != nil {
			continue
		}
		j.ID = id
		j.Pinned = pins[id]

		tl, _ := LoadTimeline(filepath.Join(jobsDir, id, "timeline.jsonl"))
		j.Timeline = tl

		jobs = append(jobs, j)
	}
	return jobs, nil
}

func LoadTimeline(path string) ([]TimelineEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []TimelineEntry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 256*1024), 256*1024)
	for sc.Scan() {
		var e TimelineEntry
		if json.Unmarshal(sc.Bytes(), &e) == nil {
			entries = append(entries, e)
		}
	}
	return entries, sc.Err()
}

func loadPins(claudeHome string) map[string]bool {
	raw, err := os.ReadFile(filepath.Join(claudeHome, "jobs", "pins.json"))
	if err != nil {
		return nil
	}
	var ids []string
	if json.Unmarshal(raw, &ids) != nil {
		return nil
	}
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}
