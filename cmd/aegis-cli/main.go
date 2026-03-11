package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"aegis-cli/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: aegis-cli <vault-file>")
		os.Exit(1)
	}

	vaultPath := os.Args[1]
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: vault file not found: %s\n", vaultPath)
		os.Exit(1)
	}

	model := tui.NewModel(vaultPath)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
