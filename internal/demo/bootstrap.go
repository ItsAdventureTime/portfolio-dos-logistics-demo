// Package demo wires the in-memory stores and seeds the demo user for local
// demo mode. This avoids needing PostgreSQL to run the demo.
package demo

import (
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/service"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/store"
	"log/slog"
)

// BootstrapResult holds the seeded services for demo mode.
type BootstrapResult struct {
	AuthService *service.AuthService
	OTPCfg      auth.OTPConfig
}

// Bootstrap creates in-memory stores, seeds a demo admin user, and returns
// a ready-to-use AuthService.
func Bootstrap(log *slog.Logger) BootstrapResult {
	users := store.NewMemUserStore()
	sessions := store.NewMemSessionStore()
	challenges := store.NewMemEmailChallengeStore()
	audit := store.NewMemAuditStore()

	// Seed the demo admin user with a known password.
	hash, _ := auth.HashPassword("Password123!")
	now := time.Now().UTC()
	admin := &domain.User{
		ID:            domain.UserID("demo-admin-001"),
		Username:      "admin",
		Email:         "admin@dosfreightflow.example",
		PasswordHash:  hash,
		DisplayName:   "Demo Administrator",
		EmailVerified: true,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	users.Seed(admin)

	otpCfg := auth.DefaultOTPConfig([]byte("demo-otp-secret-32-bytes-minimum-ok"))
	svc := service.NewAuthService(
		users, sessions, challenges, audit,
		otpCfg,
		24*time.Hour, 1*time.Hour,
		true, // dev code visible for demo
		log,
	)

	return BootstrapResult{
		AuthService: svc,
		OTPCfg:      otpCfg,
	}
}