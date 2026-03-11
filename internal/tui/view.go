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
			Background(lipgloss.Color("4"))

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
	b.WriteString(strings.Repeat("─", 78))
	b.WriteString("\n")

	// Entries
	for i, entry := range m.filteredEntries {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render(m.formatSelectedRow(entry, i)))
		} else {
			b.WriteString(m.formatEntryRow(entry, i))
		}
		b.WriteString("\n")
	}

	if len(m.filteredEntries) == 0 {
		b.WriteString(dimStyle.Render("No entries found"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.showCodes {
		b.WriteString(dimStyle.Render("j/k: navigate • y: copy • /: search • c: hide codes • q: quit"))
	} else {
		b.WriteString(dimStyle.Render("j/k: navigate • y: copy • /: search • c: show codes • q: quit"))
	}

	return b.String()
}

func (m Model) formatEntryRow(entry vault.Entry, index int) string {
	code := m.getCodeForEntry(index)
	remaining := m.getRemainingTime(index)

	note := entry.Note
	if len(note) > 18 {
		note = note[:15] + "..."
	}

	return fmt.Sprintf("%-20s %-20s %-20s %-20s %-6s",
		truncate(entry.Issuer, 20),
		truncate(entry.Name, 20),
		truncate(note, 20),
		codeStyle.Render(code),
		formatTimer(remaining),
	)
}

func (m Model) formatSelectedRow(entry vault.Entry, index int) string {
	code := m.getCodeForEntry(index)
	remaining := m.getRemainingTime(index)

	note := entry.Note
	if len(note) > 18 {
		note = note[:15] + "..."
	}

	// For selected row, use black text for code as well
	selectedCodeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0"))

	return fmt.Sprintf("%-20s %-20s %-20s %-20s %-6s",
		truncate(entry.Issuer, 20),
		truncate(entry.Name, 20),
		truncate(note, 20),
		selectedCodeStyle.Render(code),
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

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("y: copy • Enter: accept • Esc: cancel"))

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
	// Format with fixed width (5 chars: "  25s") then apply style
	timerStr := fmt.Sprintf("%4ds", remaining)
	return style.Render(timerStr)
}

// getCodeForEntry generates the current TOTP code for an entry
// When showCodes is false, all codes are masked
// When showCodes is true, only the selected entry shows its code
func (m Model) getCodeForEntry(index int) string {
	if index < 0 || index >= len(m.filteredEntries) {
		return "------"
	}
	
	entry := m.filteredEntries[index]
	digits := entry.Info.Digits
	if digits == 0 {
		digits = 6
	}
	
	// Mask codes when showCodes is disabled, or when entry is not selected
	if !m.showCodes || index != m.cursor {
		return strings.Repeat("*", digits)
	}
	
	code, err := totp.Generate(entry, m.lastUpdate)
	if err != nil {
		return "ERROR"
	}
	return code
}

// getActualCodeForEntry returns the actual TOTP code for clipboard copy
// This always returns the real code, regardless of showCodes setting
func (m Model) getActualCodeForEntry(index int) string {
	if index < 0 || index >= len(m.filteredEntries) {
		return ""
	}
	
	entry := m.filteredEntries[index]
	code, err := totp.Generate(entry, m.lastUpdate)
	if err != nil {
		return ""
	}
	return code
}

func (m Model) getRemainingTime(index int) int {
	if index < 0 || index >= len(m.filteredEntries) {
		return 0
	}
	return totp.GetRemainingTime(m.filteredEntries[index], m.lastUpdate)
}
