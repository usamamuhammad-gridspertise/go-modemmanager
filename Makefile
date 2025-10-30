# Makefile for go-modemmanager and mmctl CLI
.PHONY: all build install clean test lint fmt help

# Binary names
BINARY_NAME=mmctl
INSTALL_PATH=/usr/local/bin

# Build information
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

all: build ## Build the CLI binary

help: ## Show this help message
	@echo "$(CYAN)go-modemmanager Makefile$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build the mmctl CLI binary
	@echo "$(CYAN)Building mmctl...$(NC)"
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/mmctl
	@echo "$(GREEN)✓ Build complete: $(BINARY_NAME)$(NC)"

build-all: ## Build for multiple platforms
	@echo "$(CYAN)Building for multiple platforms...$(NC)"
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-linux-amd64 ./cmd/mmctl
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-linux-arm64 ./cmd/mmctl
	GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -o build/$(BINARY_NAME)-linux-armv7 ./cmd/mmctl
	@echo "$(GREEN)✓ Multi-platform build complete$(NC)"
	@ls -lh build/

install: build ## Install mmctl to system path
	@echo "$(CYAN)Installing mmctl to $(INSTALL_PATH)...$(NC)"
	sudo install -m 755 $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "$(GREEN)✓ Installed successfully$(NC)"
	@echo "Run 'mmctl --help' to get started"

uninstall: ## Uninstall mmctl from system
	@echo "$(CYAN)Uninstalling mmctl...$(NC)"
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled successfully$(NC)"

clean: ## Remove build artifacts
	@echo "$(CYAN)Cleaning build artifacts...$(NC)"
	rm -f $(BINARY_NAME)
	rm -rf build/
	go clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

test: ## Run tests
	@echo "$(CYAN)Running tests...$(NC)"
	go test -v ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

test-mock: ## Run tests with mock ModemManager
	@echo "$(CYAN)Starting mock ModemManager...$(NC)"
	@echo "$(YELLOW)Make sure D-Bus is running and real ModemManager is stopped$(NC)"
	cd test-environment/mock-dbus && sudo ./start-mock.sh &
	@sleep 3
	@echo "$(CYAN)Running tests...$(NC)"
	go test -v ./...
	@echo "$(GREEN)✓ Mock tests complete$(NC)"

lint: ## Run linter
	@echo "$(CYAN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Lint complete$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not found, running basic checks...$(NC)"; \
		go vet ./...; \
		echo "$(GREEN)✓ Basic checks complete$(NC)"; \
	fi

fmt: ## Format code
	@echo "$(CYAN)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Format complete$(NC)"

tidy: ## Tidy go modules
	@echo "$(CYAN)Tidying go modules...$(NC)"
	go mod tidy
	@echo "$(GREEN)✓ Tidy complete$(NC)"

generate: ## Run go generate
	@echo "$(CYAN)Running go generate...$(NC)"
	go generate ./...
	@echo "$(GREEN)✓ Generate complete$(NC)"

deps: ## Download dependencies
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

verify: lint test ## Verify code (lint + test)
	@echo "$(GREEN)✓ Verification complete$(NC)"

docker-build: ## Build Docker test environment
	@echo "$(CYAN)Building Docker test environment...$(NC)"
	cd test-environment && docker-compose build
	@echo "$(GREEN)✓ Docker build complete$(NC)"

docker-test: ## Run tests in Docker container
	@echo "$(CYAN)Running tests in Docker...$(NC)"
	cd test-environment && docker-compose run --rm modemmanager-test bash -c "cd /workspace && go test -v ./..."
	@echo "$(GREEN)✓ Docker tests complete$(NC)"

docker-shell: ## Open shell in Docker test environment
	@echo "$(CYAN)Starting Docker shell...$(NC)"
	cd test-environment && docker-compose run --rm modemmanager-test

run-example: build ## Run with example (list modems)
	@echo "$(CYAN)Running mmctl list...$(NC)"
	./$(BINARY_NAME) list || echo "$(YELLOW)Note: Requires ModemManager running and modems connected$(NC)"

demo: build ## Run CLI demo
	@echo "$(CYAN)=== mmctl Demo ===$(NC)"
	@echo ""
	@echo "$(GREEN)1. Version:$(NC)"
	./$(BINARY_NAME) --version || true
	@echo ""
	@echo "$(GREEN)2. List modems:$(NC)"
	./$(BINARY_NAME) list || echo "$(YELLOW)No modems found$(NC)"
	@echo ""
	@echo "$(GREEN)3. Help:$(NC)"
	./$(BINARY_NAME) --help
	@echo ""
	@echo "$(CYAN)=== Demo Complete ===$(NC)"

watch: ## Watch for changes and rebuild (requires entr)
	@echo "$(CYAN)Watching for changes...$(NC)"
	@if command -v entr > /dev/null; then \
		find . -name '*.go' | entr -r make build; \
	else \
		echo "$(RED)entr not found. Install with: apt-get install entr$(NC)"; \
		exit 1; \
	fi

release: clean verify build-all ## Prepare release (clean, verify, build all)
	@echo "$(GREEN)✓ Release build complete$(NC)"
	@echo "$(CYAN)Binaries in build/ directory:$(NC)"
	@ls -lh build/

info: ## Show build information
	@echo "$(CYAN)Build Information:$(NC)"
	@echo "  Binary name:   $(BINARY_NAME)"
	@echo "  Version:       $(VERSION)"
	@echo "  Git commit:    $(GIT_COMMIT)"
	@echo "  Go version:    $(shell go version)"
	@echo "  Install path:  $(INSTALL_PATH)"

check-deps: ## Check for required dependencies
	@echo "$(CYAN)Checking dependencies...$(NC)"
	@echo -n "Go:              "
	@which go > /dev/null && echo "$(GREEN)✓$(NC)" || echo "$(RED)✗ Not found$(NC)"
	@echo -n "Git:             "
	@which git > /dev/null && echo "$(GREEN)✓$(NC)" || echo "$(RED)✗ Not found$(NC)"
	@echo -n "Docker:          "
	@which docker > /dev/null && echo "$(GREEN)✓$(NC)" || echo "$(YELLOW)✗ Optional$(NC)"
	@echo -n "ModemManager:    "
	@which mmcli > /dev/null && echo "$(GREEN)✓$(NC)" || echo "$(YELLOW)✗ Optional$(NC)"
	@echo -n "golangci-lint:   "
	@which golangci-lint > /dev/null && echo "$(GREEN)✓$(NC)" || echo "$(YELLOW)✗ Optional$(NC)"

quick: ## Quick build and test
	@$(MAKE) fmt
	@$(MAKE) build
	@$(MAKE) test
	@echo "$(GREEN)✓ Quick check complete$(NC)"

.DEFAULT_GOAL := help
