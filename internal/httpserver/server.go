// Package httpserver wires the Chi router, middleware, and the health
// endpoints. It implements the health-check requirement of
// docs/spec/05-acceptance-matrix.md.
package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/config"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/httpserver/middleware"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
	"github.com/go-chi/chi/v5"
)

// Server holds the configured router and the underlying http.Server.
type Server struct {
	router *chi.Mux
	srv    *http.Server
}

// New builds the HTTP server with middleware and the health routes registered.
func New(cfg config.Config, log *slog.Logger) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recover(log))
	r.Use(middleware.Timeout(cfg.ReadTimeout, cfg.WriteTimeout))
	r.Use(middleware.SecurityHeaders)

	r.Get("/healthz", healthLiveness)
	r.Get("/healthz/ready", healthReady)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return &Server{router: r, srv: srv}
}

// Router returns the underlying Chi router so feature modules can mount
// their routes during Stage 3+.
func (s *Server) Router() *chi.Mux { return s.router }

// Start begins listening and blocks until the server stops.
func (s *Server) Start() error { return s.srv.ListenAndServe() }

// Shutdown gracefully stops the server within the given timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func writeJSON(w http.ResponseWriter, status int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func healthLiveness(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"check":  "liveness",
		"now":    time.Now().UTC().Format(time.RFC3339),
	})
}

func healthReady(w http.ResponseWriter, r *http.Request) {
	// Stage 2: always ready once the process is up. Stage 3+ will add DB
	// connectivity and migration checks here.
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"check":  "ready",
		"now":    time.Now().UTC().Format(time.RFC3339),
	})
}

// _ ensures observability import is retained for correlation helpers used by
// middleware in later stages.
var _ = observability.CorrelationFrom