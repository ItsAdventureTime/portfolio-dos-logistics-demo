// Package main is the migration runner entry point. It delegates to goose so
// migrations can be applied with: go run ./cmd/migrate up|down|status.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/config"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/store/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: migrate [up|down|status|redo]")
		os.Exit(2)
	}
	if err := migrations.Run(cfg.DatabaseURL, os.Args[1], os.Args[2:]); err != nil {
		log.Fatalf("migrate: %v", err)
	}
}