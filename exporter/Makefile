.PHONY: all build install clean test run help

# Binary name
BINARY_NAME=mm-exporter
BINARY_PATH=../cmd/mm-exporter

# Installation paths
INSTALL_PATH=/usr/local/bin
SERVICE_PATH=/etc/systemd/system

# Build variables
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Default target
all: build

# Build the exporter
build:
	@echo "Building $(BINARY_NAME)..."
	cd $(BINARY_PATH) && $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_PATH)/$(BINARY_NAME)"

# Build with debug symbols
build-debug:
	@echo "Building $(BINARY_NAME) with debug symbols..."
	cd $(BINARY_PATH) && $(GO) build $(GOFLAGS) -o $(BINARY_NAME)
	@echo "Debug build complete: $(BINARY_PATH)/$(BINARY_NAME)"

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BINARY_PATH)/$(BINARY_NAME) $(INSTALL_PATH)/
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete"
	@echo "Run 'make install-service' to install systemd service"

# Install systemd service
install-service:
	@echo "Installing systemd service..."
	@echo "[Unit]" | sudo tee $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "Description=ModemManager Prometheus Exporter" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "After=network.target ModemManager.service" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "Requires=ModemManager.service" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "[Service]" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "Type=simple" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "User=root" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "ExecStart=$(INSTALL_PATH)/$(BINARY_NAME)" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "Restart=on-failure" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "RestartSec=5s" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "[Install]" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	@echo "WantedBy=multi-user.target" | sudo tee -a $(SERVICE_PATH)/mm-exporter.service > /dev/null
	sudo systemctl daemon-reload
	@echo "Service installed. Enable with: sudo systemctl enable mm-exporter"
	@echo "Start with: sudo systemctl start mm-exporter"

# Uninstall the binary and service
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	-sudo systemctl stop mm-exporter 2>/dev/null
	-sudo systemctl disable mm-exporter 2>/dev/null
	-sudo rm -f $(SERVICE_PATH)/mm-exporter.service
	-sudo systemctl daemon-reload
	-sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Run the exporter locally
run: build
	@echo "Running $(BINARY_NAME)..."
	cd $(BINARY_PATH) && ./$(BINARY_NAME)

# Run with custom flags
run-custom: build
	@echo "Running $(BINARY_NAME) with custom settings..."
	cd $(BINARY_PATH) && ./$(BINARY_NAME) -listen-address=":9090" -signal-rate=10s

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_PATH)/$(BINARY_NAME)
	$(GO) clean -cache
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Check dependencies
deps:
	@echo "Checking dependencies..."
	$(GO) mod download
	$(GO) mod verify

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

# Show version
version: build
	cd $(BINARY_PATH) && ./$(BINARY_NAME) -version

# Check if ModemManager is running
check-modemmanager:
	@echo "Checking ModemManager status..."
	@systemctl is-active ModemManager && echo "✓ ModemManager is running" || echo "✗ ModemManager is not running"
	@echo ""
	@echo "Available modems:"
	@mmcli -L || echo "No modems found or mmcli not available"

# Test metrics endpoint
test-metrics:
	@echo "Testing metrics endpoint..."
	@curl -s http://localhost:9539/metrics | head -20
	@echo "..."
	@echo "(showing first 20 lines)"

# Show service status
status:
	@systemctl status mm-exporter

# Show service logs
logs:
	@journalctl -u mm-exporter -f

# Help
help:
	@echo "ModemManager Exporter - Makefile Help"
	@echo ""
	@echo "Available targets:"
	@echo "  make build              - Build the exporter binary"
	@echo "  make build-debug        - Build with debug symbols"
	@echo "  make install            - Install binary to $(INSTALL_PATH)"
	@echo "  make install-service    - Install systemd service"
	@echo "  make uninstall          - Remove binary and service"
	@echo "  make run                - Build and run locally"
	@echo "  make run-custom         - Run with custom flags"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Lint code (requires golangci-lint)"
	@echo "  make deps               - Download and verify dependencies"
	@echo "  make update-deps        - Update dependencies"
	@echo "  make version            - Show version"
	@echo "  make check-modemmanager - Check ModemManager status"
	@echo "  make test-metrics       - Test metrics endpoint (exporter must be running)"
	@echo "  make status             - Show service status"
	@echo "  make logs               - Follow service logs"
	@echo "  make help               - Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  1. make build"
	@echo "  2. make run"
	@echo "  3. curl http://localhost:9539/metrics"
	@echo ""
	@echo "Production deployment:"
	@echo "  1. make install"
	@echo "  2. make install-service"
	@echo "  3. sudo systemctl enable mm-exporter"
	@echo "  4. sudo systemctl start mm-exporter"
