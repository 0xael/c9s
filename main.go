package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elalwany/c9s/internal/app"
	"github.com/elalwany/c9s/internal/config"
)

func main() {
	cfg := config.Load()
	m := app.New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if model, ok := result.(app.Model); ok && model.AttachCmd() != "" {
		fmt.Printf("\nRun: %s\n", model.AttachCmd())
	}
}
