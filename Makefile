.PHONY: all bin-deps generate tidy build test migrate-up migrate-down db-up db-down

# 1. Load local .env file if it exists to import real credentials
-include .env

# Local bin directory for tool overrides
LOCAL_BIN := $(shell pwd)/bin
export PATH := $(LOCAL_BIN):$(PATH)

# 2. Build connection string using .env variables, falling back to defaults if empty
DB_USER     ?= postgres
DB_PASSWORD ?= $(APP_DATABASE__PASSWORD)
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_NAME     ?= dragon_market

ifeq ($(DB_PASSWORD),)
  $(error APP_DATABASE__PASSWORD is not set. Please check your local .env file)
endif

DB_CONN ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
all: generate build test

# Install strict tool dependencies locally
bin-deps:
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
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
	@docker compose up -d

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
