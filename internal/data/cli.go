package data

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

func FetchAgents(claudeHome string) ([]Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "claude", "agents", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sessions []Session
	if err := json.Unmarshal(out, &sessions); err != nil {
		return nil, err
	}

	for i := range sessions {
		_ = EnrichSession(claudeHome, &sessions[i])
	}

	return sessions, nil
}
