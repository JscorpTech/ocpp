.PHONY: test test-unit test-integration build run clean help

# Test commands
test: ## Run all tests
	go test ./... -v

test-unit: ## Run unit tests only
	go test ./... -v -short

test-coverage: ## Run tests with coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Build commands
build: ## Build the application
	go build -o bin/ocpp ./cmd/main.go

run: ## Run the application
	go run ./cmd/main.go

# Development commands
dev: ## Run with hot reload (requires air)
	air

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

# Utility commands
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	go mod download
	go mod tidy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
