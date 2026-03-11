package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"aegis-cli/internal/tui"
)

func main() {
	timeoutFlag := flag.Duration("timeout", 60*time.Second, "Auto-exit after duration of inactivity")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: aegis-cli <vault-file> [-timeout duration]")
		os.Exit(1)
	}

	vaultPath := args[0]
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: vault file not found: %s\n", vaultPath)
		os.Exit(1)
	}

	model := tui.NewModel(vaultPath, *timeoutFlag)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
