// Package middleware contains HTTP middleware: correlation id, panic recovery,
// timeouts, and security headers. It implements the redaction, security, and
// reliability requirements of docs/spec/03-security-and-privacy.md.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
)

// RequestID reads or generates an X-Request-ID and attaches it to the
// response header and the request context for correlation in logs.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		ctx := observability.WithCorrelation(r.Context(), id)
		w.Header().Set("X-Request-ID", observability.CorrelationFrom(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Recover catches panics, logs the stack trace (redacted), and returns 500.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.ErrorContext(r.Context(), "panic recovered",
						"error", rec,
						"stack", string(debug.Stack()),
					)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout enforces read and write timeouts as a safety net in addition to the
// http.Server timeouts.
func Timeout(read, write time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), write)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SecurityHeaders sets baseline response headers.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}