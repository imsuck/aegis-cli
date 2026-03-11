package totp_test

import (
	"testing"
	"time"
	"aegis-cli/internal/totp"
	"aegis-cli/internal/vault"
)

func TestGenerateTOTP(t *testing.T) {
	// Use a known test vector
	entry := vault.Entry{
		Type: "totp",
		Info: vault.Info{
			Secret: "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567", // JBSWY3DPEHPK3PXP
			Algo:   "SHA1",
			Digits: 6,
			Period: 30,
		},
	}

	// Test at a fixed time (Unix timestamp 0 = 1970-01-01)
	fixedTime := time.Unix(0, 0)
	code, err := totp.Generate(entry, fixedTime)
	if err != nil {
		t.Fatalf("failed to generate TOTP: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("expected 6 digit code, got %d", len(code))
	}

	// Code should be all numeric
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("expected numeric code, got %c", c)
		}
	}
}

func TestGenerateTOTPWithDifferentAlgorithms(t *testing.T) {
	tests := []struct {
		name   string
		algo   string
		secret string
	}{
		{"SHA1", "SHA1", "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"},
		{"SHA256", "SHA256", "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"},
		{"SHA512", "SHA512", "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := vault.Entry{
				Type: "totp",
				Info: vault.Info{
					Secret: tt.secret,
					Algo:   tt.algo,
					Digits: 6,
					Period: 30,
				},
			}

			code, err := totp.Generate(entry, time.Now())
			if err != nil {
				t.Fatalf("failed to generate TOTP with %s: %v", tt.algo, err)
			}

			if len(code) != 6 {
				t.Errorf("expected 6 digit code, got %d", len(code))
			}
		})
	}
}

func TestGetRemainingTime(t *testing.T) {
	entry := vault.Entry{
		Info: vault.Info{
			Period: 30,
		},
	}

	// At time 0, remaining should be 30
	remaining := totp.GetRemainingTime(entry, time.Unix(0, 0))
	if remaining != 30 {
		t.Errorf("expected 30 seconds remaining, got %d", remaining)
	}

	// At time 15, remaining should be 15
	remaining = totp.GetRemainingTime(entry, time.Unix(15, 0))
	if remaining != 15 {
		t.Errorf("expected 15 seconds remaining, got %d", remaining)
	}
}
