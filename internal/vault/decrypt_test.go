package vault_test

import (
	"testing"
	"aegis-cli/internal/vault"
)

const testVaultPath = "../../test/resources/aegis_encrypted.json"

func TestDecryptVault(t *testing.T) {
	// Password for test fixture is "test"
	result, err := vault.LoadAndDecrypt(testVaultPath, "test")
	if err != nil {
		t.Fatalf("failed to decrypt vault: %v", err)
	}

	if len(result.Content.Entries) == 0 {
		t.Error("expected at least one entry in decrypted vault")
	}
}

func TestDecryptVaultWrongPassword(t *testing.T) {
	_, err := vault.LoadAndDecrypt(testVaultPath, "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}
