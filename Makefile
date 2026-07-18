.PHONY: all bin-deps generate tidy build test migrate-up migrate-down db-up db-down run docker-run docker-down

# 1. Load local .env file if it exists to import real credentials
-include .env

# Local bin directory for tool overrides
LOCAL_BIN := $(shell pwd)/bin
export PATH := $(LOCAL_BIN):$(PATH)

# 2. Extract configuration variables, providing secure fallback defaults
DB_USER     ?= postgres
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_NAME     ?= dragon_market

# Set a fallback default password *only* if it is completely absent from the .env file
# This prevents the Makefile from throwing an unexpected build error immediately.
APP_DATABASE__PASSWORD ?= password123
DB_PASSWORD            ?= $(APP_DATABASE__PASSWORD)

DB_CONN ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

all: generate build test

# Install strict tool dependencies locally
bin-deps:
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.7.2
	@GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
	@GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

# Generate strict Go server boilerplate and sqlc queries
generate: bin-deps
	@echo "Generating API boilerplate from openapi.yaml..."
	@$(LOCAL_BIN)/oapi-codegen -config api/wire/cfg.yaml api/openapi.yaml
	@echo "Generating database models and interfaces via sqlc..."
	@$(LOCAL_BIN)/sqlc generate

tidy:
	@go mod tidy

# Build the application
build: tidy
	@mkdir -p bin
	@go build -o bin/server cmd/server/main.go

# Run database container
db-up:
	@docker compose up -d postgres

# Stop database container
db-down:
	@docker compose down

# Run migrations to bring database up to date
migrate-up: bin-deps
	@echo "Running pending migrations..."
	@$(LOCAL_BIN)/migrate -path migration -database "$(DB_CONN)" up

# Rollback the last migration step
migrate-down: bin-deps
	@echo "Rolling back migrations..."
	@$(LOCAL_BIN)/migrate -path migration -database "$(DB_CONN)" down 1

# ==============================================================================
# 🚀 Local & Docker Deployment Workflows
# ==============================================================================

# Option A: Run locally (Explicitly forwards verified Makefile variables down into Go runtime environment)
run: generate db-up migrate-up
	@echo "Starting application server locally..."
	@APP_DATABASE__PASSWORD=$(DB_PASSWORD) DB_HOST=$(DB_HOST) go run cmd/server/main.go

# Option B: Run everything in Docker (Forwards matching host secrets to Docker engine orchestration context)
docker-run: generate
	@echo "Building and launching all services (App, DB, SwaggerUI) inside Docker..."
	@APP_DATABASE__PASSWORD=$(DB_PASSWORD) docker compose up --build -d
	@echo "Applying migrations to the dockerized database..."
	@sleep 2
	@$(MAKE) migrate-up

# Tear down all Docker services and clear storage spaces
docker-down:
	@docker compose down -v
