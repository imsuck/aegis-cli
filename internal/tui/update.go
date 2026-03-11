package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"aegis-cli/internal/vault"
)

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	return tick()
}

// TickMsg represents a timer tick
type TickMsg time.Time

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case TickMsg:
		return m.handleTick(msg)
	case passwordSubmittedMsg:
		return m.handlePasswordSubmit(string(msg))
	case vaultLoadedMsg:
		return m.handleVaultLoaded(msg)
	case vaultErrorMsg:
		return m.handleVaultError(msg)
	}

	return m, nil
}

func (m Model) handleKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModePassword:
		return m.handlePasswordKey(key)
	case ModeTable:
		return m.handleTableKey(key)
	case ModeSearch:
		return m.handleSearchKey(key)
	case ModeCodeDisplay:
		return m.handleCodeDisplayKey(key)
	}
	return m, nil
}

func (m Model) handlePasswordKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.Type {
	case tea.KeyEnter:
		return m, submitPassword(m.passwordInput.Value())
	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.passwordInput, cmd = m.passwordInput.Update(key)
	return m, cmd
}

func (m Model) handleTableKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		if m.cursor < len(m.filteredEntries)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "g":
		m.cursor = 0
	case "G":
		m.cursor = len(m.filteredEntries) - 1
	case "/":
		m.mode = ModeSearch
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		return m, textinput.Blink
	case "c":
		m.mode = ModeCodeDisplay
		return m, nil
	case "y":
		if len(m.filteredEntries) > 0 {
			return m, copyToClipboard(m.getCodeForEntry(m.cursor))
		}
	}
	return m, nil
}

func (m Model) handleSearchKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.Type {
	case tea.KeyEsc:
		m.mode = ModeTable
		m.searchInput.SetValue("")
		m.filteredEntries = m.entries
		m.cursor = 0
		return m, nil
	case tea.KeyEnter:
		m.mode = ModeTable
		return m, nil
	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	// Update search input
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(key)

	// Filter entries based on query
	query := m.searchInput.Value()
	m.filteredEntries = filterEntries(m.entries, query)
	m.cursor = 0

	return m, cmd
}

func (m Model) handleCodeDisplayKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "c":
		m.mode = ModeTable
		return m, nil
	case "j", "down":
		if m.cursor < len(m.filteredEntries)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "y":
		if len(m.filteredEntries) > 0 {
			return m, copyToClipboard(m.getCodeForEntry(m.cursor))
		}
	}
	return m, nil
}

func (m Model) handleTick(t TickMsg) (tea.Model, tea.Cmd) {
	m.lastUpdate = time.Time(t)
	return m, tick()
}

func (m Model) handlePasswordSubmit(password string) (tea.Model, tea.Cmd) {
	// Start loading the vault
	return m, loadVaultAsync(m.vaultPath, password)
}

func (m Model) handleVaultLoaded(msg vaultLoadedMsg) (tea.Model, tea.Cmd) {
	m.content = msg.content
	m.entries = msg.content.Entries
	m.filteredEntries = msg.content.Entries
	m.mode = ModeTable
	m.cursor = 0
	return m, tick()
}

func (m Model) handleVaultError(msg vaultErrorMsg) (tea.Model, tea.Cmd) {
	m.err = msg.err
	m.passwordInput.SetValue("")
	m.passwordInput.Placeholder = "Wrong password, try again"
	return m, nil
}

// filterEntries filters entries based on search query
func filterEntries(entries []vault.Entry, query string) []vault.Entry {
	if query == "" {
		return entries
	}

	result := []vault.Entry{}
	for _, entry := range entries {
		searchable := entry.Issuer + " " + entry.Name + " " + entry.Note
		if containsIgnoreCase(searchable, query) {
			result = append(result, entry)
		}
	}
	return result
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		(len(substr) == 0 || 
			findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// tick returns a command that sends a TickMsg every second
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// submitPassword returns a command to submit the password
func submitPassword(password string) tea.Cmd {
	return func() tea.Msg {
		return passwordSubmittedMsg(password)
	}
}

type passwordSubmittedMsg string

// loadVaultAsync returns a command to load and decrypt the vault asynchronously
func loadVaultAsync(path, password string) tea.Cmd {
	return func() tea.Msg {
		result, err := vault.LoadAndDecrypt(path, password)
		if err != nil {
			return vaultErrorMsg{err: err}
		}
		return vaultLoadedMsg{content: &result.Content}
	}
}

type vaultLoadedMsg struct {
	content *vault.Content
}

type vaultErrorMsg struct {
	err error
}

// copyToClipboard returns a command to copy text to clipboard
func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		return copyMsg{text: text}
	}
}

type copyMsg struct {
	text string
}
