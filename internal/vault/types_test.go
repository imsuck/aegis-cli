package vault_test

import (
	"encoding/json"
	"testing"
	"aegis-cli/internal/vault"
)

func TestParseVaultHeader(t *testing.T) {
	jsonStr := `{
		"slots": [{
			"type": 1,
			"uuid": "a8325752-c1be-458a-9b3e-5e0a8154d9ec",
			"key": "491d44550430ba248986b904b8cffd3a",
			"key_params": {
				"nonce": "e9705513ba4951fa7a0608d2",
				"tag": "931237af257b83c693ddb8f9a7eddaf0"
			},
			"n": 32768,
			"r": 8,
			"p": 1,
			"salt": "27ea9ae53fa2f08a8dcd201615a8229422647b3058f9f36b08f9457e62888be1"
		}],
		"params": {
			"nonce": "095fd13dee336fa56b4634ff",
			"tag": "5db2470edf2d12f82a89ae7f48ccd50c"
		}
	}`

	var header vault.Header
	err := json.Unmarshal([]byte(jsonStr), &header)
	if err != nil {
		t.Fatalf("failed to unmarshal header: %v", err)
	}

	if len(header.Slots) != 1 {
		t.Errorf("expected 1 slot, got %d", len(header.Slots))
	}
	if header.Slots[0].Type != 1 {
		t.Errorf("expected slot type 1, got %d", header.Slots[0].Type)
	}
	if header.Params == nil {
		t.Errorf("expected params to be non-nil")
	}
}

func TestParseVaultContent(t *testing.T) {
	jsonStr := `{
		"version": 3,
		"entries": [{
			"type": "totp",
			"uuid": "01234567-89ab-cdef-0123-456789abcdef",
			"name": "Alice",
			"issuer": "Google",
			"note": "Work account",
			"favorite": false,
			"icon": null,
			"icon_mime": null,
			"icon_hash": null,
			"info": {
				"secret": "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567",
				"algo": "SHA1",
				"digits": 6,
				"period": 30
			},
			"groups": []
		}],
		"groups": []
	}`

	var content vault.Content
	err := json.Unmarshal([]byte(jsonStr), &content)
	if err != nil {
		t.Fatalf("failed to unmarshal content: %v", err)
	}

	if len(content.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(content.Entries))
	}
	if content.Entries[0].Issuer != "Google" {
		t.Errorf("expected issuer Google, got %s", content.Entries[0].Issuer)
	}
}

func TestParseUnencryptedVault(t *testing.T) {
	jsonStr := `{
		"slots": [{
			"type": 1,
			"uuid": "a8325752-c1be-458a-9b3e-5e0a8154d9ec",
			"key": "491d44550430ba248986b904b8cffd3a",
			"key_params": {
				"nonce": "e9705513ba4951fa7a0608d2",
				"tag": "931237af257b83c693ddb8f9a7eddaf0"
			},
			"n": 32768,
			"r": 8,
			"p": 1,
			"salt": "27ea9ae53fa2f08a8dcd201615a8229422647b3058f9f36b08f9457e62888be1"
		}],
		"params": null
	}`

	var header vault.Header
	err := json.Unmarshal([]byte(jsonStr), &header)
	if err != nil {
		t.Fatalf("failed to unmarshal header: %v", err)
	}

	if header.Params != nil {
		t.Errorf("expected params to be nil for unencrypted vault")
	}
}

func TestParseSlotWithRepairedField(t *testing.T) {
	jsonStr := `{
		"type": 1,
		"uuid": "a8325752-c1be-458a-9b3e-5e0a8154d9ec",
		"key": "491d44550430ba248986b904b8cffd3a",
		"key_params": {
			"nonce": "e9705513ba4951fa7a0608d2",
			"tag": "931237af257b83c693ddb8f9a7eddaf0"
		},
		"n": 32768,
		"r": 8,
		"p": 1,
		"salt": "27ea9ae53fa2f08a8dcd201615a8229422647b3058f9f36b08f9457e62888be1",
		"repaired": true
	}`

	var slot vault.Slot
	err := json.Unmarshal([]byte(jsonStr), &slot)
	if err != nil {
		t.Fatalf("failed to unmarshal slot: %v", err)
	}

	if !slot.Repaired {
		t.Errorf("expected repaired to be true")
	}
}
