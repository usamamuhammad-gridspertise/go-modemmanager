# CLI Transformation Guide

## Overview

The `go-modemmanager` library has been transformed into a powerful command-line tool called `mmctl`. This guide explains the transformation, architecture, and how to use and extend the CLI.

## Table of Contents

- [What Changed](#what-changed)
- [Architecture](#architecture)
- [Building the CLI](#building-the-cli)
- [Using the CLI](#using-the-cli)
- [Command Reference](#command-reference)
- [Extending the CLI](#extending-the-cli)
- [Comparison with Library](#comparison-with-library)

---

## What Changed

### Before: Library Only

```go
// Using the library in your Go code
package main

import "github.com/maltegrosse/go-modemmanager"

func main() {
    mm, _ := modemmanager.NewModemManager()
    modems, _ := mm.GetModems()
    // ... more code
}
```

### After: CLI + Library

```bash
# Now you can use the CLI directly
mmctl list
mmctl modem info -m 0
mmctl connect -m 0 --apn internet
```

**The library is still available for Go developers**, but now there's also a standalone CLI tool!

---

## Architecture

### Directory Structure

```
go-modemmanager/
├── cmd/                        # CLI applications
│   └── mmctl/                  # Main CLI tool
│       ├── main.go             # Entry point
│       ├── cmd/                # Command implementations
│       │   ├── root.go         # Root command & global flags
│       │   ├── list.go         # List modems command
│       │   ├── modem.go        # Modem management commands
│       │   ├── connect.go      # Connection commands
│       │   └── sms.go          # SMS commands
│       └── README.md           # CLI documentation
├── mocks/                      # Mock implementations for testing
│   ├── mock_modem.go           # Modem mocks
│   └── example_test.go         # Example tests
├── test-environment/           # Testing infrastructure
│   ├── mock-dbus/              # D-Bus mock service
│   └── ...
├── *.go                        # Original library files
├── go.mod                      # Module definition
├── Makefile                    # Build automation
└── CLI_GUIDE.md               # This file
```

### Technology Stack

- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - Modern CLI library
- **Output Formatting**: `text/tabwriter` for tables, `encoding/json` for JSON
- **Library**: Original `go-modemmanager` D-Bus bindings
- **Build Tool**: Make + Go build
- **Testing**: Go test + Mock implementations

---

## Building the CLI

### Quick Build

```bash
# Simple build
go build -o mmctl ./cmd/mmctl

# Using Make (recommended)
make build
```

### Full Build with Make

```bash
# Show all available commands
make help

# Build
make build              # Build binary
make build-all          # Build for multiple platforms (Linux amd64, arm64, armv7)

# Install
make install            # Install to /usr/local/bin (requires sudo)
make uninstall          # Remove from system

# Test
make test               # Run tests
make test-coverage      # Generate coverage report
make test-mock          # Test with mock ModemManager

# Code quality
make fmt                # Format code
make lint               # Run linter
make verify             # Lint + test

# Docker
make docker-build       # Build test environment
make docker-test        # Run tests in Docker
make docker-shell       # Interactive Docker shell

# Development
make quick              # Format + build + test
make watch              # Auto-rebuild on changes (requires entr)
make demo               # Run demo

# Release
make release            # Full release build (all platforms)
make clean              # Clean build artifacts
```

### Manual Build

```bash
# Basic build
go build -o mmctl ./cmd/mmctl

# Build with version info
go build -ldflags "-X main.version=0.1.0" -o mmctl ./cmd/mmctl

# Build for specific platform
GOOS=linux GOARCH=arm64 go build -o mmctl-arm64 ./cmd/mmctl

# Build for Raspberry Pi
GOOS=linux GOARCH=arm GOARM=7 go build -o mmctl-pi ./cmd/mmctl
```

---

## Using the CLI

### Installation

#### Option 1: Build and Install

```bash
git clone https://github.com/maltegrosse/go-modemmanager.git
cd go-modemmanager
make install
```

#### Option 2: Go Install

```bash
go install github.com/maltegrosse/go-modemmanager/cmd/mmctl@latest
```

#### Option 3: Download Binary

```bash
# Download from releases (when available)
curl -L https://github.com/maltegrosse/go-modemmanager/releases/download/v0.1.0/mmctl-linux-amd64 -o mmctl
chmod +x mmctl
sudo mv mmctl /usr/local/bin/
```

### Basic Usage

```bash
# Get help
mmctl --help
mmctl <command> --help

# List modems
mmctl list

# Get modem info
mmctl modem info -m 0

# Enable modem
mmctl modem enable -m 0

# Connect
mmctl connect -m 0 --apn internet

# Send SMS
mmctl sms send -m 0 --number +1234567890 --text "Hello"

# Check status
mmctl status -m 0
```

### Output Formats

#### Human-Readable (Default)

```bash
mmctl list
```

Output:
```
INDEX  MANUFACTURER  MODEL    STATE       SIGNAL  IMEI             PORT
-----  ------------  -----    -----       ------  ----             ----
0      Quectel       EC25     Registered  85%     123456789012345  ttyUSB0
```

#### JSON Format

```bash
mmctl list --json
```

Output:
```json
[
  {
    "index": 0,
    "manufacturer": "Quectel",
    "model": "EC25",
    "state": "MmModemStateRegistered",
    "signal_quality": 85
  }
]
```

### Scripting Examples

#### Monitor Signal Quality

```bash
#!/bin/bash
while true; do
    signal=$(mmctl modem signal -m 0 --json | jq -r '.quality')
    echo "$(date): Signal: ${signal}%"
    sleep 60
done
```

#### Auto-Reconnect

```bash
#!/bin/bash
while true; do
    connected=$(mmctl status -m 0 --json | jq -r '.connected')
    if [ "$connected" != "true" ]; then
        echo "Reconnecting..."
        mmctl connect -m 0 --apn internet
    fi
    sleep 300
done
```

---

## Command Reference

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--modem` | `-m` | Modem index (0, 1, 2, ...) |
| `--path` | `-p` | Modem D-Bus path |
| `--json` | `-j` | Output in JSON format |
| `--verbose` | `-v` | Verbose output |
| `--help` | `-h` | Show help |

### Commands

#### List Modems

```bash
mmctl list [flags]
```

Lists all detected modems with basic information.

#### Modem Commands

```bash
mmctl modem info -m <index>      # Detailed info
mmctl modem enable -m <index>    # Enable modem
mmctl modem disable -m <index>   # Disable modem
mmctl modem reset -m <index>     # Reset modem
mmctl modem signal -m <index>    # Signal quality
mmctl modem command -m <index> "AT_CMD"  # AT command
```

#### Connection Commands

```bash
mmctl connect -m <index> --apn <apn> [flags]
mmctl disconnect -m <index>
mmctl status -m <index>
```

Connect flags:
- `--apn` - Access Point Name (required)
- `--user` - Username
- `--password` - Password
- `--ip-type` - IP type (ipv4/ipv6/ipv4v6)
- `--allow-roaming` - Allow roaming

#### SMS Commands

```bash
mmctl sms send -m <index> --number <phone> --text <message>
mmctl sms list -m <index>
mmctl sms read -m <index> --sms-index <idx>
mmctl sms delete -m <index> --sms-index <idx>
```

---

## Extending the CLI

### Adding a New Command

1. **Create command file** in `cmd/mmctl/cmd/`:

```go
// cmd/mmctl/cmd/location.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var locationCmd = &cobra.Command{
    Use:   "location",
    Short: "Get modem location",
    Long:  "Get GPS location from modem.",
    RunE:  runLocation,
}

func init() {
    rootCmd.AddCommand(locationCmd)
}

func runLocation(cmd *cobra.Command, args []string) error {
    modem, err := getModem()
    if err != nil {
        return err
    }

    location, err := modem.GetLocation()
    if err != nil {
        return fmt.Errorf("failed to get location: %w", err)
    }

    // Get location data
    locations, err := location.GetLocation()
    if err != nil {
        return fmt.Errorf("failed to get location data: %w", err)
    }

    // Output location
    fmt.Printf("Latitude:  %.6f\n", locations.GpsRawLatitude)
    fmt.Printf("Longitude: %.6f\n", locations.GpsRawLongitude)

    return nil
}
```

2. **Test the command**:

```bash
make build
./mmctl location -m 0
```

3. **Add tests**:

```go
// cmd/mmctl/cmd/location_test.go
package cmd

import (
    "testing"
)

func TestLocationCommand(t *testing.T) {
    // Test implementation
}
```

### Adding Flags

```go
var (
    locationFormat string
)

func init() {
    rootCmd.AddCommand(locationCmd)
    
    // Add command-specific flags
    locationCmd.Flags().StringVar(&locationFormat, "format", "decimal", "Output format (decimal, dms)")
}
```

### Adding Subcommands

```go
var (
    locationCmd = &cobra.Command{
        Use:   "location",
        Short: "Location management",
    }

    locationGetCmd = &cobra.Command{
        Use:   "get",
        Short: "Get current location",
        RunE:  runLocationGet,
    }

    locationSetupCmd = &cobra.Command{
        Use:   "setup",
        Short: "Setup location sources",
        RunE:  runLocationSetup,
    }
)

func init() {
    rootCmd.AddCommand(locationCmd)
    locationCmd.AddCommand(locationGetCmd)
    locationCmd.AddCommand(locationSetupCmd)
}
```

Usage:
```bash
mmctl location get -m 0
mmctl location setup -m 0 --sources gps,agps
```

### Custom Output Formatting

```go
func outputLocationJSON(loc Location) error {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    return encoder.Encode(loc)
}

func outputLocationTable(loc Location) error {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    defer w.Flush()

    fmt.Fprintf(w, "Property\tValue\n")
    fmt.Fprintf(w, "--------\t-----\n")
    fmt.Fprintf(w, "Latitude\t%.6f\n", loc.Latitude)
    fmt.Fprintf(w, "Longitude\t%.6f\n", loc.Longitude)
    
    return nil
}

func runLocationGet(cmd *cobra.Command, args []string) error {
    // ... get location ...
    
    if jsonOutput {
        return outputLocationJSON(location)
    }
    return outputLocationTable(location)
}
```

---

## Comparison with Library

### Library Usage (Go Code)

```go
package main

import (
    "fmt"
    "github.com/maltegrosse/go-modemmanager"
)

func main() {
    // Connect to ModemManager
    mm, err := modemmanager.NewModemManager()
    if err != nil {
        panic(err)
    }

    // Get modems
    modems, err := mm.GetModems()
    if err != nil {
        panic(err)
    }

    // Use first modem
    modem := modems[0]

    // Enable modem
    if err := modem.Enable(true); err != nil {
        panic(err)
    }

    // Get signal quality
    signal, err := modem.GetSignalQuality()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Signal: %d%%\n", signal.Quality)

    // Create bearer
    simple, _ := modem.GetSimpleModem()
    bearer, err := simple.Connect(modemmanager.SimpleProperty{
        Apn: "internet",
    })
    if err != nil {
        panic(err)
    }

    // Connect
    if err := bearer.Connect(); err != nil {
        panic(err)
    }

    fmt.Println("Connected!")
}
```

### CLI Usage (Command Line)

```bash
# Enable modem
mmctl modem enable -m 0

# Get signal quality
mmctl modem signal -m 0

# Connect
mmctl connect -m 0 --apn internet
```

### When to Use Each

#### Use the Library When:
- Building Go applications
- Need fine-grained control
- Integrating with other Go code
- Building custom tools
- Need to handle D-Bus events
- Complex logic required

#### Use the CLI When:
- Interactive terminal usage
- Shell scripting
- Quick operations
- System administration
- Testing/debugging
- No coding required
- CI/CD pipelines

---

## Best Practices

### Error Handling

```go
func runMyCommand(cmd *cobra.Command, args []string) error {
    modem, err := getModem()
    if err != nil {
        return err  // Cobra handles the error display
    }

    // Use fmt.Errorf for context
    if err := modem.Enable(true); err != nil {
        return fmt.Errorf("failed to enable modem: %w", err)
    }

    return nil
}
```

### Verbose Output

```go
if verbose {
    fmt.Printf("Connecting to %s...\n", apn)
}

// Do operation

if verbose {
    fmt.Println("Connection established")
}
```

### JSON Output Support

```go
// Always support both output formats
if jsonOutput {
    return outputJSON(data)
}
return outputTable(data)
```

### Flag Validation

```go
func init() {
    myCmd.Flags().StringVar(&myFlag, "my-flag", "", "Description")
    myCmd.MarkFlagRequired("my-flag")
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Additional validation
    if myFlag == "" || len(myFlag) > 100 {
        return fmt.Errorf("invalid flag value: %s", myFlag)
    }
    // ...
}
```

---

## Testing the CLI

### Unit Tests

```go
func TestMyCommand(t *testing.T) {
    // Use mocks
    mockModem := mocks.NewMockModem()
    mockModem.SignalQualityValue = modemmanager.SignalQuality{
        Quality: 85,
        Recent:  true,
    }

    // Test your command logic
}
```

### Integration Tests

```bash
# Start mock ModemManager
cd test-environment/mock-dbus
sudo ./start-mock.sh &

# Run CLI tests
make build
./mmctl list
./mmctl modem info -m 0
```

### Docker Tests

```bash
# Build test environment
make docker-build

# Run tests in container
make docker-test

# Interactive testing
make docker-shell
# Inside container:
mmctl list
```

---

## Troubleshooting

### Build Issues

**Problem**: `cannot find package`

```bash
# Solution: Update dependencies
go mod tidy
go mod download
```

**Problem**: `undefined: cobra`

```bash
# Solution: Install cobra
go get github.com/spf13/cobra
```

### Runtime Issues

**Problem**: `Permission denied`

```bash
# Solution: Add user to dialout group
sudo usermod -a -G dialout $USER
# Or run with sudo
sudo mmctl list
```

**Problem**: `No modems found`

```bash
# Check ModemManager status
systemctl status ModemManager

# Check USB devices
lsusb

# Scan for modems
mmcli -S
```

---

## Future Enhancements

### Planned Features

- [ ] Voice call support (`mmctl call`)
- [ ] USSD support (`mmctl ussd`)
- [ ] GPS location tracking (`mmctl location`)
- [ ] Firmware updates (`mmctl firmware`)
- [ ] Interactive mode (`mmctl -i`)
- [ ] Configuration file support (`~/.mmctl.yaml`)
- [ ] Bash/Zsh completion
- [ ] Signal history graphs
- [ ] Daemon mode for monitoring

### Community Contributions

We welcome contributions! To add a feature:

1. Fork the repository
2. Create a feature branch
3. Implement the command in `cmd/mmctl/cmd/`
4. Add tests
5. Update documentation
6. Submit pull request

---

## Resources

- **Library Documentation**: https://pkg.go.dev/github.com/maltegrosse/go-modemmanager
- **Cobra Documentation**: https://github.com/spf13/cobra
- **ModemManager API**: https://www.freedesktop.org/software/ModemManager/api/latest/
- **D-Bus Specification**: https://dbus.freedesktop.org/doc/dbus-specification.html

---

## Summary

The CLI transformation adds:

✅ **User-friendly interface** - Easy to use without coding
✅ **Multiple output formats** - Human-readable and JSON
✅ **Shell scripting support** - Easy automation
✅ **Comprehensive commands** - Cover common use cases
✅ **Good documentation** - Extensive help and examples
✅ **Testing infrastructure** - Mocks and Docker support
✅ **Build automation** - Makefile for easy building

The library remains intact and fully functional for Go developers, while the CLI provides a convenient tool for system administrators, scripters, and interactive use.

**Happy modem managing! 📱**