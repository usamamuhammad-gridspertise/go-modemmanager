# mmctl - ModemManager CLI Tool

A powerful command-line interface for managing cellular modems via ModemManager.

## Overview

`mmctl` provides a user-friendly way to interact with ModemManager and control cellular modems. It's built using the `go-modemmanager` library and offers an alternative to the standard `mmcli` tool with additional features and better usability.

## Features

- 📱 **Modem Management** - List, enable, disable, and reset modems
- 🌐 **Network Connectivity** - Connect and disconnect from mobile networks
- 📨 **SMS Messaging** - Send, receive, and manage text messages
- 📊 **Signal Quality** - Monitor signal strength and network status
- 💻 **AT Commands** - Send raw AT commands to modems
- 🔍 **Detailed Information** - Get comprehensive modem and SIM details
- 📋 **Multiple Output Formats** - Human-readable tables or JSON for scripting
- 🎯 **Simple Interface** - Intuitive commands with helpful examples

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/maltegrosse/go-modemmanager.git
cd go-modemmanager

# Build the CLI
go build -o mmctl ./cmd/mmctl

# Install to system (optional)
sudo mv mmctl /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/maltegrosse/go-modemmanager/cmd/mmctl@latest
```

## Quick Start

```bash
# List all modems
mmctl list

# Get detailed modem information
mmctl modem info -m 0

# Enable a modem
mmctl modem enable -m 0

# Connect to network
mmctl connect -m 0 --apn internet

# Send an SMS
mmctl sms send -m 0 --number +1234567890 --text "Hello World"

# Check connection status
mmctl status -m 0
```

## Commands

### Global Flags

Available for all commands:

- `-m, --modem <index>` - Modem index (0, 1, 2, etc.)
- `-p, --path <path>` - Modem D-Bus path (alternative to index)
- `-j, --json` - Output in JSON format
- `-v, --verbose` - Verbose output with additional details
- `--help` - Show help for any command

### List Modems

List all available modems detected by ModemManager.

```bash
mmctl list [flags]

# Examples:
mmctl list                    # List all modems
mmctl list --json             # JSON output
mmctl list --verbose          # Detailed output
```

**Output:**
```
INDEX  MANUFACTURER  MODEL           STATE       SIGNAL  IMEI             PORT
-----  ------------  -----           -----       ------  ----             ----
0      Quectel       EC25            Registered  85%     123456789012345  ttyUSB0
```

### Modem Commands

Manage and query modem devices.

#### Get Modem Information

```bash
mmctl modem info -m <index> [flags]

# Examples:
mmctl modem info -m 0
mmctl modem info -m 0 --json
mmctl modem info -m 0 --verbose
```

Shows comprehensive modem details including hardware info, capabilities, state, SIM details, and network information.

#### Enable/Disable Modem

```bash
# Enable modem
mmctl modem enable -m <index>

# Disable modem
mmctl modem disable -m <index>

# Examples:
mmctl modem enable -m 0
mmctl modem disable -m 0 --verbose
```

#### Reset Modem

```bash
mmctl modem reset -m <index>

# Example:
mmctl modem reset -m 0
```

Resets the modem to its initial state.

#### Get Signal Quality

```bash
mmctl modem signal -m <index> [flags]

# Examples:
mmctl modem signal -m 0
mmctl modem signal -m 0 --json
```

**Output:**
```
Signal Quality: 85% (recent)
Signal Bars:    [████░]
```

#### Send AT Command

```bash
mmctl modem command -m <index> "<AT_COMMAND>" [flags]

# Examples:
mmctl modem command -m 0 "ATI"
mmctl modem command -m 0 "AT+CSQ" --timeout 5
```

**Warning:** Sending incorrect AT commands can disrupt modem operation.

### Connection Commands

Manage mobile data connections.

#### Connect to Network

```bash
mmctl connect -m <index> --apn <apn> [flags]

# Flags:
#   --apn string         Access Point Name (required)
#   --user string        Username for authentication
#   --password string    Password for authentication
#   --ip-type string     IP type: ipv4, ipv6, ipv4v6 (default "ipv4")
#   --allow-roaming      Allow connection while roaming

# Examples:
mmctl connect -m 0 --apn internet
mmctl connect -m 0 --apn internet --user myuser --password mypass
mmctl connect -m 0 --apn internet --ip-type ipv4v6
mmctl connect -m 0 --apn internet --allow-roaming
```

Creates a data connection and displays IP configuration.

#### Disconnect from Network

```bash
mmctl disconnect -m <index>

