package config

import "testing"

func TestLoad_RequiresMinimum(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	t.Setenv("SESSION_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("OTP_SECRET", "0123456789abcdef0123456789abcdef")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPAddr != "127.0.0.1:8080" {
		t.Errorf("HTTPAddr default = %q, want 127.0.0.1:8080", cfg.HTTPAddr)
	}
	if !cfg.IsProduction() && cfg.IsProduction() {
		t.Errorf("IsProduction should be false in development")
	}
}

func TestLoad_RejectsShortSecrets(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("DATABASE_URL", "postgres://u:p@h:5432/d")
	t.Setenv("SESSION_SECRET", "short")
	t.Setenv("OTP_SECRET", "0123456789abcdef0123456789abcdef")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for short SESSION_SECRET, got nil")
	}
}

func TestLoad_RejectsDevCodeVisibleInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://u:p@h:5432/d")
	t.Setenv("SESSION_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("OTP_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("APP_DEV_CODE_VISIBLE", "true")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for dev code visible in production, got nil")
	}
}

func TestLoad_RejectsInvalidEnv(t *testing.T) {
	t.Setenv("APP_ENV", "nonsense")
	t.Setenv("DATABASE_URL", "postgres://u:p@h:5432/d")
	t.Setenv("SESSION_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("OTP_SECRET", "0123456789abcdef0123456789abcdef")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid APP_ENV, got nil")
	}
}