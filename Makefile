# GoPackKit Makefile

.PHONY: help build test lint fmt vet clean deps coverage security release

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Project info
PROJECT_NAME=gopackkit
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Colors for terminal output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Display this help message
	@echo "$(GREEN)GoPackKit - Development Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(BLUE)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download and verify dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated successfully!$(NC)"

fmt: ## Format Go code
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@echo "$(GREEN)Code formatted successfully!$(NC)"

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	$(GOCMD) vet ./...
	@echo "$(GREEN)Vet check passed!$(NC)"

lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	$(GOLINT) run --timeout=5m
	@echo "$(GREEN)Linting passed!$(NC)"

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests passed!$(NC)"

test-coverage: test ## Run tests with coverage report
	@echo "$(YELLOW)Generating coverage report...$(NC)"
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

test-short: ## Run tests (short mode)
	@echo "$(YELLOW)Running tests (short mode)...$(NC)"
	$(GOTEST) -short ./...
	@echo "$(GREEN)Short tests passed!$(NC)"

test-verbose: ## Run tests with verbose output
	@echo "$(YELLOW)Running tests (verbose)...$(NC)"
	$(GOTEST) -v ./...
	@echo "$(GREEN)Verbose tests passed!$(NC)"

security: ## Run security scan
	@echo "$(YELLOW)Running security scan...$(NC)"
	gosec ./...
	@echo "$(GREEN)Security scan completed!$(NC)"

vulnerability: ## Check for vulnerabilities
	@echo "$(YELLOW)Checking for vulnerabilities...$(NC)"
	govulncheck ./...
	@echo "$(GREEN)Vulnerability check completed!$(NC)"

build: ## Build the project
	@echo "$(YELLOW)Building project...$(NC)"
	$(GOBUILD) -v ./...
	@echo "$(GREEN)Build successful!$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean completed!$(NC)"

install-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) -u github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) -u golang.org/x/vuln/cmd/govulncheck@latest
	@echo "$(GREEN)Development tools installed!$(NC)"

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "$(GREEN)All checks passed!$(NC)"

ci: deps check security ## Run CI pipeline locally
	@echo "$(GREEN)CI pipeline completed successfully!$(NC)"

release-check: ## Check if ready for release
	@echo "$(YELLOW)Checking release readiness...$(NC)"
	@if [ -z "$(shell git status --porcelain)" ]; then \
		echo "$(GREEN)Working directory is clean$(NC)"; \
	else \
		echo "$(RED)Working directory has uncommitted changes$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Running full test suite...$(NC)"
	@$(MAKE) ci
	@echo "$(GREEN)Ready for release!$(NC)"

info: ## Show project information
	@echo "$(BLUE)Project Information:$(NC)"
	@echo "  Name: $(PROJECT_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Go Version: $(shell $(GOCMD) version)"

benchmark: ## Run benchmarks
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...
	@echo "$(GREEN)Benchmarks completed!$(NC)"

profile: ## Run with profiling
	@echo "$(YELLOW)Running with CPU profiling...$(NC)"
	$(GOTEST) -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo "$(GREEN)Profiling completed! Files: cpu.prof, mem.prof$(NC)"

doc: ## Generate and serve documentation
	@echo "$(YELLOW)Starting documentation server...$(NC)"
	@echo "$(GREEN)Documentation available at: http://localhost:6060/pkg/github.com/saipulimdn/gopackkit/$(NC)"
	godoc -http=:6060

update-deps: ## Update all dependencies
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

docker-test: ## Run tests in Docker container
	@echo "$(YELLOW)Running tests in Docker...$(NC)"
	docker run --rm -v $(PWD):/app -w /app golang:1.21 make test
	@echo "$(GREEN)Docker tests completed!$(NC)"

# Example targets for different environments
test-dev: ## Run tests for development environment
	@echo "$(YELLOW)Running development tests...$(NC)"
	$(GOTEST) -tags=dev ./...

test-prod: ## Run tests for production environment
	@echo "$(YELLOW)Running production tests...$(NC)"
	$(GOTEST) -tags=prod ./...

# Git hooks
install-hooks: ## Install git pre-commit hooks
	@echo "$(YELLOW)Installing git hooks...$(NC)"
	@echo '#!/bin/sh\nmake check' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Git hooks installed!$(NC)"
