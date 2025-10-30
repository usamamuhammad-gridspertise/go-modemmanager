# Makefile for go-modemmanager test environment
# Simplifies Docker container management and testing workflows

.PHONY: help build up down shell test clean rebuild logs check-modem install-deps vm-setup

# Default target
.DEFAULT_GOAL := help

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
RESET := \033[0m

help: ## Show this help message
	@echo "$(CYAN)go-modemmanager Test Environment$(RESET)"
	@echo ""
	@echo "$(GREEN)Available targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Quick Start:$(RESET)"
	@echo "  1. make build          # Build the container"
	@echo "  2. make shell          # Start interactive shell"
	@echo "  3. make test           # Run tests"

build: ## Build the Docker container
	@echo "$(CYAN)Building ModemManager test container...$(RESET)"
	docker-compose build modemmanager-test
	@echo "$(GREEN)✓ Build complete$(RESET)"

up: ## Start the container in background
	@echo "$(CYAN)Starting ModemManager container...$(RESET)"
	docker-compose up -d modemmanager-test
	@echo "$(GREEN)✓ Container started$(RESET)"
	@echo "Use 'make shell' to connect or 'make logs' to view output"

down: ## Stop and remove the container
	@echo "$(CYAN)Stopping ModemManager container...$(RESET)"
	docker-compose down
	@echo "$(GREEN)✓ Container stopped$(RESET)"

shell: ## Start interactive shell in container
	@echo "$(CYAN)Starting interactive shell...$(RESET)"
	docker-compose run --rm modemmanager-test /bin/bash

test: ## Run Go tests inside container
	@echo "$(CYAN)Running tests...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go test -v ./..."

test-build: ## Test if code compiles
	@echo "$(CYAN)Testing compilation...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go build ./..."
	@echo "$(GREEN)✓ Compilation successful$(RESET)"

test-examples: ## Run example programs
	@echo "$(CYAN)Running examples...$(RESET)"
	@echo "$(YELLOW)Checking interfaces example:$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace/examples && go run check_interfaces.go" || true
	@echo ""
	@echo "$(YELLOW)Note: Examples may fail without physical modem hardware$(RESET)"

check-modem: ## Check if modems are detected in container
	@echo "$(CYAN)Checking for modems...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "mmcli -L"

logs: ## Show container logs
	@echo "$(CYAN)Showing container logs...$(RESET)"
	docker-compose logs -f modemmanager-test

clean: ## Remove container and images
	@echo "$(CYAN)Cleaning up...$(RESET)"
	docker-compose down -v
	docker-compose rm -f
	@echo "$(GREEN)✓ Cleanup complete$(RESET)"

rebuild: clean build ## Clean and rebuild container
	@echo "$(GREEN)✓ Rebuild complete$(RESET)"

lint: ## Run Go linter
	@echo "$(CYAN)Running golangci-lint...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go fmt ./... && go vet ./..."

fmt: ## Format Go code
	@echo "$(CYAN)Formatting Go code...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go fmt ./..."

mod-tidy: ## Run go mod tidy
	@echo "$(CYAN)Tidying Go modules...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go mod tidy"

generate: ## Run go generate
	@echo "$(CYAN)Running go generate...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go generate ./..."

debug: ## Start container with debugging enabled
	@echo "$(CYAN)Starting container in debug mode...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "/usr/sbin/ModemManager --debug & bash"

dbus-status: ## Check D-Bus status in container
	@echo "$(CYAN)Checking D-Bus status...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "ps aux | grep dbus && dbus-send --system --print-reply --dest=org.freedesktop.DBus / org.freedesktop.DBus.ListNames" || true

mm-version: ## Show ModemManager version
	@echo "$(CYAN)ModemManager version:$(RESET)"
	docker-compose run --rm modemmanager-test mmcli --version

install-deps: ## Install development dependencies on macOS
	@echo "$(CYAN)Installing macOS dependencies...$(RESET)"
	@command -v docker >/dev/null 2>&1 || { echo "$(YELLOW)Installing Docker Desktop...$(RESET)"; brew install --cask docker; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(YELLOW)Installing docker-compose...$(RESET)"; brew install docker-compose; }
	@echo "$(GREEN)✓ Dependencies installed$(RESET)"
	@echo ""
	@echo "$(YELLOW)Note: Make sure Docker Desktop is running!$(RESET)"

vm-setup: ## Instructions for VM-based testing setup
	@echo "$(CYAN)VM-Based Testing Setup$(RESET)"
	@echo ""
	@echo "For full hardware testing, consider using a VM:"
	@echo ""
	@echo "$(GREEN)Option 1: UTM (Recommended for Mac)$(RESET)"
	@echo "  1. Install UTM: brew install --cask utm"
	@echo "  2. Download Ubuntu Server 22.04 ISO"
	@echo "  3. Create new VM (2GB RAM, 20GB disk)"
	@echo "  4. Enable USB passthrough in VM settings"
	@echo "  5. Install ModemManager: sudo apt install modemmanager"
	@echo ""
	@echo "$(GREEN)Option 2: VirtualBox$(RESET)"
	@echo "  1. Install VirtualBox: brew install --cask virtualbox"
	@echo "  2. Follow similar steps as UTM"
	@echo ""
	@echo "$(GREEN)Option 3: Remote Linux Machine$(RESET)"
	@echo "  1. Use Raspberry Pi or cloud Linux instance"
	@echo "  2. Connect via SSH for development"
	@echo ""
	@echo "See test-environment/README.md for detailed instructions"

quick-test: build test-build ## Quick build and compile test
	@echo "$(GREEN)✓ Quick test passed$(RESET)"

full-test: build test-build test ## Full test suite
	@echo "$(GREEN)✓ Full test passed$(RESET)"

# Advanced targets

mock-setup: ## Set up mock ModemManager for testing without hardware
	@echo "$(CYAN)Setting up mock ModemManager...$(RESET)"
	@echo "$(YELLOW)TODO: Implement mock D-Bus interface$(RESET)"
	@echo "This will allow testing without physical modem hardware"

benchmark: ## Run benchmarks
	@echo "$(CYAN)Running benchmarks...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go test -bench=. -benchmem ./..."

coverage: ## Generate test coverage report
	@echo "$(CYAN)Generating coverage report...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "cd /workspace && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"

usb-devices: ## List USB devices (requires privileged container)
	@echo "$(CYAN)Listing USB devices...$(RESET)"
	docker-compose run --rm --privileged modemmanager-test lsusb || echo "$(YELLOW)USB access may not work on macOS Docker$(RESET)"

watch-logs: ## Watch ModemManager logs in real-time
	@echo "$(CYAN)Watching ModemManager logs...$(RESET)"
	docker-compose run --rm modemmanager-test bash -c "tail -f /var/log/syslog | grep ModemManager" || docker-compose logs -f

interactive: build up ## Build, start, and attach to container
	@echo "$(CYAN)Starting interactive session...$(RESET)"
	docker-compose exec modemmanager-test /bin/bash

status: ## Show container status
	@echo "$(CYAN)Container status:$(RESET)"
	@docker-compose ps

# Info targets

info: ## Show environment information
	@echo "$(CYAN)Environment Information$(RESET)"
	@echo ""
	@echo "$(GREEN)Host System:$(RESET)"
	@uname -a
	@echo ""
	@echo "$(GREEN)Docker Version:$(RESET)"
	@docker --version || echo "$(RED)Docker not installed$(RESET)"
	@echo ""
	@echo "$(GREEN)Docker Compose Version:$(RESET)"
	@docker-compose --version || echo "$(RED)Docker Compose not installed$(RESET)"
	@echo ""
	@echo "$(GREEN)Container Status:$(RESET)"
	@docker-compose ps || echo "No containers running"

doctor: ## Check if environment is properly set up
	@echo "$(CYAN)Running environment checks...$(RESET)"
	@echo ""
	@echo -n "Docker installed: "
	@command -v docker >/dev/null 2>&1 && echo "$(GREEN)✓$(RESET)" || echo "$(RED)✗$(RESET)"
	@echo -n "Docker running: "
	@docker ps >/dev/null 2>&1 && echo "$(GREEN)✓$(RESET)" || echo "$(RED)✗$(RESET)"
	@echo -n "Docker Compose installed: "
	@command -v docker-compose >/dev/null 2>&1 && echo "$(GREEN)✓$(RESET)" || echo "$(RED)✗$(RESET)"
	@echo -n "Go project files present: "
	@test -f ../go.mod && echo "$(GREEN)✓$(RESET)" || echo "$(RED)✗$(RESET)"
	@echo ""
	@docker ps >/dev/null 2>&1 || echo "$(YELLOW)⚠ Docker is not running. Start Docker Desktop first.$(RESET)"
