.PHONY: build test test-unit test-integration lint vet typecheck fmt migrate-up migrate-down migrate-status run dev web-dev web-build web-test secret-scan all

GOPATH_BIN := $(shell go env GOPATH)/bin

build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate

test:
	go test ./...

test-unit:
	go test -short ./...

test-integration:
	go test -run Integration ./...

vet:
	go vet ./...

lint: vet
	$(GOPATH_BIN)/golangci-lint run ./...

typecheck:
	cd web && npx tsc --noEmit -p tsconfig.app.json

fmt:
	gofmt -w .
	gofumpt -w .

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-status:
	go run ./cmd/migrate status

run:
	go run ./cmd/api

dev:
	go run ./cmd/api &

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

web-test:
	cd web && npm run test

web-lint:
	cd web && npm run lint

web-typecheck:
	cd web && npm run typecheck

secret-scan:
	$(GOPATH_BIN)/gitleaks detect --source . --no-banner

all: lint test web-typecheck web-test web-build build