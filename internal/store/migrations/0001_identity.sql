-- +goose Up
-- +goose StatementBegin

-- Users: the single authenticated account type (Administrator).
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        TEXT NOT NULL UNIQUE,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    display_name    TEXT NOT NULL,
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sessions: server-managed, hash-only storage (never store the plaintext token).
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      TEXT NOT NULL UNIQUE,
    role_preview    TEXT,  -- null = default Administrator; preview role otherwise
    ip_address      TEXT,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Email challenges: one-time OTP codes for email verification.
CREATE TABLE email_challenges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash       TEXT NOT NULL,
    purpose         TEXT NOT NULL DEFAULT 'email_verification',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL,
    consumed_at     TIMESTAMPTZ
);

CREATE INDEX idx_email_challenges_user_id ON email_challenges(user_id);

-- Audit events: append-only. No UPDATE or DELETE ever permitted.
-- Enforced via trigger + REVOKE at the end of this migration.
CREATE TABLE audit_events (
    id              BIGSERIAL PRIMARY KEY,
    correlation_id  TEXT NOT NULL,
    actor_user_id   UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_role      TEXT NOT NULL,
    action          TEXT NOT NULL,
    entity_type     TEXT NOT NULL,
    entity_id       TEXT,
    details         JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_events_correlation_id ON audit_events(correlation_id);
CREATE INDEX idx_audit_events_actor_user_id ON audit_events(actor_user_id);
CREATE INDEX idx_audit_events_entity_type ON audit_events(entity_type);
CREATE INDEX idx_audit_events_created_at ON audit_events(created_at);

-- Append-only enforcement: block any UPDATE or DELETE on audit_events.
CREATE OR REPLACE FUNCTION block_audit_mutation() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'audit_events is append-only: UPDATE and DELETE are not permitted';
END;
$$;

CREATE TRIGGER no_update_audit_events
    BEFORE UPDATE ON audit_events
    FOR EACH ROW EXECUTE FUNCTION block_audit_mutation();

CREATE TRIGGER no_delete_audit_events
    BEFORE DELETE ON audit_events
    FOR EACH ROW EXECUTE FUNCTION block_audit_mutation();

-- Explicitly revoke UPDATE and DELETE from all roles.
REVOKE UPDATE, DELETE ON audit_events FROM PUBLIC;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_events;
DROP TABLE IF EXISTS email_challenges;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS block_audit_mutation;
-- +goose StatementEnd