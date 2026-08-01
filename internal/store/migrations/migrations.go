// Package migrations wraps goose so the migration runner and tests can apply
// migrations programmatically. Migration SQL files live in this directory.
package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Run executes a goose command against the given database URL.
// command is one of: up, down, status, redo, reset.
func Run(databaseURL, command string, args []string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	if err := goose.RunContext(context.Background(), command, db, ".", args...); err != nil {
		return fmt.Errorf("goose %s: %w", command, err)
	}
	return nil
}