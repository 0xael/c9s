package control

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

func run(ctx context.Context, dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "claude", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("claude %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func LaunchBackground(ctx context.Context, prompt string, cwd string) (string, error) {
	return run(ctx, cwd, "--bg", prompt)
}

func StopSession(ctx context.Context, id string) error {
	_, err := run(ctx, "", "stop", id)
	return err
}

func RemoveJob(ctx context.Context, id string) error {
	_, err := run(ctx, "", "rm", id)
	return err
}

func RespawnJob(ctx context.Context, id string) error {
	_, err := run(ctx, "", "respawn", id)
	return err
}

func AttachCommand(id string) string {
	return fmt.Sprintf("claude attach %s", id)
}
