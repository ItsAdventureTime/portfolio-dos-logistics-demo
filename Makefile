.PHONY: build test test-unit test-integration lint vet typecheck fmt migrate-up migrate-down migrate-status run web-dev web-build web-test web-lint web-typecheck secret-scan all

DEV := ./scripts/dos-dev

build:
	$(DEV) go build -o bin/api ./cmd/api
	$(DEV) go build -o bin/migrate ./cmd/migrate

test:
	$(DEV) go test ./...

test-unit:
	$(DEV) go test -short ./...

test-integration:
	$(DEV) go test -run Integration ./...

vet:
	$(DEV) go vet ./...

lint: vet
	$(DEV) golangci-lint run ./...

typecheck:
	$(DEV) web npx tsc --noEmit -p tsconfig.app.json

fmt:
	$(DEV) gofmt -w .

migrate-up:
	$(DEV) goose -dir internal/store/migrations "$(DATABASE_URL)" up

migrate-down:
	$(DEV) goose -dir internal/store/migrations "$(DATABASE_URL)" down

migrate-status:
	$(DEV) goose -dir internal/store/migrations "$(DATABASE_URL)" status

run:
	go run ./cmd/api

web-dev:
	cd web && npm run dev

web-build:
	$(DEV) web npm run build

web-test:
	$(DEV) web npx vitest run

web-lint:
	$(DEV) web npx oxlint src

web-typecheck:
	$(DEV) web npx tsc --noEmit -p tsconfig.app.json

secret-scan:
	$(DEV) gitleaks detect --source . --no-banner --config .gitleaks.toml

all: lint test web-typecheck web-test web-build build