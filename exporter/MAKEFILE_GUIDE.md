# Makefile Usage Guide for ModemManager Exporter

## Overview

The Makefile in the `exporter/` directory provides an easy way to build, install, run, and manage the ModemManager Prometheus exporter. This guide explains all available commands and common workflows.

## Location

```bash
cd ZedProjects/go-modemmanager/exporter
```

**Important**: Always run `make` commands from the `exporter/` directory, not from `cmd/mm-exporter/`.

## Quick Reference

```bash
make                    # Same as 'make build'
make help               # Show all available commands
make build              # Build the exporter
make run                # Build and run locally
make install            # Install system-wide
make install-service    # Install systemd service
make uninstall          # Remove everything
```

## Development Workflow

### 1. Initial Setup and Testing

```bash
# Navigate to exporter directory
cd exporter

# Build the binary
make build

# Run it locally (Ctrl+C to stop)
make run
```

Expected output:
```
Building mm-exporter...
Build complete: ../cmd/mm-exporter/mm-exporter
Running mm-exporter...
Starting ModemManager Exporter v1.0.0
Listening on :9539
```

### 2. Test the Exporter

While the exporter is running (in another terminal):

```bash
# Check if it's responding
curl http://localhost:9539/health

# View metrics
curl http://localhost:9539/metrics | head -20

# Or use the Makefile (requires exporter to be running)
make test-metrics
```

### 3. Check ModemManager Status

```bash
# Verify ModemManager is running and modems are detected
make check-modemmanager
```

Expected output:
```
Checking ModemManager status...
✓ ModemManager is running

Available modems:
    /org/freedesktop/ModemManager1/Modem/0 [Quectel] EC25
```

## Production Deployment

### Option 1: Manual Deployment (Recommended for Testing)

```bash
# 1. Build the binary
make build

# 2. Run with custom settings
cd ../cmd/mm-exporter
./mm-exporter -listen-address=":9090" -signal-rate=15s
```

### Option 2: System-Wide Installation (Recommended for Production)

```bash
# 1. Install binary to /usr/local/bin/
make install

# 2. Install systemd service
make install-service

# 3. Enable service to start on boot
sudo systemctl enable mm-exporter

# 4. Start the service
sudo systemctl start mm-exporter

# 5. Check if it's running
make status

# 6. View logs
make logs
```

## All Available Commands

### Building

#### `make build` (or just `make`)
Builds the exporter binary with optimizations (smaller size, stripped symbols).

```bash
make build
# Output: ../cmd/mm-exporter/mm-exporter
```

#### `make build-debug`
Builds with debug symbols (larger size, better for debugging).

```bash
make build-debug
```

#### `make clean`
Removes built binaries and cleans Go cache.

```bash
make clean
```

### Running

#### `make run`
Builds and runs the exporter locally with default settings.

```bash
make run
# Starts exporter on :9539
# Press Ctrl+C to stop
```

#### `make run-custom`
Runs with custom settings (port 9090, 10s signal rate).

```bash
make run-custom
# Edit Makefile to change the custom flags
```

To run with your own custom flags:
```bash
make build
cd ../cmd/mm-exporter
./mm-exporter -listen-address=":8080" -signal-rate=30s -metrics-path="/custom-metrics"
```

### Installation

#### `make install`
Installs the binary to `/usr/local/bin/mm-exporter` (requires sudo).

```bash
make install
# Binary is now available system-wide as 'mm-exporter'
```

After installation, you can run it from anywhere:
```bash
mm-exporter
mm-exporter -version
```

#### `make install-service`
Creates and installs a systemd service file (requires sudo).

```bash
make install-service
# Creates /etc/systemd/system/mm-exporter.service
```

After installing the service:
```bash
# Enable to start on boot
sudo systemctl enable mm-exporter

# Start now
sudo systemctl start mm-exporter

# Check status
sudo systemctl status mm-exporter
```

#### `make uninstall`
Removes both the binary and systemd service (requires sudo).

```bash
make uninstall
# Stops service, removes binary and service file
```

### Monitoring & Debugging

#### `make status`
Shows the systemd service status.

```bash
make status
```

Output shows:
- Service state (active/inactive/failed)
- Uptime
- Recent log entries
- Memory usage

#### `make logs`
Follows the service logs in real-time (like `tail -f`).

```bash
make logs
# Press Ctrl+C to stop following
```

#### `make test-metrics`
Tests the metrics endpoint (exporter must be running).

```bash
make test-metrics
# Shows first 20 lines of metrics output
```

#### `make check-modemmanager`
Checks if ModemManager is running and lists available modems.

```bash
make check-modemmanager
```

### Testing & Quality

#### `make test`
Runs Go tests (if any test files exist).

```bash
make test
```

#### `make test-coverage`
Runs tests and generates coverage report.

```bash
make test-coverage
# Creates coverage.html - open in browser
```

#### `make fmt`
Formats Go code according to standard style.

```bash
make fmt
```

#### `make lint`
Runs code linter (requires `golangci-lint` to be installed).

```bash
# Install golangci-lint first:
# https://golangci-lint.run/usage/install/

make lint
```

### Dependencies

#### `make deps`
Downloads and verifies Go dependencies.

