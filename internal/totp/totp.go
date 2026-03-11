package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"hash"
	"math"
	"time"

	"aegis-cli/internal/vault"
)

// Generate generates a TOTP code for the given entry at the specified time
func Generate(entry vault.Entry, t time.Time) (string, error) {
	// Decode base32 secret
	secret, err := base32.StdEncoding.DecodeString(entry.Info.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}

	// Get hash function
	hasher := getHasher(entry.Info.Algo)
	if hasher == nil {
		return "", fmt.Errorf("unsupported algorithm: %s", entry.Info.Algo)
	}

	// Calculate time counter
	period := entry.Info.Period
	if period == 0 {
		period = 30
	}
	counter := uint64(t.Unix()) / uint64(period)

	// Generate HOTP code
	return generateHOTP(secret, counter, hasher, entry.Info.Digits)
}

// GetRemainingTime returns the seconds remaining until the next code
func GetRemainingTime(entry vault.Entry, t time.Time) int {
	period := entry.Info.Period
	if period == 0 {
		period = 30
	}
	return period - int(t.Unix()%int64(period))
}

// getHasher returns the appropriate hash function based on algorithm name
func getHasher(algo string) func() hash.Hash {
	switch algo {
	case "SHA1":
		return sha1.New
	case "SHA256":
		return sha256.New
	case "SHA512":
		return sha512.New
	default:
		return nil
	}
}

// generateHOTP generates an HOTP code
func generateHOTP(secret []byte, counter uint64, hasher func() hash.Hash, digits int) (string, error) {
	// Create HMAC
	h := hmac.New(hasher, secret)

	// Write counter as big-endian 8-byte array
	var counterBytes [8]byte
	for i := 7; i >= 0; i-- {
		counterBytes[i] = byte(counter & 0xff)
		counter >>= 8
	}
	h.Write(counterBytes[:])

	// Get HMAC result
	sum := h.Sum(nil)

	// Dynamic truncation (RFC 4226)
	offset := sum[len(sum)-1] & 0x0f
	code := int(sum[offset]&0x7f)<<24 |
		int(sum[offset+1])<<16 |
		int(sum[offset+2])<<8 |
		int(sum[offset+3])

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, code%int(math.Pow10(digits))), nil
}
