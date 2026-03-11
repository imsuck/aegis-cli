package search

import (
	"strings"
	"aegis-cli/internal/vault"
)

// SearchEntries filters entries based on query matching issuer, name, and note
// The search is case-insensitive and matches concatenated "issuer name note" fields
func SearchEntries(entries []vault.Entry, query string) []vault.Entry {
	if query == "" {
		return entries
	}

	query = strings.ToLower(query)
	var results []vault.Entry

	for _, entry := range entries {
		// Concatenate fields for matching: "issuer name note"
		searchable := strings.ToLower(entry.Issuer + " " + entry.Name + " " + entry.Note)

		// Simple substring match
		if strings.Contains(searchable, query) {
			results = append(results, entry)
		}
	}

	return results
}
