package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

// GenerateSessionToken creates a 256-bit random session token. Returns
// (plaintext, hash, error). Only the hash is stored; the plaintext is
// given to the client via a cookie and never persisted server-side.
func GenerateSessionToken() (plaintext, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}
	plaintext = hex.EncodeToString(b)
	hash = HashToken(plaintext)
	return plaintext, hash, nil
}

// HashToken returns the SHA-256 hex hash of a token. Used for session
// token and CSRF token hashing.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// GenerateCSRFToken creates a random 256-bit CSRF token bound to the
// session via HMAC. The token is sent as a cookie and must be echoed in
// the X-CSRF-Token header for state-changing requests (double-submit).
func GenerateCSRFToken(sessionHash string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate csrf token: %w", err)
	}
	// Combine random bytes with the session hash so the CSRF token is
	// bound to the session without a shared secret beyond the session itself.
	input := append(b, []byte(sessionHash)...)
	h := sha256.Sum256(input)
	return hex.EncodeToString(h[:]), nil
}

// VerifyCSRFToken checks a submitted CSRF token against the expected value
// in constant time.
func VerifyCSRFToken(submitted, expected string) bool {
	if len(submitted) != len(expected) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(submitted), []byte(expected)) == 1
}