package mocks_test

import (
	"testing"

	"github.com/maltegrosse/go-modemmanager/mocks"
)

// TestMockModemManager demonstrates how to use the mock ModemManager
func TestMockModemManager(t *testing.T) {
	// Create a mock ModemManager
	mockMM := mocks.NewMockModemManager()

	// Test GetVersion
	version, err := mockMM.GetVersion()
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}
	if version != "1.12.8-mock" {
		t.Errorf("Expected version '1.12.8-mock', got '%s'", version)
	}

	// Test GetModems
	modems, err := mockMM.GetModems()
	if err != nil {
		t.Fatalf("GetModems failed: %v", err)
	}
	if len(modems) != 1 {
		t.Errorf("Expected 1 modem, got %d", len(modems))
	}

	// Test ScanDevices (should not error by default)
	err = mockMM.ScanDevices()
	if err != nil {
		t.Errorf("ScanDevices failed: %v", err)
	}
}

// TestMockModem demonstrates how to use the mock Modem
func TestMockModem(t *testing.T) {
	// Create a mock modem
	mockModem := mocks.NewMockModem()

	// Test getting manufacturer
	manufacturer, err := mockModem.GetManufacturer()
	if err != nil {
		t.Fatalf("GetManufacturer failed: %v", err)
	}
	if manufacturer != "MockModem Inc." {
		t.Errorf("Expected 'MockModem Inc.', got '%s'", manufacturer)
	}

	// Test getting model
	model, err := mockModem.GetModel()
	if err != nil {
		t.Fatalf("GetModel failed: %v", err)
	}
	if model != "MockModem X1000" {
		t.Errorf("Expected 'MockModem X1000', got '%s'", model)
	}

	// Test getting state
	state, err := mockModem.GetState()
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	t.Logf("Modem state: %s", state.String())

	// Test enabling modem
	err = mockModem.Enable(true)
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	// Verify state changed
	state, err = mockModem.GetState()
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	t.Logf("Modem state after enable: %s", state.String())
}

// TestMockModemWithErrors demonstrates testing error cases
func TestMockModemWithErrors(t *testing.T) {
	mockModem := mocks.NewMockModem()

	// Configure the mock to return an error
	mockModem.EnableError = &MockError{msg: "simulated enable error"}

	// Test that error is returned
	err := mockModem.Enable(true)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "simulated enable error" {
		t.Errorf("Expected 'simulated enable error', got '%s'", err.Error())
	}
}

// TestMockModem3gpp demonstrates testing 3GPP functionality
func TestMockModem3gpp(t *testing.T) {
	mock3gpp := mocks.NewMockModem3gpp()

	// Test getting IMEI
	imei, err := mock3gpp.GetImei()
	if err != nil {
		t.Fatalf("GetImei failed: %v", err)
	}
	if len(imei) != 15 {
		t.Errorf("IMEI should be 15 digits, got %d: %s", len(imei), imei)
	}

	// Test getting operator code
	opCode, err := mock3gpp.GetOperatorCode()
	if err != nil {
		t.Fatalf("GetOperatorCode failed: %v", err)
	}
	if opCode != "310260" {
		t.Errorf("Expected operator code '310260', got '%s'", opCode)
	}

	// Test scanning networks
	networks, err := mock3gpp.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if len(networks) == 0 {
		t.Error("Expected at least one network")
	}
	t.Logf("Found %d networks", len(networks))
	for _, net := range networks {
		t.Logf("  - %s (%s) [%s]", net.OperatorLong, net.OperatorShort, net.OperatorCode)
	}

	// Test registration
	err = mock3gpp.Register("")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
}

// TestMockBearer demonstrates testing bearer functionality
func TestMockBearer(t *testing.T) {
	mockBearer := mocks.NewMockBearer()

	// Initially should not be connected
	connected, err := mockBearer.GetConnected()
	if err != nil {
		t.Fatalf("GetConnected failed: %v", err)
	}
	if connected {
		t.Error("Bearer should not be connected initially")
	}

	// Connect the bearer
	err = mockBearer.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Should now be connected
	connected, err = mockBearer.GetConnected()
	if err != nil {
		t.Fatalf("GetConnected failed: %v", err)
	}
	if !connected {
		t.Error("Bearer should be connected after Connect()")
	}

	// Get interface name
	iface, err := mockBearer.GetInterface()
	if err != nil {
		t.Fatalf("GetInterface failed: %v", err)
	}
	t.Logf("Bearer interface: %s", iface)

	// Get IPv4 config
	ipv4Config, err := mockBearer.GetIp4Config()
	if err != nil {
		t.Fatalf("GetIp4Config failed: %v", err)
	}
	t.Logf("IPv4 Config: %s/%d, Gateway: %s",
		ipv4Config.Address, ipv4Config.Prefix, ipv4Config.Gateway)

	// Get stats
	stats, err := mockBearer.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	t.Logf("Stats: RX=%d bytes, TX=%d bytes", stats.BytesRx, stats.BytesTx)

	// Disconnect
	err = mockBearer.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}

	// Verify disconnected
	connected, err = mockBearer.GetConnected()
	if err != nil {
		t.Fatalf("GetConnected failed: %v", err)
	}
	if connected {
		t.Error("Bearer should not be connected after Disconnect()")
	}
}

