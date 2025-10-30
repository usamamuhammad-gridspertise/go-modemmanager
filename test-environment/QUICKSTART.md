# Quick Start Guide - Testing go-modemmanager on macOS

This is a **5-minute quick start** to get you testing `go-modemmanager` on your Mac.

## Prerequisites

- macOS (any recent version)
- Docker Desktop installed and running

## Step 1: Install Docker (if needed)

```bash
brew install --cask docker
```

Then launch Docker Desktop from Applications.

## Step 2: Navigate to Test Environment

```bash
cd test-environment
```

## Step 3: Build the Container

```bash
make build
```

This will:
- Pull Ubuntu 22.04
- Install ModemManager and dependencies
- Install Go 1.21
- Configure D-Bus

**Time**: ~2-3 minutes

## Step 4: Start Testing

### Quick Compilation Test (No Hardware Needed)

```bash
make test-build
```

This verifies the code compiles successfully.

### Run Unit Tests

```bash
make test
```

⚠️ **Note**: Tests requiring actual modem hardware will fail, but this validates the code structure.

### Interactive Shell

```bash
make shell
```

Inside the container:
```bash
# Check ModemManager version
mmcli --version

# Check for modems (will show none without hardware)
mmcli -L

# Build the library
cd /workspace
go build ./...

# Try running an example
cd examples
go run check_interfaces.go
```

## Step 5: Development Workflow

1. **Edit code** on your Mac using your favorite editor (VS Code, etc.)
2. **Test in container**: `make test-build
` or `make shell`
3. **Iterate**: Changes are live-synced via Docker volume

## Understanding Test Results

### ✅ Expected Successes (No Hardware)
- Code compilation
- Import validation
- Type checking
- Mock-based tests

### ❌ Expected Failures (No Hardware)
- Modem detection tests
- D-Bus communication with actual modems
- Hardware-specific feature tests

## What Next?

### Option A: Continue Without Hardware
Perfect for:
- Code refactoring
- Adding new interfaces
- Documentation improvements
- Type system updates

### Option B: Get Real Hardware Testing

**Recommended**: Set up a Linux VM with USB passthrough

1. Install UTM (free):
   ```bash
   brew install --cask utm
   ```

2. Create Ubuntu VM:
   - Download Ubuntu Server 22.04
   - Allocate 2GB RAM, 20GB disk
   - Enable USB controller in settings

3. Install ModemManager in VM:
   ```bash
   sudo apt update
   sudo apt install modemmanager golang-1.21
   ```

4. Connect USB modem to VM (in UTM: Devices → USB)

5. Clone and test:
   ```bash
   git clone <your-repo>
   cd go-modemmanager
   mmcli -L  # Should show modem!
   go test -v
   ```

**Alternative**: Use a Raspberry Pi (~$50) with a USB modem for dedicated testing.

## Common Commands

```bash
make help           # Show all available commands
make build          # Build container
make shell          # Interactive shell
make test           # Run tests
make test-build     # Compilation test
make check-modem    # Check for modems
make logs           # View logs
make clean          # Remove container
make doctor         # Check environment setup
```

## Troubleshooting

### "Docker not running"
→ Launch Docker Desktop application

### "Cannot connect to Docker daemon"
→ Make sure Docker Desktop is fully started (icon in menu bar)

### "make: command not found"
→ Use docker-compose directly:
```bash
docker-compose build
docker-compose run modemmanager-test
```

### "No modems detected"
→ Expected! You need:
- Physical USB modem, OR
- Linux VM with USB passthrough, OR
- Mock D-Bus interface (coming soon)

## Need Real Testing?

For comprehensive testing with actual modems, you have these options:

1. **Linux VM** (UTM/VirtualBox) - Best for Mac
2. **Raspberry Pi** (~$50) - Dedicated test hardware
3. **Cloud Linux** (AWS/DigitalOcean) - Remote testing
4. **Dual-boot Linux** - Native performance

See `README.md` for detailed setup instructions for each option.

## Support

- **Issues**: Check existing issues or create new one
- **Documentation**: See `test-environment/README.md` for detailed info
- **Examples**: Check `examples/` directory

## Summary

✅ **What you can do now (no hardware)**:
- Compile and test code structure
- Refactor and improve code
- Add new interfaces
- Update to newer Go versions
- Improve documentation

⚠️ **What you need hardware for**:
- Testing actual modem operations
- Validating D-Bus communication
- Hardware-specific features
- End-to-end integration tests

---

**Ready to upgrade the library?** Start with code improvements that don't require hardware:
1. Update Go version (1.13 → 1.21+)
2. Add comprehensive unit tests with mocks
3. Improve documentation
4. Modernize code patterns

Then set up hardware testing when needed!