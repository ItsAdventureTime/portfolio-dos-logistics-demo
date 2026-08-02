// Package main is the composition root for the API server. It loads config,
// builds the server, and handles graceful shutdown on SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/config"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/demo"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/httpserver"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/httpserver/middleware"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
)

func main() {
	if err := run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := observability.Logger(cfg.LogLevel)

	if cfg.DemoMode {
		middleware.InDemoMode = true
		log.Info("starting in demo mode (in-memory stores, no database)",
			"env", cfg.Env, "addr", cfg.HTTPAddr)
	}

	srv := httpserver.New(cfg, log)

	if cfg.DemoMode {
		bootstrapped := demo.Bootstrap(log)
		srv.MountAuth(bootstrapped.AuthService, bootstrapped.OTPCfg)
		srv.MountWorkflow(bootstrapped.QuotationService, bootstrapped.WorkflowService, bootstrapped.AuthService.ValidateSession)
		log.Info("demo auth and workflow routes mounted",
			"username", "admin",
			"password", "Password123!",
		)

		// Auto-reset every 30 minutes so each prospect starts fresh.
		resetInterval := 30 * time.Minute
		go func() {
			ticker := time.NewTicker(resetInterval)
			defer ticker.Stop()
			for range ticker.C {
				bootstrapped.ResetFn()
			}
		}()
		log.Info("demo auto-reset scheduled", "interval", resetInterval.String())
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		log.Info("shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return srv.Shutdown(ctx)
}