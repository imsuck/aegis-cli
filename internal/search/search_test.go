package search_test

import (
	"testing"
	"aegis-cli/internal/search"
	"aegis-cli/internal/vault"
)

func TestSearchEntries(t *testing.T) {
	entries := []vault.Entry{
		{Issuer: "Google", Name: "Alice", Note: "Work"},
		{Issuer: "GitHub", Name: "Bob", Note: "Personal"},
		{Issuer: "Amazon", Name: "Charlie", Note: "Shopping"},
	}

	// Test issuer match
	results := search.SearchEntries(entries, "google")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'google', got %d", len(results))
	}

	// Test name match
	results = search.SearchEntries(entries, "bob")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'bob', got %d", len(results))
	}

	// Test note match
	results = search.SearchEntries(entries, "shopping")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'shopping', got %d", len(results))
	}

	// Test partial match
	results = search.SearchEntries(entries, "goo")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'goo', got %d", len(results))
	}

	// Test no match
	results = search.SearchEntries(entries, "xyz")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'xyz', got %d", len(results))
	}
}

func TestSearchConcatenatedFields(t *testing.T) {
	entries := []vault.Entry{
		{Issuer: "Google", Name: "Alice", Note: "Work"},
	}

	// Search should match concatenated "issuer name note"
	results := search.SearchEntries(entries, "google alice")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'google alice', got %d", len(results))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	entries := []vault.Entry{
		{Issuer: "Google", Name: "Alice", Note: "Work"},
	}

	// Search should be case-insensitive
	results := search.SearchEntries(entries, "GOOGLE")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'GOOGLE', got %d", len(results))
	}

	results = search.SearchEntries(entries, "alice")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'alice', got %d", len(results))
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	entries := []vault.Entry{
		{Issuer: "Google", Name: "Alice", Note: "Work"},
	}

	// Empty query should return all entries
	results := search.SearchEntries(entries, "")
	if len(results) != 1 {
		t.Errorf("expected 1 result for empty query, got %d", len(results))
	}
}