# Example:
mmctl disconnect -m 0
```

Disconnects all active bearers.

#### Get Connection Status

```bash
mmctl status -m <index> [flags]

# Examples:
mmctl status -m 0
mmctl status -m 0 --json
mmctl status -m 0 --verbose
```

**Output:**
```
Connection Status
=================

State:          Connected
Signal:         85%
Registration:   Home
Operator:       T-Mobile

Data Connection:

Bearer 0:
  Status:       ✓ Connected
  Interface:    wwan0
  APN:          internet
  IPv4:         10.20.30.40/24
  Gateway:      10.20.30.1
  DNS:          [8.8.8.8 8.8.4.4]
  RX:           1024000 bytes
  TX:           512000 bytes
  Duration:     2h30m15s
```

### SMS Commands

Send, receive, and manage text messages.

#### Send SMS

```bash
mmctl sms send -m <index> --number <phone> --text <message> [flags]

# Flags:
#   --number string      Recipient phone number (required)
#   --text string        Message text (required)
#   --validity int       Message validity in minutes (0 = default)

# Examples:
mmctl sms send -m 0 --number +1234567890 --text "Hello World"
mmctl sms send -m 0 -n +1234567890 -t "Test message" --verbose
```

#### List SMS Messages

```bash
mmctl sms list -m <index> [flags]

# Examples:
mmctl sms list -m 0
mmctl sms list -m 0 --json
mmctl sms list -m 0 --verbose
```

**Output:**
```
INDEX  NUMBER          STATE     TIMESTAMP         MESSAGE
-----  ------          -----     ---------         -------
0      +1234567890     Received  2024-01-15 14:30  Hello, this is a test message
1      +0987654321     Sent      2024-01-15 15:45  Reply message
```

#### Read SMS Message

```bash
mmctl sms read -m <index> --sms-index <sms_index> [flags]

# Examples:
mmctl sms read -m 0 --sms-index 0
mmctl sms read -m 0 -i 0 --json
```

Displays full message details including sender, timestamp, and complete text.

#### Delete SMS Message

```bash
mmctl sms delete -m <index> --sms-index <sms_index>

# Example:
mmctl sms delete -m 0 --sms-index 0
```

### Help and Version

```bash
# Get help
mmctl --help
mmctl <command> --help
mmctl <command> <subcommand> --help

# Get version
mmctl --version
```

## Output Formats

### Human-Readable (Default)

Formatted tables and text output designed for terminal viewing:

```bash
mmctl list
```

### JSON Format

Machine-readable JSON output for scripting:

```bash
mmctl list --json
```

Example JSON output:
```json
[
  {
    "index": 0,
    "path": "/org/freedesktop/ModemManager1/Modem/0",
    "manufacturer": "Quectel",
    "model": "EC25",
    "state": "MmModemStateRegistered",
    "signal_quality": 85,
    "equipment_identifier": "123456789012345",
    "device": "usb:003:006",
    "primary_port": "ttyUSB0"
  }
]
```

## Usage Examples

### Basic Workflow

```bash
# 1. Check available modems
mmctl list

# 2. Get modem details
mmctl modem info -m 0

# 3. Enable modem
mmctl modem enable -m 0

# 4. Check signal quality
mmctl modem signal -m 0

# 5. Connect to internet
mmctl connect -m 0 --apn internet

# 6. Check connection status
mmctl status -m 0

# 7. Disconnect when done
mmctl disconnect -m 0
```

### Scripting Example

Monitor signal quality and log to file:

```bash
#!/bin/bash
while true; do
    signal=$(mmctl modem signal -m 0 --json | jq -r '.quality')
    echo "$(date): Signal quality: ${signal}%" >> signal.log
    sleep 60
done
```

### SMS Notification

Send SMS when an event occurs:

```bash
#!/bin/bash
if [ "$BACKUP_FAILED" = "true" ]; then
    mmctl sms send -m 0 \
        --number "+1234567890" \
        --text "ALERT: Backup failed on $(hostname)"
fi
```

### Connection Monitoring

Check if connected and reconnect if needed:

```bash
#!/bin/bash
status=$(mmctl status -m 0 --json | jq -r '.connected')
if [ "$status" != "true" ]; then
    echo "Not connected, attempting to reconnect..."
    mmctl connect -m 0 --apn internet
