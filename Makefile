.PHONY: all bin-deps generate tidy build test migrate-up migrate-down

# Local bin directory for tool overrides
LOCAL_BIN := $(shell pwd)/bin
export PATH := $(LOCAL_BIN):$(PATH)

all: generate build test

# Install strict tool dependencies locally
bin-deps:
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

# Generate strict Go server boilerplate from OpenAPI spec
generate: bin-deps
	@echo "Generating API boilerplate from openapi.yaml..."
	@$(LOCAL_BIN)/oapi-codegen -config api/wire/cfg.yaml api/openapi.yaml

tidy:
	@go mod tidy

# Build the application
build: tidy
	@mkdir -p bin
	@go build -o bin/server cmd/server/main.go

# Run tests
test:
	@go test -v -race ./...
