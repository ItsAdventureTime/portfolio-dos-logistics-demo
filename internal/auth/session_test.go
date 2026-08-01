package auth

import (
	"testing"
)

func TestGenerateSessionToken_UniqueAndHashed(t *testing.T) {
	t1, h1, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken: %v", err)
	}
	t2, h2, _ := GenerateSessionToken()
	if t1 == t2 {
		t.Error("tokens should be unique")
	}
	if h1 == h2 {
		t.Error("hashes should differ")
	}
	if t1 == h1 {
		t.Error("plaintext should not equal hash")
	}
	if len(t1) != 64 {
		t.Errorf("token length = %d, want 64", len(t1))
	}
}

func TestHashToken_Deterministic(t *testing.T) {
	h1 := HashToken("abc")
	h2 := HashToken("abc")
	if h1 != h2 {
		t.Error("HashToken should be deterministic")
	}
	if HashToken("abc") == HashToken("abd") {
		t.Error("HashToken should differ for different inputs")
	}
}

func TestGenerateCSRFToken_SessionBound(t *testing.T) {
	csrf1, err := GenerateCSRFToken("session-hash-1")
	if err != nil {
		t.Fatalf("GenerateCSRFToken: %v", err)
	}
	csrf2, _ := GenerateCSRFToken("session-hash-2")
	if csrf1 == csrf2 {
		t.Error("CSRF tokens should differ for different sessions")
	}
}

func TestVerifyCSRFToken(t *testing.T) {
	token := "deadbeef12345678"
	if !VerifyCSRFToken(token, token) {
		t.Error("VerifyCSRFToken failed for matching tokens")
	}
	if VerifyCSRFToken(token, "different") {
		t.Error("VerifyCSRFToken accepted mismatch")
	}
	if VerifyCSRFToken("short", "different-length") {
		t.Error("VerifyCSRFToken accepted different-length tokens")
	}
}