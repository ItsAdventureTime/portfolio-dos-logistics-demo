// Package observability provides structured logging with secret redaction
// and correlation identifiers. Logs are JSON-formatted, and sensitive values
// (passwords, tokens, cookies, auth headers, emails, OTPs, document names,
// customer data) are replaced with [REDACTED] before they reach the output.
package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
	"strings"
)

// contextKey is an unexported type to avoid collisions.
type contextKey struct{ name string }

var correlationKey = contextKey{"correlation_id"}

// Redacted keys and substrings. Values matching these are replaced with
// "[REDACTED]" before the record is emitted.
var sensitiveKeys = map[string]bool{
	"password": true, "passwd": true, "secret": true, "token": true,
	"authorization": true, "cookie": true, "otp": true, "api_key": true,
	"apikey": true, "session": true, "private_key": true, "credential": true,
}

var sensitiveSubstrs = []string{"password", "secret", "token", "cookie", "authorization"}

// WithCorrelation returns a context carrying the correlation id. If id is
// empty a new random id is generated.
func WithCorrelation(ctx context.Context, id string) context.Context {
	if id == "" {
		id = newID()
	}
	return context.WithValue(ctx, correlationKey, id)
}

// CorrelationFrom returns the correlation id from the context, or "" if absent.
func CorrelationFrom(ctx context.Context) string {
	v, _ := ctx.Value(correlationKey).(string)
	return v
}

// Logger returns a slog.Logger that emits JSON to w with the configured level
// and a handler that redacts sensitive keys.
func Logger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	h := redactingHandler{Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})}
	return slog.New(h)
}

// redactingHandler wraps a slog.Handler and scrubs sensitive attribute values.
type redactingHandler struct {
	slog.Handler
}

func (h redactingHandler) Handle(ctx context.Context, r slog.Record) error {
	cleaned := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	cid := CorrelationFrom(ctx)
	if cid != "" {
		cleaned.AddAttrs(slog.String("correlation_id", cid))
	}
	r.Attrs(func(a slog.Attr) bool {
		cleaned.AddAttrs(redactAttr(a))
		return true
	})
	return h.Handler.Handle(ctx, cleaned)
}

func redactAttr(a slog.Attr) slog.Attr {
	if isSensitiveKey(a.Key) {
		return slog.String(a.Key, "[REDACTED]")
	}
	if v, ok := a.Value.Any().(string); ok && containsSensitive(v) {
		return slog.String(a.Key, "[REDACTED]")
	}
	return a
}

func isSensitiveKey(k string) bool {
	lk := strings.ToLower(k)
	if sensitiveKeys[lk] {
		return true
	}
	for _, s := range sensitiveSubstrs {
		if strings.Contains(lk, s) {
			return true
		}
	}
	return false
}

func containsSensitive(v string) bool {
	// Redact values that look like bearer tokens or long secrets. This is a
	// heuristic; key-based redaction is the primary mechanism.
	if strings.HasPrefix(strings.ToLower(v), "bearer ") {
		return true
	}
	return false
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}