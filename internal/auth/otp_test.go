package auth

import (
	"testing"
	"time"
)

func TestGenerateOTP_VerifyOTP_RoundTrip(t *testing.T) {
	cfg := DefaultOTPConfig([]byte("0123456789abcdef0123456789abcdef"))
	now := time.Now()
	code, err := GenerateOTP(cfg, now)
	if err != nil {
		t.Fatalf("GenerateOTP: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("OTP length = %d, want 6", len(code))
	}
	if !VerifyOTP(cfg, code, now) {
		t.Error("VerifyOTP failed for valid code")
	}
}

func TestVerifyOTP_RejectsWrongCode(t *testing.T) {
	cfg := DefaultOTPConfig([]byte("0123456789abcdef0123456789abcdef"))
	now := time.Now()
	if VerifyOTP(cfg, "000000", now) && GenerateOTPDummy(cfg, now) != "000000" {
		t.Error("VerifyOTP accepted wrong code")
	}
}

func TestVerifyOTP_AcceptsPriorStep(t *testing.T) {
	cfg := DefaultOTPConfig([]byte("0123456789abcdef0123456789abcdef"))
	prior := time.Now().Add(-cfg.Step)
	code, _ := GenerateOTP(cfg, prior)
	now := time.Now()
	if !VerifyOTP(cfg, code, now) {
		t.Error("VerifyOTP should accept code from prior step when verifying at current time")
	}
}

func TestVerifyOTP_RejectsExpiredCode(t *testing.T) {
	cfg := DefaultOTPConfig([]byte("0123456789abcdef0123456789abcdef"))
	cfg.Window = 0
	now := time.Now()
	code, _ := GenerateOTP(cfg, now)
	old := now.Add(-2 * cfg.Step)
	if VerifyOTP(cfg, code, old) {
		t.Error("VerifyOTP accepted expired code")
	}
}

func TestHashOTPCode_VerifyOTPHash(t *testing.T) {
	code := "123456"
	hash := HashOTPCode(code)
	if hash == code {
		t.Error("hash should not equal plaintext")
	}
	if !VerifyOTPHash(code, hash) {
		t.Error("VerifyOTPHash failed for valid code")
	}
	if VerifyOTPHash("999999", hash) {
		t.Error("VerifyOTPHash accepted wrong code")
	}
}

// GenerateOTPDummy is a test helper that generates an OTP for comparison.
func GenerateOTPDummy(cfg OTPConfig, now time.Time) string {
	code, _ := GenerateOTP(cfg, now)
	return code
}