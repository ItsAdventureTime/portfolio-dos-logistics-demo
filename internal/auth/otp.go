package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"time"
)

// OTPConfig controls OTP generation and verification.
type OTPConfig struct {
	Digits    int           // 6 per RFC 6238
	Step      time.Duration // 30s per RFC 6238
	Window    int           // how many steps back to accept (1 = current + 1 prior)
	SecretKey []byte         // HMAC signing key (>=32 bytes)
}

// DefaultOTPConfig returns RFC 6238-compliant defaults.
func DefaultOTPConfig(secret []byte) OTPConfig {
	return OTPConfig{
		Digits:    6,
		Step:      30 * time.Second,
		Window:    1,
		SecretKey: secret,
	}
}

// GenerateOTP produces a time-based 6-digit code and returns it as a zero-padded
// string. The code is valid for the current time step.
func GenerateOTP(cfg OTPConfig, now time.Time) (string, error) {
	counter := uint64(now.Unix()) / uint64(cfg.Step.Seconds())
	return generateTOTP(cfg.SecretKey, counter, cfg.Digits)
}

// VerifyOTP checks a code against the current and prior time steps.
func VerifyOTP(cfg OTPConfig, code string, now time.Time) bool {
	counter := uint64(now.Unix()) / uint64(cfg.Step.Seconds())
	for i := 0; i <= cfg.Window; i++ {
		expected, err := generateTOTP(cfg.SecretKey, counter-uint64(i), cfg.Digits)
		if err != nil {
			continue
		}
		if hmac.Equal([]byte(expected), []byte(code)) {
			return true
		}
	}
	return false
}

// HashOTPCode hashes a code with SHA-256 for storage so the plaintext code
// is not persisted. Returns a hex string.
func HashOTPCode(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}

// VerifyOTPHash compares a plaintext code against a stored hash in constant time.
func VerifyOTPHash(code, storedHash string) bool {
	h := sha256.Sum256([]byte(code))
	return hmac.Equal([]byte(hex.EncodeToString(h[:])), []byte(storedHash))
}

// GenerateOTPSecret returns a random base32-encoded secret for TOTP.
// Not used for email-verification OTP (which uses the app-wide OTP_SECRET),
// but available for future authenticator-app flows.
func GenerateOTPSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate otp secret: %w", err)
	}
	return base32.StdEncoding.EncodeToString(b), nil
}

func generateTOTP(secret []byte, counter uint64, digits int) (string, error) {
	msg := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		msg[i] = byte(counter & 0xFF)
		counter >>= 8
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(msg)
	hash := mac.Sum(nil)
	offset := int(hash[len(hash)-1] & 0x0F)
	binary := (uint32(hash[offset])&0x7F)<<24 |
		uint32(hash[offset+1])<<16 |
		uint32(hash[offset+2])<<8 |
		uint32(hash[offset+3])
	otp := binary % uint32(pow10(digits))
	return fmt.Sprintf("%0*d", digits, otp), nil
}

func pow10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}