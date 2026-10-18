package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"
)

type daemonLock struct {
	PID       int    `json:"pid"`
	Version   string `json:"version"`
	StartedAt int64  `json:"startedAt"`
}

type rosterFile struct {
	Workers map[string]json.RawMessage `json:"workers"`
}

func LoadDaemonInfo(claudeHome string) DaemonInfo {
	var info DaemonInfo

	raw, err := os.ReadFile(filepath.Join(claudeHome, "daemon.lock"))
	if err != nil {
		return info
	}
	var lock daemonLock
	if json.Unmarshal(raw, &lock) != nil {
		return info
	}

	info.PID = lock.PID
	info.Version = lock.Version
	info.StartedAt = lock.StartedAt
	info.Running = isProcessAlive(lock.PID)

	raw, err = os.ReadFile(filepath.Join(claudeHome, "daemon", "roster.json"))
	if err == nil {
		var roster rosterFile
		if json.Unmarshal(raw, &roster) == nil {
			info.Workers = len(roster.Workers)
		}
	}

	return info
}

func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}
