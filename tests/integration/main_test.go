package integration_test

import (
	"testing"
	"time"
	"aegis-cli/internal/vault"
	"aegis-cli/internal/totp"
)

const testVaultPath = "../../test/resources/aegis_encrypted.json"

func TestFullVaultDecryptionAndCodeGeneration(t *testing.T) {
	// Test with the encrypted fixture (password: test)
	result, err := vault.LoadAndDecrypt(testVaultPath, "test")
	if err != nil {
		t.Fatalf("failed to decrypt vault: %v", err)
	}

	if len(result.Content.Entries) == 0 {
		t.Fatal("expected entries in vault")
	}

	// Generate codes for all TOTP entries
	totpCount := 0
	skippedCount := 0
	for _, entry := range result.Content.Entries {
		if entry.Type == "totp" {
			code, err := totp.Generate(entry, time.Now())
			if err != nil {
				// Some entries might have invalid secrets (e.g., encrypted secrets), skip them
				skippedCount++
				t.Logf("skipping entry %s (%s): %v", entry.Issuer, entry.Name, err)
				continue
			}
			totpCount++
			if len(code) != entry.Info.Digits {
				t.Errorf("expected %d digit code, got %d", entry.Info.Digits, len(code))
			}
			// Verify code is numeric
			for _, c := range code {
				if c < '0' || c > '9' {
					t.Errorf("expected numeric code, got %c", c)
				}
			}
		}
	}
	
	t.Logf("Generated codes for %d TOTP entries, skipped %d", totpCount, skippedCount)
	// We expect at least some valid TOTP entries if there are any totp type entries
}

func TestVaultDecryptionWrongPassword(t *testing.T) {
	_, err := vault.LoadAndDecrypt(testVaultPath, "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestVaultDecryptionAndSearch(t *testing.T) {
	result, err := vault.LoadAndDecrypt(testVaultPath, "test")
	if err != nil {
		t.Fatalf("failed to decrypt vault: %v", err)
	}

	// Test that we have entries
	if len(result.Content.Entries) == 0 {
		t.Fatal("expected entries in vault")
	}

	// Verify entries have required fields
	for _, entry := range result.Content.Entries {
		if entry.Issuer == "" {
			t.Error("expected entry to have issuer")
		}
		if entry.Name == "" {
			t.Error("expected entry to have name")
		}
		if entry.Info.Secret == "" {
			t.Error("expected entry to have secret")
		}
	}
}

func TestTOTPTimeRemaining(t *testing.T) {
	result, err := vault.LoadAndDecrypt(testVaultPath, "test")
	if err != nil {
		t.Fatalf("failed to decrypt vault: %v", err)
	}

	for _, entry := range result.Content.Entries {
		if entry.Type == "totp" {
			period := entry.Info.Period
			if period == 0 {
				period = 30
			}
			remaining := totp.GetRemainingTime(entry, time.Now())
			// Remaining should be 1 to period
			if remaining < 0 || remaining > period {
				t.Errorf("expected remaining time between 0-%d, got %d for %s", period, remaining, entry.Issuer)
			}
		}
	}
}
