package observability

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestRedact_SensitiveKey(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)
	l := slog.New(redactingHandler{Handler: h})
	l.Info("test", "password", "hunter2", "token", "abc")
	out := buf.String()
	if strings.Contains(out, "hunter2") {
		t.Errorf("password value leaked: %s", out)
	}
	if strings.Contains(out, "abc") {
		t.Errorf("token value leaked: %s", out)
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Errorf("expected redaction marker: %s", out)
	}
}

func TestRedact_BearerToken(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)
	l := slog.New(redactingHandler{Handler: h})
	l.Info("req", "header", "Bearer supersecret")
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("bearer token leaked: %s", out)
	}
}

func TestCorrelationID_AttachedToLog(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(redactingHandler{Handler: h})
	ctx := WithCorrelation(context.Background(), "corr-123")
	r := slog.NewRecord(slog.Record{}.Time, slog.LevelInfo, "msg", 0)
	_ = logger.Handler().Handle(ctx, r)
	// The record is empty; test via direct handle path by logging through logger.
	logger.InfoContext(ctx, "msg")
	out := buf.String()
	if !strings.Contains(out, "corr-123") {
		t.Errorf("correlation id missing: %s", out)
	}
}

func TestCorrelationID_GeneratedIfEmpty(t *testing.T) {
	ctx := WithCorrelation(context.Background(), "")
	id := CorrelationFrom(ctx)
	if id == "" {
		t.Error("expected generated correlation id, got empty")
	}
	if len(id) != 32 {
		t.Errorf("generated id length = %d, want 32", len(id))
	}
}