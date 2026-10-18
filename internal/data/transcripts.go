package data

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type usageEnvelope struct {
	Type    string `json:"type"`
	Message struct {
		Usage struct {
			InputTokens              int64 `json:"input_tokens"`
			OutputTokens             int64 `json:"output_tokens"`
			CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
			OutputTokensDetails      struct {
				ThinkingTokens int64 `json:"thinking_tokens"`
			} `json:"output_tokens_details"`
		} `json:"usage"`
	} `json:"message"`
}

func LoadSessionTokens(claudeHome, sessionID string) TokenUsage {
	path, err := FindTranscriptPath(claudeHome, sessionID)
	if err != nil {
		return TokenUsage{}
	}
	return sumTokens(path)
}

func sumTokens(path string) TokenUsage {
	f, err := os.Open(path)
	if err != nil {
		return TokenUsage{}
	}
	defer f.Close()

	var total TokenUsage
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) < 30 {
			continue
		}
		var e usageEnvelope
		if json.Unmarshal(line, &e) != nil || e.Type != "assistant" {
			continue
		}
		total.InputTokens += e.Message.Usage.InputTokens
		total.OutputTokens += e.Message.Usage.OutputTokens
		total.CacheCreationInputTokens += e.Message.Usage.CacheCreationInputTokens
		total.CacheReadInputTokens += e.Message.Usage.CacheReadInputTokens
		total.ThinkingTokens += e.Message.Usage.OutputTokensDetails.ThinkingTokens
	}
	return total
}

func LoadTranscript(path string, maxLines int) ([]TranscriptEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var all []TranscriptEntry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var e TranscriptEntry
		if json.Unmarshal(line, &e) == nil && e.Type != "" {
			all = append(all, e)
		}
	}
	if err := sc.Err(); err != nil {
		return all, err
	}

	if maxLines > 0 && len(all) > maxLines {
		all = all[len(all)-maxLines:]
	}
	return all, nil
}

func FindTranscriptPath(claudeHome, sessionID string) (string, error) {
	projectsDir := filepath.Join(claudeHome, "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return "", err
	}

	suffix := sessionID + ".jsonl"
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		projDir := filepath.Join(projectsDir, e.Name())
		files, err := os.ReadDir(projDir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
				return filepath.Join(projDir, f.Name()), nil
			}
		}
	}
	return "", os.ErrNotExist
}
