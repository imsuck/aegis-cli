package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"aegis-cli/internal/vault"
)

// Mode represents the current TUI mode
type Mode int

const (
	ModePassword Mode = iota
	ModeTable
	ModeSearch
)

// Model represents the TUI application state
type Model struct {
	// State
	mode            Mode
	entries         []vault.Entry
	content         *vault.Content
	err             error
	vaultPath       string
	lastCopyTime    time.Time
	copySuccess     bool
	showCodes       bool

	// Password input
	passwordInput textinput.Model

	// Search input
	searchInput textinput.Model

	// Table state
	cursor          int
	filteredEntries []vault.Entry

	// Timing
	lastUpdate     time.Time
	lastActivity   time.Time
	timeout        time.Duration
	warningShown   bool
}

// NewModel creates a new TUI model
func NewModel(vaultPath string, timeout time.Duration) Model {
	// Setup password input
	ti := textinput.New()
	ti.Placeholder = "Enter vault password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'
	ti.Focus()

	// Setup search input
	si := textinput.New()
	si.Placeholder = "Search entries..."

	return Model{
		mode:            ModePassword,
		vaultPath:       vaultPath,
		passwordInput:   ti,
		searchInput:     si,
		lastActivity:    time.Now(),
		timeout:         timeout,
		filteredEntries: []vault.Entry{},
	}
}
