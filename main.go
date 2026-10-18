package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/0xael/c9s/internal/app"
	"github.com/0xael/c9s/internal/config"
	"github.com/0xael/c9s/internal/web"
)

func main() {
	cfg := config.Load()

	if len(os.Args) > 1 && os.Args[1] == "web" {
		port := 9731
		noOpen := false
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--port":
				if i+1 < len(os.Args) {
					port, _ = strconv.Atoi(os.Args[i+1])
					i++
				}
			case "--no-open":
				noOpen = true
			}
		}
		web.Run(cfg, port, noOpen)
		return
	}

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
