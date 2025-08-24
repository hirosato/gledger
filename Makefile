# Makefile for gledger - Go implementation of ledger-cli

# Variables
BINARY_NAME=gledger
MAIN_PATH=cmd/gledger/main.go
BUILD_DIR=build
COVERAGE_FILE=coverage.out
LEDGER_CMD=ledger

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Build flags
LDFLAGS=-ldflags "-s -w"
TESTFLAGS=-v -race -coverprofile=$(COVERAGE_FILE)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test test-unit test-integration test-compare coverage fmt vet lint install deps help

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'

## all: Format, vet, lint, test, and build
all: fmt vet lint test build

## build: Build the gledger binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## clean: Remove build artifacts and test coverage
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE)
	@echo "$(GREEN)Clean complete$(NC)"

## test: Run all tests
test: test-unit test-integration test-specs

## test-unit: Run unit tests
test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GOTEST) -v -race ./domain/... ./application/... ./infrastructure/... ./interfaces/...

## test-integration: Run integration tests
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GOTEST) -v ./test/integration/...

## test-compare: Run comparison tests with ledger-cli
test-compare: build
	@echo "$(GREEN)Running comparison tests with ledger-cli...$(NC)"
	@if ! command -v $(LEDGER_CMD) > /dev/null; then \
		echo "$(RED)Error: ledger-cli not found. Please install ledger-cli first.$(NC)"; \
		exit 1; \
	fi
	$(GOTEST) -v ./test/integration/compare/...

## test-specs: Run spec-based tests from original ledger
test-specs: build
	@echo "$(GREEN)Running spec-based tests...$(NC)"
	$(GOTEST) -v ./test/specs/...

## test-runner: Run spec tests using test runner
test-runner: build
	@echo "$(GREEN)Building test runner...$(NC)"
	$(GOBUILD) -o $(BUILD_DIR)/testrunner cmd/testrunner/main.go
	@echo "$(GREEN)Running spec tests with test runner...$(NC)"
	$(BUILD_DIR)/testrunner -dir ../ledger/test/baseline -list

## coverage: Run tests with coverage and display report
coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GOTEST) $(TESTFLAGS) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

## fmt: Format all Go files
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) ./...
	@echo "$(GREEN)Format complete$(NC)"

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GOVET) ./...
	@echo "$(GREEN)Vet complete$(NC)"

## lint: Run golangci-lint
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if ! command -v $(GOLINT) > /dev/null; then \
		echo "$(YELLOW)Installing golangci-lint...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	$(GOLINT) run ./...
	@echo "$(GREEN)Lint complete$(NC)"

## deps: Download and tidy dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

## install: Install gledger to GOPATH/bin
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	$(GOCMD) install $(MAIN_PATH)
	@echo "$(GREEN)$(BINARY_NAME) installed to GOPATH/bin$(NC)"

## benchmark: Run benchmarks
benchmark:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...

## import-tests: Import test files from original ledger
import-tests:
	@echo "$(GREEN)Importing test files from original ledger...$(NC)"
	@./scripts/import-tests.sh
	@echo "$(GREEN)Test import complete$(NC)"

## watch: Run tests on file change (requires entr)
watch:
	@if ! command -v entr > /dev/null; then \
		echo "$(RED)Error: entr not found. Install with: brew install entr$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Watching for changes...$(NC)"
	@find . -name "*.go" | entr -c make test

## ci: Run CI pipeline locally
ci: deps fmt vet lint test build
	@echo "$(GREEN)CI pipeline complete$(NC)"

# Default target
.DEFAULT_GOAL := help