// TestMockSim demonstrates testing SIM functionality
func TestMockSim(t *testing.T) {
	mockSim := mocks.NewMockSim()

	// Test getting SIM identifier
	simId, err := mockSim.GetSimIdentifier()
	if err != nil {
		t.Fatalf("GetSimIdentifier failed: %v", err)
	}
	t.Logf("SIM ID: %s", simId)

	// Test getting IMSI
	imsi, err := mockSim.GetImsi()
	if err != nil {
		t.Fatalf("GetImsi failed: %v", err)
	}
	if len(imsi) < 14 || len(imsi) > 15 {
		t.Errorf("IMSI should be 14-15 digits, got %d: %s", len(imsi), imsi)
	}

	// Test getting operator
	opName, err := mockSim.GetOperatorName()
	if err != nil {
		t.Fatalf("GetOperatorName failed: %v", err)
	}
	t.Logf("Operator: %s", opName)

	// Test sending PIN (should succeed with mock)
	err = mockSim.SendPin("1234")
	if err != nil {
		t.Fatalf("SendPin failed: %v", err)
	}
}

// TestMockModemSimple demonstrates testing Simple interface
func TestMockModemSimple(t *testing.T) {
	mockSimple := mocks.NewMockModemSimple()

	// Test getting status
	status, err := mockSimple.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	t.Logf("Status: %+v", status)

	// Test connecting (returns a bearer)
	bearer, err := mockSimple.Connect(status)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	if bearer == nil {
		t.Fatal("Expected bearer, got nil")
	}

	// Verify bearer
	bearerPath := bearer.GetObjectPath()
	t.Logf("Bearer created at: %s", bearerPath)

	// Test disconnecting
	err = mockSimple.Disconnect(bearerPath)
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}
}

// TestMockCustomization demonstrates customizing mock behavior
func TestMockCustomization(t *testing.T) {
	mockModem := mocks.NewMockModem()

	// Customize the mock's return values
	mockModem.ManufacturerValue = "Custom Manufacturer"
	mockModem.ModelValue = "Custom Model"
	mockModem.StateValue = 6 // MM_MODEM_STATE_ENABLED

	// Test customized values
	manufacturer, _ := mockModem.GetManufacturer()
	if manufacturer != "Custom Manufacturer" {
		t.Errorf("Expected 'Custom Manufacturer', got '%s'", manufacturer)
	}

	model, _ := mockModem.GetModel()
	if model != "Custom Model" {
		t.Errorf("Expected 'Custom Model', got '%s'", model)
	}
}

// TestIntegrationScenario demonstrates a complete workflow
func TestIntegrationScenario(t *testing.T) {
	// Setup: Create mock ModemManager
	mockMM := mocks.NewMockModemManager()

	// Step 1: Get version
	version, err := mockMM.GetVersion()
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}
	t.Logf("ModemManager version: %s", version)

	// Step 2: List modems
	modems, err := mockMM.GetModems()
	if err != nil {
		t.Fatalf("GetModems failed: %v", err)
	}
	if len(modems) == 0 {
		t.Fatal("No modems found")
	}
	t.Logf("Found %d modem(s)", len(modems))

	// Step 3: Get modem details
	modem := modems[0]
	manufacturer, _ := modem.GetManufacturer()
	model, _ := modem.GetModel()
	t.Logf("Modem: %s %s", manufacturer, model)

	// Step 4: Enable modem
	err = modem.Enable(true)
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}
	t.Log("Modem enabled")

	// Step 5: Get SIM
	sim, err := modem.GetSim()
	if err != nil {
		t.Fatalf("GetSim failed: %v", err)
	}
	imsi, _ := sim.GetImsi()
	t.Logf("SIM IMSI: %s", imsi)

	// Step 6: Get 3GPP interface
	modem3gpp, err := modem.Get3gpp()
	if err != nil {
		t.Fatalf("Get3gpp failed: %v", err)
	}
	operatorName, _ := modem3gpp.GetOperatorName()
	t.Logf("Registered to: %s", operatorName)

	// Step 7: Create and connect bearer
	mockModem := modem.(*mocks.MockModem)
	bearer, err := mockModem.CreateBearer(mocks.NewMockBearer().GetProperties())
	if err != nil {
		t.Fatalf("CreateBearer failed: %v", err)
	}

	err = bearer.Connect()
	if err != nil {
		t.Fatalf("Bearer.Connect failed: %v", err)
	}
	t.Log("Bearer connected")

	connected, _ := bearer.GetConnected()
	if !connected {
		t.Error("Bearer should be connected")
	}

	// Cleanup
	err = bearer.Disconnect()
	if err != nil {
		t.Fatalf("Bearer.Disconnect failed: %v", err)
	}
	t.Log("Bearer disconnected")
}

// MockError is a simple error type for testing
type MockError struct {
	msg string
}

func (e *MockError) Error() string {
	return e.msg
}
