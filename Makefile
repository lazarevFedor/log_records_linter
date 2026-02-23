.PHONY: help build build-plugin test lint clean

PLUGIN_NAME := logs
PLUGIN_OUT := bin/$(PLUGIN_NAME).so
GO := go

help:
	@echo "Available targets:"
	@echo "  make build-plugin    - Build the golangci-lint plugin (Module Plugin System)"
	@echo "  make test            - Run all tests"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make help            - Show this help message"
	@echo "  make verify-plugin   - Verify the plugin builds correctly"

# Build the golangci-lint plugin using Module Plugin System
build-plugin:
	@echo "Building golangci-lint plugin..."
	mkdir -p bin
	$(GO) build -buildmode=plugin -o $(PLUGIN_OUT) ./plugin

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Verify the plugin builds correctly
verify-plugin: build-plugin
	@echo "Plugin built successfully: $(PLUGIN_OUT)"
	@ls -lh $(PLUGIN_OUT)

.DEFAULT_GOAL := help

