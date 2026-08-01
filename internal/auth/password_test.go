package auth

import (
	"testing"
)

func TestHashPassword_VerifyPassword_RoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Error("hash should not equal plaintext")
	}
	ok, err := VerifyPassword("correct horse battery staple", hash)
	if err != nil || !ok {
		t.Errorf("VerifyPassword correct: ok=%v err=%v", ok, err)
	}
	ok, err = VerifyPassword("wrong password", hash)
	if err != nil || ok {
		t.Errorf("VerifyPassword wrong: ok=%v err=%v", ok, err)
	}
}

func TestVerifyPassword_MalformedHash(t *testing.T) {
	ok, err := VerifyPassword("pw", "not-a-valid-hash")
	if ok {
		t.Error("malformed hash should not verify")
	}
	if err == nil {
		t.Error("malformed hash should return error")
	}
}

func TestDummyHash_DoesNotVerify(t *testing.T) {
	dh := DummyHash()
	ok, err := VerifyPassword("anything", dh)
	if err != nil {
		t.Errorf("dummy hash verify error: %v", err)
	}
	if ok {
		t.Error("dummy hash should not verify against any password")
	}
}

func TestHashPassword_DifferentSalts(t *testing.T) {
	h1, _ := HashPassword("same")
	h2, _ := HashPassword("same")
	if h1 == h2 {
		t.Error("same password should produce different hashes (different salts)")
	}
}