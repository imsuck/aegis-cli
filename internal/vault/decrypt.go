package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/scrypt"
)

// DecryptResult holds the decrypted vault with raw and parsed content
type DecryptResult struct {
	Raw     Vault
	Content Content
}

// LoadAndDecrypt loads an encrypted vault file and decrypts it with the given password
func LoadAndDecrypt(path, password string) (*DecryptResult, error) {
	// Read the vault file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	// Parse JSON
	var vault Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, fmt.Errorf("failed to parse vault JSON: %w", err)
	}

	// Find password slot
	var passwordSlot *Slot
	for i := range vault.Header.Slots {
		if vault.Header.Slots[i].Type == 1 {
			passwordSlot = &vault.Header.Slots[i]
			break
		}
	}
	if passwordSlot == nil {
		return nil, fmt.Errorf("no password slot found in vault")
	}

	// Check if vault is encrypted
	if vault.Header.Params == nil {
		return nil, fmt.Errorf("vault is not encrypted")
	}

	// Derive key from password using scrypt
	salt, err := hex.DecodeString(passwordSlot.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	derivedKey, err := scrypt.Key(
		[]byte(password),
		salt,
		passwordSlot.N,
		passwordSlot.R,
		passwordSlot.P,
		32, // 256-bit key
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Decrypt master key
	masterKey, err := decryptKey(derivedKey, passwordSlot)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt master key: %w", err)
	}

	// Decrypt vault contents
	content, err := decryptContent(masterKey, &vault)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault content: %w", err)
	}

	return &DecryptResult{
		Raw:     vault,
		Content: content,
	}, nil
}

// decryptKey decrypts the master key from a slot
func decryptKey(derivedKey []byte, slot *Slot) ([]byte, error) {
	keyBytes, err := hex.DecodeString(slot.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	nonce, err := hex.DecodeString(slot.KeyParams.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	tag, err := hex.DecodeString(slot.KeyParams.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tag: %w", err)
	}

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Append tag to ciphertext for Go's Open
	ciphertext := append(keyBytes, tag...)
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// decryptContent decrypts the vault database
func decryptContent(masterKey []byte, vault *Vault) (Content, error) {
	// Decode base64 encoded content
	// json.RawMessage includes quotes, use json.Unmarshal to properly decode
	var dbStr string
	if err := json.Unmarshal(vault.DB, &dbStr); err != nil {
		return Content{}, fmt.Errorf("failed to unmarshal DB string: %w", err)
	}
	contentBytes, err := base64.StdEncoding.DecodeString(dbStr)
	if err != nil {
		return Content{}, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	nonce, err := hex.DecodeString(vault.Header.Params.Nonce)
	if err != nil {
		return Content{}, fmt.Errorf("failed to decode nonce: %w", err)
	}

	tag, err := hex.DecodeString(vault.Header.Params.Tag)
	if err != nil {
		return Content{}, fmt.Errorf("failed to decode tag: %w", err)
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return Content{}, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Content{}, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Append tag to ciphertext
	ciphertext := append(contentBytes, tag...)
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return Content{}, fmt.Errorf("failed to decrypt content: %w", err)
	}

	// Parse JSON content
	var content Content
	if err := json.Unmarshal(plaintext, &content); err != nil {
		return Content{}, fmt.Errorf("failed to parse content JSON: %w", err)
	}

	return content, nil
}
