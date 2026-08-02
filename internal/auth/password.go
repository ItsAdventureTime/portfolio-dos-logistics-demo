// Package auth implements password hashing (Argon2id), OTP generation and
// verification (RFC 6238 TOTP), session token generation, and CSRF token
// helpers. All dependencies are MIT, BSD, or Apache-2.0 licensed.
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters. These follow OWASP recommendations.
const (
	argon2Memory      = 64 * 1024 // 64 MB
	argon2Iterations  = 3
	argon2Parallelism  = 2
	argon2SaltLength  = 16
	argon2KeyLength   = 32
)

// HashPassword creates an Argon2id hash in PHC string format.
func HashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLength)
	return encodePHC(salt, hash), nil
}

// VerifyPassword checks a password against an Argon2id PHC hash in constant
// time. Returns (true, nil) on match, (false, nil) on mismatch, (false, err)
// if the hash is malformed.
func VerifyPassword(password, encodedHash string) (bool, error) {
	salt, hash, err := decodePHC(encodedHash)
	if err != nil {
		return false, err
	}
	other := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLength)
	if subtle.ConstantTimeCompare(hash, other) != 1 {
		return false, nil
	}
	return true, nil
}

// NeedsRehash reports whether the hash uses outdated parameters. Since we
// use fixed OWASP params, this always returns false for now. It exists so
// the service layer can call it when params are tuned in the future.
func NeedsRehash(encodedHash string) bool {
	return false
}

// DummyHash returns a pre-computed hash so that login timing stays consistent
// even when the username doesn't exist. This keeps auth errors neutral. The
// caller can't tell whether the account was found.
func DummyHash() string {
	salt := make([]byte, argon2SaltLength)
	for i := range salt {
		salt[i] = 0
	}
	hash := argon2.IDKey([]byte("dummy"), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLength)
	return encodePHC(salt, hash)
}

func encodePHC(salt, hash []byte) string {
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2Memory, argon2Iterations, argon2Parallelism, b64Salt, b64Hash)
}

func decodePHC(encoded string) (salt, hash []byte, err error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return nil, nil, errors.New("invalid hash format")
	}
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, fmt.Errorf("decode salt: %w", err)
	}
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, fmt.Errorf("decode hash: %w", err)
	}
	return salt, hash, nil
}