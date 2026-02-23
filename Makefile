.PHONY: help build-cli build-plugin test lint

PLUGIN_NAME := logs
GO := go
CUSTOM_BIN := custom-gcl

help:
	@echo "Available targets:"
	@echo "  make build-cli       - Build the CLI tool"
	@echo "  make vet    		  - Run linter as CLI tool through go vet"
	@echo "  make vet-fix    	  - Run linter as CLI tool through go vet with fix option"
	@echo "  make build-plugin    - Build custom golangci-lint binary with embedded logs plugin"
	@echo "  make lint            - Run linting with the custom binary"
	@echo "  make lint-fix        - Run linting with the custom binary and apply fixes"
	@echo "  make test            - Run all tests"
	@echo "  make help            - Show this help message"

build-cli:
	@echo "Building CLI tool..."
	$(GO) build -o logs ./cmd/main.go

vet:
	@echo "Running go vet..."
	$(GO) vet -vettool=./logs ./...

vet-fix:
	@echo "Running go vet with fix option..."
	$(GO) vet -vettool=./logs -fix ./...

build-plugin:
	@echo "Building custom golangci-lint binary with logs plugin..."
	golangci-lint custom

lint:
	@echo "Running linting with custom binary..."
	./$(CUSTOM_BIN) run -v

lint-fix:
	@echo "Running linting with custom binary and applying fixes..."
	./$(CUSTOM_BIN) run -v --fix

test:
	@echo "Running tests..."
	$(GO) test -v -cover ./...

.DEFAULT_GOAL := help

