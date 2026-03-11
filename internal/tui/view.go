package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"aegis-cli/internal/vault"
	"aegis-cli/internal/totp"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("13")).
			Padding(0, 1)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	codeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10"))

	timerGoodStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	timerWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))

	timerBadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("1"))
)

func (m Model) View() string {
	switch m.mode {
	case ModePassword:
		return m.passwordView()
	case ModeTable:
		return m.tableView()
	case ModeSearch:
		return m.searchView()
	case ModeCodeDisplay:
		return m.codeDisplayView()
	default:
		return "Unknown mode"
	}
}

func (m Model) passwordView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Aegis Vault Unlock"))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v\n\n", m.err)))
	}

	b.WriteString("Enter your vault password:\n\n")
	b.WriteString(m.passwordInput.View())
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("Press Enter to unlock, Ctrl+C to quit"))

	return b.String()
}

func (m Model) tableView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Aegis Vault"))
	b.WriteString("\n\n")

	// Table header
	header := fmt.Sprintf("%-20s %-20s %-20s %-10s %-6s",
		"ISSUER", "NAME", "NOTE", "CODE", "TIME")
	b.WriteString(dimStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 70))
	b.WriteString("\n")

	// Entries
	for i, entry := range m.filteredEntries {
		row := m.formatEntryRow(entry, i)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render(row))
		} else {
			b.WriteString(row)
		}
		b.WriteString("\n")
	}

	if len(m.filteredEntries) == 0 {
		b.WriteString(dimStyle.Render("No entries found"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("j/k: navigate • y: copy code • /: search • c: code view • q: quit"))

	return b.String()
}

func (m Model) formatEntryRow(entry vault.Entry, index int) string {
	code := m.getCodeForEntry(index)
	remaining := m.getRemainingTime(index)

	note := entry.Note
	if len(note) > 18 {
		note = note[:15] + "..."
	}

	return fmt.Sprintf("%-20s %-20s %-20s %-10s %s",
		truncate(entry.Issuer, 20),
		truncate(entry.Name, 20),
		truncate(note, 20),
		codeStyle.Render(code),
		formatTimer(remaining),
	)
}

func (m Model) searchView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Search"))
	b.WriteString("\n\n")
	b.WriteString("Type to filter entries (Esc to cancel):\n\n")
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	// Show filtered results
	for i, entry := range m.filteredEntries {
		row := m.formatEntryRow(entry, i)
		b.WriteString(row)
		b.WriteString("\n")
	}

	if len(m.filteredEntries) == 0 {
		b.WriteString(dimStyle.Render("No matches found"))
	}

	return b.String()
}

func (m Model) codeDisplayView() string {
	var b strings.Builder

	if len(m.filteredEntries) == 0 {
		return "No entries"
	}

	entry := m.filteredEntries[m.cursor]
	code := m.getCodeForEntry(m.cursor)
	remaining := m.getRemainingTime(m.cursor)

	b.WriteString(titleStyle.Render("Code Display"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Issuer: %s\n", entry.Issuer))
	b.WriteString(fmt.Sprintf("Name:   %s\n", entry.Name))
	if entry.Note != "" {
		b.WriteString(fmt.Sprintf("Note:   %s\n", entry.Note))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Code: %s\n", codeStyle.Render(code)))
	b.WriteString(fmt.Sprintf("Refreshes in: %s\n", formatTimer(remaining)))
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("y: copy • j/k: navigate • esc/c: back to table • q: quit"))

	return b.String()
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func formatTimer(remaining int) string {
	var style lipgloss.Style
	switch {
	case remaining > 15:
		style = timerGoodStyle
	case remaining > 5:
		style = timerWarnStyle
	default:
		style = timerBadStyle
	}
	return style.Render(fmt.Sprintf("%2ds", remaining))
}

// getCodeForEntry generates the current TOTP code for an entry
func (m Model) getCodeForEntry(index int) string {
	if index < 0 || index >= len(m.filteredEntries) {
		return "------"
	}
	code, err := totp.Generate(m.filteredEntries[index], m.lastUpdate)
	if err != nil {
		return "ERROR"
	}
	return code
}

func (m Model) getRemainingTime(index int) int {
	if index < 0 || index >= len(m.filteredEntries) {
		return 0
	}
	return totp.GetRemainingTime(m.filteredEntries[index], m.lastUpdate)
}