```bash
make deps
```

#### `make update-deps`
Updates all dependencies to latest versions.

```bash
make update-deps
```

### Other

#### `make version`
Builds and shows the exporter version.

```bash
make version
# Output: mm-exporter version 1.0.0
```

#### `make help`
Shows help with all available commands.

```bash
make help
```

## Common Workflows

### Development Cycle

```bash
# 1. Make code changes
nano ../exporter/handler.go

# 2. Format code
make fmt

# 3. Build and test
make build
make run

# 4. In another terminal, test it
curl http://localhost:9539/metrics

# 5. Stop with Ctrl+C and repeat
```

### First Time Setup on a New Server

```bash
# 1. Ensure ModemManager is installed and running
sudo systemctl status ModemManager
mmcli -L

# 2. Clone/copy the project to the server

# 3. Navigate to exporter directory
cd go-modemmanager/exporter

# 4. Install everything
make install
make install-service

# 5. Start the service
sudo systemctl enable mm-exporter
sudo systemctl start mm-exporter

# 6. Verify it's working
make status
make test-metrics

# 7. Check logs if there are issues
make logs
```

### Updating the Exporter

```bash
# 1. Stop the service
sudo systemctl stop mm-exporter

# 2. Pull latest code / make changes

# 3. Rebuild and reinstall
make install

# 4. Restart the service
sudo systemctl start mm-exporter

# 5. Verify
make status
make logs
```

### Troubleshooting

```bash
# Check if binary exists
ls -lh ../cmd/mm-exporter/mm-exporter

# Check if ModemManager is running
make check-modemmanager

# Check service status
make status

# View detailed logs
make logs

# Test manually
make build
cd ../cmd/mm-exporter
./mm-exporter  # Run in foreground to see all output

# Check port availability
netstat -tlnp | grep 9539

# Test metrics endpoint
curl http://localhost:9539/health
curl http://localhost:9539/metrics
```

## Configuration

### Changing Default Settings

The Makefile uses these default locations:
- Binary: `/usr/local/bin/mm-exporter`
- Service file: `/etc/systemd/system/mm-exporter.service`

To change these, edit the Makefile:

```makefile
# At the top of the Makefile
INSTALL_PATH=/opt/bin              # Change installation directory
SERVICE_PATH=/etc/systemd/system   # Change service directory
```

### Customizing Service Settings

After running `make install-service`, you can edit the service file:

```bash
sudo nano /etc/systemd/system/mm-exporter.service
```

Change the `ExecStart` line to add flags:
```ini
ExecStart=/usr/local/bin/mm-exporter -listen-address=":9090" -signal-rate=15s
```

Then reload:
```bash
sudo systemctl daemon-reload
sudo systemctl restart mm-exporter
```

## Tips & Best Practices

### 1. **Always work from the exporter directory**
```bash
cd exporter/  # Correct
make build

# NOT from:
cd cmd/mm-exporter/  # Wrong
```

### 2. **Use `make run` for development**
- Quick iteration
- Easy to stop and restart
- See all output immediately

### 3. **Use `make install-service` for production**
- Starts automatically on boot
- Restarts on failure
- Managed by systemd
- Logs to journal

### 4. **Check logs when troubleshooting**
```bash
make logs              # Systemd logs
journalctl -xe         # All system logs
dmesg | grep -i modem  # Kernel messages
```

### 5. **Test before deploying**
```bash
make build
make run
# In another terminal:
make test-metrics
```

### 6. **Keep dependencies updated**
```bash
make update-deps
make test
make build
```

## Environment Requirements

The Makefile requires:
- **Go 1.13+**: For building
- **Make**: For running Makefile
- **ModemManager**: For the exporter to function
- **sudo**: For installation commands
- **systemctl**: For service management (Linux only)
- **curl**: For `make test-metrics`

Check requirements:
```bash
go version          # Should be 1.13 or higher
make --version      # Should be installed
mmcli --version     # ModemManager CLI
systemctl --version # Systemd (Linux)
```

## Platform Notes

### Linux (Ubuntu/Debian)
All commands work as documented.

### Linux (Other Distros)
May need to adjust service paths or use different init system.

### macOS
- Can build: ✅ (`make build`)
- Can run locally: ⚠️ (no ModemManager)
- Cannot install service: ❌ (no systemd)

Use for development/testing only.

### Docker/Container
Skip `install-service`, use Docker or run binary directly.

## Summary

| Use Case | Commands |
|----------|----------|
| Quick test | `make run` |
| Development | `make build` → edit → `make run` |
| Production install | `make install` → `make install-service` |
| Check status | `make status` or `make logs` |
| Update | Stop → `make install` → Start |
| Remove | `make uninstall` |
| Help | `make help` |

## Getting Help

```bash
# Show all commands
make help

# Check this guide
cat MAKEFILE_GUIDE.md

# Check main docs
cat README.md
cat QUICKSTART.md
```

For issues, check:
1. `make check-modemmanager` - Is ModemManager running?
2. `make status` - Is the exporter running?
3. `make logs` - What do the logs say?
4. `make test-metrics` - Can you reach the metrics?

---

**Remember**: Always run `make` commands from the `exporter/` directory! 📂