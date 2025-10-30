# Mocking Guide for go-modemmanager

This guide explains how to mock ModemManager for testing without physical hardware.

## Table of Contents

- [Overview](#overview)
- [Approach 1: D-Bus Mock Service](#approach-1-dbus-mock-service)
- [Approach 2: Go Interface Mocks](#approach-2-go-interface-mocks)
- [Approach 3: Test Fixtures](#approach-3-test-fixtures)
- [Which Approach to Use?](#which-approach-to-use)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## Overview

There are three main approaches to mock ModemManager:

1. **D-Bus Mock Service** - Most realistic, simulates actual D-Bus interface
2. **Go Interface Mocks** - Best for unit tests, pure Go implementation
3. **Test Fixtures** - Simplest, uses recorded data

## Approach 1: D-Bus Mock Service

### Description

A Python-based D-Bus service that mimics ModemManager's D-Bus interface. This allows testing the entire stack including D-Bus communication.

### Advantages

✅ Most realistic - tests actual D-Bus communication
✅ No code changes required in go-modemmanager
✅ Can be used with mmcli and other tools
✅ Tests signal/event handling

### Disadvantages

❌ More complex setup
❌ Requires Python dependencies
❌ Slower than pure Go mocks
❌ Needs D-Bus daemon running

### Setup

#### Prerequisites

```bash
# In container or Linux system
apt-get install python3 python3-dbus python3-gi
```

#### Starting the Mock Service

```bash
cd test-environment/mock-dbus

# Stop real ModemManager first
sudo systemctl stop ModemManager

# Start mock service
sudo ./start-mock.sh
```

Or manually:

```bash
cd test-environment/mock-dbus
sudo python3 mock_modemmanager.py
```

### Testing with Mock Service

```bash
# Verify it's running
mmcli -L
# Output: /org/freedesktop/ModemManager1/Modem/0 [MockModem Inc.] MockModem X1000

# Get modem details
mmcli -m 0

# Now run Go tests
cd /workspace
go test -v
```

### Customizing Mock Behavior

Edit `mock_modemmanager.py` to change:

```python
# In MockModem.__init__():
self.manufacturer = "MyCustom Manufacturer"
self.model = "MyCustom Model"
self.state = MM_MODEM_STATE_REGISTERED
```

### Features Supported

The mock D-Bus service supports:

- ✅ ModemManager interface (version, scan, list modems)
- ✅ Modem interface (enable, disable, properties)
- ✅ Modem.Simple interface (connect, disconnect, status)
- ✅ Modem.Modem3gpp interface (register, scan, IMEI)
- ✅ Bearer interface (connect, disconnect, IP config)
- ✅ SIM interface (PIN operations, IMSI, ICCID)
- ✅ StateChanged signals
- ✅ PropertiesChanged signals

### Limitations

- ⚠️ No SMS/Messaging support yet
- ⚠️ No Voice/Call support yet
- ⚠️ No OMA support
- ⚠️ Simplified location support

---

## Approach 2: Go Interface Mocks

### Description

Pure Go mock implementations of all go-modemmanager interfaces. Best for unit testing individual components.

### Advantages

✅ Fast - no D-Bus overhead
✅ Easy to customize per test
✅ No external dependencies
✅ Works on any platform (including macOS)
✅ Fine-grained control over behavior

### Disadvantages

❌ Doesn't test D-Bus communication
❌ Requires code awareness of mocks
❌ Need to update when interfaces change

### Usage

#### Basic Example

```go
package mypackage_test

import (
    "testing"
    "github.com/maltegrosse/go-modemmanager/mocks"
)

func TestMyFunction(t *testing.T) {
    // Create mock ModemManager
    mockMM := mocks.NewMockModemManager()

    // Use it like the real thing
    version, err := mockMM.GetVersion()
    if err != nil {
        t.Fatalf("GetVersion failed: %v", err)
    }

    if version != "1.12.8-mock" {
        t.Errorf("Unexpected version: %s", version)
    }
}
```

#### Customizing Return Values

```go
func TestCustomBehavior(t *testing.T) {
    mockModem := mocks.NewMockModem()

    // Customize properties
    mockModem.ManufacturerValue = "Test Manufacturer"
    mockModem.ModelValue = "Test Model"
    mockModem.StateValue = mm.MmModemStateEnabled

    // Verify
    manufacturer, _ := mockModem.GetManufacturer()
    if manufacturer != "Test Manufacturer" {
        t.Errorf("Expected 'Test Manufacturer', got '%s'", manufacturer)
    }
}
```

#### Testing Error Cases

```go
func TestErrorHandling(t *testing.T) {
    mockModem := mocks.NewMockModem()

    // Configure to return error
    mockModem.EnableError = errors.New("simulated failure")

    // Test error handling
    err := mockModem.Enable(true)
    if err == nil {
        t.Fatal("Expected error, got nil")
    }
}
```

#### Integration Test Example

```go
func TestCompleteWorkflow(t *testing.T) {
    // Setup
    mockMM := mocks.NewMockModemManager()

    // Get modems
    modems, err := mockMM.GetModems()
    require.NoError(t, err)
    require.Len(t, modems, 1)

    modem := modems[0]

    // Enable modem
    err = modem.Enable(true)
    require.NoError(t, err)

    // Get 3GPP interface
    modem3gpp, err := modem.Get3gpp()
    require.NoError(t, err)

    // Check operator
    opName, err := modem3gpp.GetOperatorName()
    require.NoError(t, err)
    assert.Equal(t, "T-Mobile", opName)

    // Create bearer
    bearer, err := modem.CreateBearer(mm.BearerProperty{
        Apn: "internet",
    })
    require.NoError(t, err)

    // Connect
    err = bearer.Connect()
    require.NoError(t, err)

    // Verify connected
    connected, err := bearer.GetConnected()
    require.NoError(t, err)
    assert.True(t, connected)
}
```

### Available Mocks

The `mocks` package provides:

- `MockModemManager` - Main ModemManager interface
- `MockModem` - Modem interface
- `MockModemSimple` - Simple interface
- `MockModem3gpp` - 3GPP interface
- `MockBearer` - Bearer interface
- `MockSim` - SIM interface

More mocks can be added as needed.

### Creating Custom Mocks

```go
type MyCustomMock struct {
    // Add fields for state
    callCount int
    lastArg   string
}

func (m *MyCustomMock) SomeMethod(arg string) error {
    m.callCount++
    m.lastArg = arg
    return nil
}

func TestWithCustomMock(t *testing.T) {
    mock := &MyCustomMock{}

    // Use mock
    mock.SomeMethod("test")

    // Verify behavior
    assert.Equal(t, 1, mock.callCount)
    assert.Equal(t, "test", mock.lastArg)
}
```

---

## Approach 3: Test Fixtures

### Description

Use recorded D-Bus responses or pre-defined data structures for testing.

### Advantages

✅ Simplest approach
✅ Fast
✅ Good for regression testing
✅ Can use real modem data

### Disadvantages

❌ Static - no dynamic behavior
❌ Doesn't test actual operations
❌ Limited flexibility

### Creating Fixtures

#### Recording Real Data

```bash
# Capture D-Bus messages from real modem
dbus-monitor --system "type='signal',interface='org.freedesktop.ModemManager1.Modem'"
 > modem_signals.txt

# Capture mmcli output
mmcli -m 0 --