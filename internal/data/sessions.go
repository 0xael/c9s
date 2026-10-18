package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type sessionFile struct {
	MessagingSocketPath string            `json:"messagingSocketPath"`
	Version             string            `json:"version"`
	Kind                string            `json:"kind"`
	Status              string            `json:"status"`
	ProviderEnv         map[string]string `json:"providerEnv"`
}

func EnrichSession(claudeHome string, s *Session) error {
	pattern := filepath.Join(claudeHome, "sessions", fmt.Sprintf("%d.json", s.PID))
	raw, err := os.ReadFile(pattern)
	if err != nil {
		return err
	}

	var sf sessionFile
	if err := json.Unmarshal(raw, &sf); err != nil {
		return err
	}

	s.SocketPath = sf.MessagingSocketPath
	if sf.Version != "" {
		s.Version = sf.Version
	}
	if model, ok := sf.ProviderEnv["ANTHROPIC_MODEL"]; ok {
		s.Model = model
	}
	return nil
}
