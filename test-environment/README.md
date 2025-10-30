# ModemManager Test Environment for macOS

This directory contains a Docker-based testing environment for `go-modemmanager` development on macOS (or any system without direct ModemManager support).

## Overview

Since ModemManager is Linux-specific and requires D-Bus, this containerized environment provides:
- **Ubuntu 22.04** base with ModemManager installed
- **D-Bus system bus** properly configured
- **Go 1.21** for building and testing
- **USB device support** (with limitations on macOS)
- **Development workspace** mounted from host

## Quick Start

### Prerequisites

1. **Docker Desktop** for Mac (or Podman)
   ```bash
   brew install --cask docker
   ```

2. **Optional**: USB modem hardware (Quectel, Huawei, Sierra Wireless, etc.)

### Build and Run

```bash
# Navigate to test environment directory
cd test-environment

# Build the container
docker-compose build

# Start the container with interactive shell
docker-compose run modemmanager-test
```

Inside the container:
```bash
# Check ModemManager status
mmcli --version
mmcli -L

# Test the library
cd /workspace
go test -v

# Run examples
cd /workspace/examples
go run check_interfaces.go
```

## Testing Scenarios

### Scenario 1: Code Compilation & Unit Tests (No Hardware)

**Purpose**: Validate code compiles and basic unit tests pass

```bash
docker-compose run modemmanager-test bash -c "cd /workspace && go build ./..."
```

**What works**:
- ✅ Code compilation
- ✅ Import validation
- ✅ Type checking
- ✅ Mock-based unit tests

**Limitations**:
- ❌ No actual modem detection
- ❌ D-Bus calls will fail without hardware

### Scenario 2: With USB Modem (Limited on macOS)

**Challenge**: Docker on macOS runs in a VM, making USB passthrough difficult.

**Workaround Options**:

#### Option A: USB/IP Network Forwarding
```bash
# On Linux host with modem:
sudo apt-get install usbip
sudo modprobe usbip_host
usbipd -D
usbip list -l
usbip bind -b <busid>

# On Mac in container:
# Connect via USB/IP protocol
```

#### Option B: Serial Device Forwarding
If your modem creates serial devices (`/dev/ttyUSB*`):
```bash
# Use socat to forward serial port
socat TCP-LISTEN:54321,reuseaddr /dev/ttyUSB0
```

### Scenario 3: Integration Testing with VM

**Recommended for serious development**:

1. **Install UTM or VirtualBox**:
   ```bash
   brew install --cask utm
   ```

2. **Create Ubuntu VM**:
   - Download Ubuntu Server 22.04
   - Allocate 2GB RAM, 20GB disk
   - Enable USB passthrough in settings

3. **Install ModemManager in VM**:
   ```bash
   sudo apt update
   sudo apt install modemmanager libmm-glib-dev
   ```

4. **Mount project directory**:
   - Use shared folders or SSHFS
   - Clone go-modemmanager repo in VM

5. **Connect USB modem to VM**:
   - In UTM: Devices → USB → Connect modem
   - Verify: `mmcli -L`

## Directory Structure

```
test-environment/
├── Dockerfile              # Container definition
├── docker-compose.yml      # Service orchestration
├── README.md              # This file
└── mock-scripts/          # Mock D-Bus responses (future)
    └── mock-modem.py      # Python-based modem simulator
```

## Testing Without Hardware

### Option 1: Mock D-Bus Interface

Create mock responses for testing (TODO):

```python
# mock-scripts/mock-modem.py
# Simulates ModemManager D-Bus interface
```

### Option 2: Unit Tests with Mocked DBus

```go
// Create interface mocks for testing
type MockModemManager struct {
    mock.Mock
}

func (m *MockModemManager) GetModems() ([]Modem, error) {
    args := m.Called()
    return args.Get(0).([]Modem), args.Error(1)
}
```

## Troubleshooting

### Issue: D-Bus connection fails

```bash
# Check D-Bus is running
ps aux | grep dbus

# Restart D-Bus
rm -f /var/run/dbus/pid
dbus-daemon --system --fork
```

### Issue: ModemManager not starting

```bash
# Check logs
journalctl -u ModemManager -f

# Or run in foreground for debugging
/usr/sbin/ModemManager --debug
```

### Issue: No modems detected

```bash
# List USB devices
lsusb

# Check if modem is in correct mode
usb_modeswitch -l

# Check kernel modules
lsmod | grep -E 'qmi|cdc|option'

# Manually trigger scan
mmcli -S
```

### Issue: Permission denied

```bash
# Container needs privileged mode for device access
docker-compose run --privileged modemmanager-test
```

## Development Workflow

### 1. Make Code Changes on Host (Mac)

Edit files in your favorite IDE on macOS. Changes are automatically reflected in container via volume mount.

### 2. Test in Container

```bash
# Start container
docker-compose run modemmanager-test

# Inside container - quick test
go build ./...
go test -v ./...

# Run specific example
cd examples
go run check_interfaces.go
```

### 3. Iterate

The `/workspace` directory in the container is mounted from the parent directory, so all changes are live.

## Real Hardware Testing Recommendations

For comprehensive testing with actual modems:

### Budget Option: Raspberry Pi Setup

**Cost**: ~$50-80

1. Raspberry Pi 4 (2GB) - $35
2. USB 4G/LTE modem - $20-40
   - Quectel EC25
   - Huawei E3372
   - Sierra Wireless MC7455

**Setup**:
```bash
# On Raspberry Pi
sudo apt install modemmanager libmm-glib-dev golang
git clone <your-fork>
cd go-modemmanager
go test -v
```

### Remote Testing Option

**Use GitHub Actions with self-hosted runner**:

1. Set up Raspberry Pi or Linux machine as runner
2. Connect USB modem
3. Run tests on every push

```yaml
# .github/workflows/hardware-test.yml
name: Hardware Tests
on: [push]
jobs:
  test-with-modem:
    runs-on: [self-hosted, linux, arm64, modem]
    steps:
      - uses: actions/checkout@v3
      - name: Test with real modem
        run: |
          mmcli -L
          go test -v ./...
```

## Known Limitations on macOS/Docker

1. **USB Device Access**: Docker Desktop on Mac has limited USB support
2. **System D-Bus**: Container D-Bus is isolated from host
3. **Kernel Modules**: Container uses host kernel, limited modem driver availability
4. **Performance**: Running in VM layer adds overhead

## Alternatives to Docker

### Lima VM (Lightweight)

```bash
brew install lima
limactl start --name=modemmanager ubuntu-lts
lima sudo apt install modemmanager
```

### Podman (Docker alternative)

```bash
brew install podman
podman machine init --rootful
podman-compose up
```

## Next Steps

1. **Create comprehensive unit tests** with mocked D-Bus
2. **Build modem simulator** for hardware-free testing
3. **Set up CI/CD** with hardware testing
4. **Document tested hardware** compatibility matrix

## Resources

- [ModemManager Documentation](https://www.freedesktop.org/software/ModemManager/)
- [D-Bus Specification](https://dbus.freedesktop.org/doc/dbus-specification.html)
- [USB Modem Mode Switching](https://www.draisberghof.de/usb_modeswitch/)

## Support

For issues specific to this test environment, please check existing issues or create a new one describing:
- Your host OS (macOS version)
- Docker version
- Modem hardware (if any)
- Error messages and logs