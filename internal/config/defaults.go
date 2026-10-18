package config

import "github.com/spf13/viper"

func setDefaults() {
	viper.SetDefault("claude.home", "~/.claude")
	viper.SetDefault("claude.cli", "claude")

	viper.SetDefault("refresh.fast", "2s")
	viper.SetDefault("refresh.slow", "5s")

	viper.SetDefault("ui.theme", "default")
	viper.SetDefault("ui.timestamp_format", "relative")
	viper.SetDefault("ui.transcript_initial_lines", 200)
	viper.SetDefault("ui.show_pid", true)
}
