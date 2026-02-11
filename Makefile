.PHONY: all build test clean coverage lint fmt vet install help

# Variables
BINARY_NAME=knife
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO=go
GOFLAGS=-v

# Default target
all: fmt vet test build

## build: Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) ./...

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race ./...

## test-cover: Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean ./...
	rm -f coverage.out coverage.html
	rm -rf dist/

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## lint: Run golangci-lint (if available)
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GO) mod tidy

## install: Install the project
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install ./...

## release: Run goreleaser (for releases)
release:
	@echo "Running goreleaser..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
	else \
		echo "goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest"; \
	fi

## release-snapshot: Run goreleaser in snapshot mode
release-snapshot:
	@echo "Running goreleaser (snapshot)..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest"; \
	fi

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