fi
```

## Troubleshooting

### No Modems Found

```bash
$ mmctl list
No modems found
```

**Solutions:**
- Check if ModemManager is running: `systemctl status ModemManager`
- Check if modem is detected: `lsusb`
- Try scanning: `mmcli -S` (with mmcli)
- Check modem is not inhibited

### Permission Denied

```bash
Error: failed to connect to ModemManager: access denied
```

**Solutions:**
- Run with sudo: `sudo mmctl list`
- Add user to dialout group: `sudo usermod -a -G dialout $USER`
- Check D-Bus permissions

### Modem Not Responding

```bash
Error: command failed: timeout
```

**Solutions:**
- Reset modem: `mmctl modem reset -m 0`
- Check modem state: `mmctl modem info -m 0`
- Power cycle modem
- Check USB connection

### Connection Fails

```bash
Error: failed to connect: invalid APN
```

**Solutions:**
- Verify APN with carrier: `mmctl connect -m 0 --apn correct-apn`
- Check SIM is registered: `mmctl modem info -m 0`
- Try with authentication: `--user <user> --password <pass>`
- Check roaming status: `--allow-roaming`

### SMS Send Fails

```bash
Error: failed to send SMS: not registered
```

**Solutions:**
- Check registration: `mmctl modem info -m 0`
- Check signal: `mmctl modem signal -m 0`
- Wait for registration
- Enable modem first: `mmctl modem enable -m 0`

## Comparison with mmcli

| Feature | mmctl | mmcli |
|---------|-------|-------|
| User-friendly commands | ✅ | ❌ |
| JSON output | ✅ | ✅ |
| Table formatting | ✅ | ❌ |
| Signal visualization | ✅ | ❌ |
| Simplified connection | ✅ | ❌ |
| AT commands | ✅ | ✅ |
| All MM features | 🚧 | ✅ |

**Note:** `mmctl` focuses on common use cases with better UX, while `mmcli` provides complete ModemManager API access.

## Requirements

- **OS:** Linux (ModemManager uses D-Bus)
- **ModemManager:** Version 1.10.0 or later
- **Go:** Version 1.21 or later (for building)
- **Permissions:** User must be in `dialout` group or run with sudo

## Configuration

`mmctl` uses the system D-Bus to communicate with ModemManager. No additional configuration is required.

## Development

### Building

```bash
go build -o mmctl ./cmd/mmctl
```

### Testing

```bash
# With mock
cd /workspace/test-environment/mock-dbus
./start-mock.sh

# In another terminal
go test ./cmd/mmctl/...
```

### Adding New Commands

1. Create new file in `cmd/mmctl/cmd/`
2. Define cobra command
3. Add to root command in `init()`
4. Implement command logic
5. Update this README

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Update documentation
5. Submit a pull request

## License

MIT License - see [LICENSE](../../LICENSE) file for details.

## Related Projects

- [go-modemmanager](https://github.com/maltegrosse/go-modemmanager) - The underlying Go library
- [ModemManager](https://www.freedesktop.org/wiki/Software/ModemManager/) - The modem management daemon
- [mmcli](https://www.freedesktop.org/software/ModemManager/) - Official ModemManager CLI

## Support

- **Issues:** [GitHub Issues](https://github.com/maltegrosse/go-modemmanager/issues)
- **Documentation:** [go-modemmanager docs](https://pkg.go.dev/github.com/maltegrosse/go-modemmanager)
- **ModemManager API:** [D-Bus API Reference](https://www.freedesktop.org/software/ModemManager/api/latest/)

## Changelog

### v0.1.0 (Initial Release)
- List modems with detailed information
- Modem management (enable, disable, reset)
- Signal quality monitoring
- Data connection management
- SMS send, receive, and management
- AT command support
- JSON output format
- Verbose mode

## Roadmap

- [ ] Voice call support
- [ ] USSD support
- [ ] GPS/location tracking
- [ ] Firmware update support
- [ ] OMA device management
- [ ] Interactive mode
- [ ] Configuration file support
- [ ] Signal history/graphs
- [ ] Multi-modem operations
- [ ] Bash completion

---

**Made with ❤️ using [go-modemmanager](https://github.com/maltegrosse/go-modemmanager)**