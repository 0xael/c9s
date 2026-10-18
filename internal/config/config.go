package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Claude  ClaudeConfig
	Refresh RefreshConfig
	UI      UIConfig
}

type ClaudeConfig struct {
	Home string
	CLI  string
}

type RefreshConfig struct {
	Fast time.Duration
	Slow time.Duration
}

type UIConfig struct {
	Theme                  string
	TimestampFormat        string `mapstructure:"timestamp_format"`
	TranscriptInitialLines int    `mapstructure:"transcript_initial_lines"`
	ShowPID                bool   `mapstructure:"show_pid"`
}

func Load() Config {
	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(xdgConfigHome(), "c9s"))
	viper.AddConfigPath("$HOME/.c9s")

	viper.SetEnvPrefix("C9S")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	var cfg Config
	cfg.Claude.Home = viper.GetString("claude.home")
	cfg.Claude.CLI = viper.GetString("claude.cli")
	cfg.Refresh.Fast = viper.GetDuration("refresh.fast")
	cfg.Refresh.Slow = viper.GetDuration("refresh.slow")
	cfg.UI.Theme = viper.GetString("ui.theme")
	cfg.UI.TimestampFormat = viper.GetString("ui.timestamp_format")
	cfg.UI.TranscriptInitialLines = viper.GetInt("ui.transcript_initial_lines")
	cfg.UI.ShowPID = viper.GetBool("ui.show_pid")

	if cfg.Claude.Home == "~/.claude" {
		home, err := os.UserHomeDir()
		if err == nil {
			cfg.Claude.Home = filepath.Join(home, ".claude")
		}
	}

	return cfg
}

func xdgConfigHome() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}
