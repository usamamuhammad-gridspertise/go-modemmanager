// Package mocks provides mock implementations of go-modemmanager interfaces for testing
package mocks

import (
	"encoding/json"
	"time"

	"github.com/godbus/dbus/v5"
	mm "github.com/maltegrosse/go-modemmanager"
)

// MockModemManager is a mock implementation of the ModemManager interface
type MockModemManager struct {
	// Configurable return values
	VersionValue       string
	ModemsValue        []mm.Modem
	ScanDevicesError   error
	SetLoggingError    error
	ReportEventError   error
	InhibitDeviceError error
	GetVersionError    error
	GetModemsError     error
	SignalChan         chan *dbus.Signal
}

// NewMockModemManager creates a new mock ModemManager with default values
func NewMockModemManager() *MockModemManager {
	return &MockModemManager{
		VersionValue: "1.12.8-mock",
		ModemsValue:  []mm.Modem{NewMockModem()},
		SignalChan:   make(chan *dbus.Signal, 10),
	}
}

func (m *MockModemManager) ScanDevices() error {
	return m.ScanDevicesError
}

func (m *MockModemManager) GetModems() ([]mm.Modem, error) {
	return m.ModemsValue, m.GetModemsError
}

func (m *MockModemManager) SetLogging(level mm.MMLoggingLevel) error {
	return m.SetLoggingError
}

func (m *MockModemManager) ReportKernelEvent(props mm.EventProperties) error {
	return m.ReportEventError
}

func (m *MockModemManager) InhibitDevice(uid string, inhibit bool) error {
	return m.InhibitDeviceError
}

func (m *MockModemManager) GetVersion() (string, error) {
	return m.VersionValue, m.GetVersionError
}

func (m *MockModemManager) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Version": m.VersionValue,
	})
}

func (m *MockModemManager) SubscribePropertiesChanged() <-chan *dbus.Signal {
	return m.SignalChan
}

func (m *MockModemManager) ParsePropertiesChanged(v *dbus.Signal) (interfaceName string, changedProperties map[string]dbus.Variant, invalidatedProperties []string, err error) {
	return "", nil, nil, nil
}

func (m *MockModemManager) Unsubscribe() {}

// MockModem is a mock implementation of the Modem interface
type MockModem struct {
	// Configurable return values
	ObjectPathValue            dbus.ObjectPath
	ManufacturerValue          string
	ModelValue                 string
	RevisionValue              string
	EquipmentIdentifierValue   string
	DeviceIdentifierValue      string
	StateValue                 mm.MMModemState
	SignalQualityValue         mm.SignalQuality
	AccessTechnologiesValue    []mm.MMModemAccessTechnology
	UnlockRequiredValue        mm.MMModemLock
	PowerStateValue            mm.MMModemPowerState
	SupportedCapabilitiesValue [][]mm.MMModemCapability
	CurrentCapabilitiesValue   []mm.MMModemCapability
	SupportedModesValue        []mm.Mode
	CurrentModesValue          mm.Mode
	SupportedBandsValue        []mm.MMModemBand
	CurrentBandsValue          []mm.MMModemBand

	// Error values
	EnableError            error
	ListBearsError         error
	CreateBearerError      error
	DeleteBearerError      error
	ResetError             error
	FactoryResetError      error
	SetPowerStateError     error
	SetCapabilitiesError   error
	SetModesError          error
	SetBandsError          error
	CommandError           error
	GetSimpleModemError    error
	Get3gppError           error
	GetCdmaError           error
	GetTimeError           error
	GetFirmwareError       error
	GetSignalError         error
	GetOmaError            error
	GetLocationError       error
	GetMessagingError      error
	GetVoiceError          error
	GetSimError            error
	GetPropertiesError     error
	GetStateError          error
	GetMaxBearsError       error
	GetMaxActiveBearsError error
}

// NewMockModem creates a new mock Modem with default values
func NewMockModem() *MockModem {
	return &MockModem{
		ObjectPathValue:            "/org/freedesktop/ModemManager1/Modem/0",
		ManufacturerValue:          "MockModem Inc.",
		ModelValue:                 "MockModem X1000",
		RevisionValue:              "1.0.0",
		EquipmentIdentifierValue:   "IMEI123456789012345",
		DeviceIdentifierValue:      "mock-0000",
		StateValue:                 mm.MmModemStateRegistered,
		SignalQualityValue:         mm.SignalQuality{Quality: 75, Recent: true},
		AccessTechnologiesValue:    []mm.MMModemAccessTechnology{mm.MmModemAccessTechnologyLte},
		UnlockRequiredValue:        mm.MmModemLockNone,
		PowerStateValue:            mm.MmModemPowerStateOn,
		SupportedCapabilitiesValue: [][]mm.MMModemCapability{{mm.MmModemCapabilityLte}},
		CurrentCapabilitiesValue:   []mm.MMModemCapability{mm.MmModemCapabilityLte},
		SupportedModesValue: []mm.
			Mode{{AllowedModes: []mm.MMModemMode{mm.MmModemModeAny}}},
		CurrentModesValue:   mm.Mode{AllowedModes: []mm.MMModemMode{mm.MmModemMode4g}},
		SupportedBandsValue: []mm.MMModemBand{mm.MmModemBandEutran1, mm.MmModemBandEutran2},
		CurrentBandsValue:   []mm.MMModemBand{mm.MmModemBandEutran1},
	}
}

func (m *MockModem) GetObjectPath() dbus.ObjectPath {
	return m.ObjectPathValue
}

func (m *MockModem) GetSimpleModem() (mm.ModemSimple, error) {
	if m.GetSimpleModemError != nil {
		return nil, m.GetSimpleModemError
	}
	return NewMockModemSimple(), nil
}

func (m *MockModem) Get3gpp() (mm.Modem3gpp, error) {
	if m.Get3gppError != nil {
		return nil, m.Get3gppError
	}
	return NewMockModem3gpp(), nil
}

func (m *MockModem) GetCdma() (mm.ModemCdma, error) {
	return nil, m.GetCdmaError
}

func (m *MockModem) GetTime() (mm.ModemTime, error) {
	return nil, m.GetTimeError
}

func (m *MockModem) GetFirmware() (mm.ModemFirmware, error) {
	return nil, m.GetFirmwareError
}

func (m *MockModem) GetSignal() (mm.ModemSignal, error) {
	return nil, m.GetSignalError
}

func (m *MockModem) GetOma() (mm.ModemOma, error) {
	return nil, m.GetOmaError
}

func (m *MockModem) GetLocation() (mm.ModemLocation, error) {
	return nil, m.GetLocationError
}

func (m *MockModem) GetMessaging() (mm.ModemMessaging, error) {
	return nil, m.GetMessagingError
}

func (m *MockModem) GetVoice() (mm.ModemVoice, error) {
	return nil, m.GetVoiceError
}

func (m *MockModem) Enable(enable bool) error {
	if enable {
		m.StateValue = mm.MmModemStateEnabled
	} else {
		m.StateValue = mm.MmModemStateDisabled
	}
	return m.EnableError
}

func (m *MockModem) ListBearers() ([]mm.Bearer, error) {
	return []mm.Bearer{}, m.ListBearsError
}

func (m *MockModem) CreateBearer(property mm.BearerProperty) (mm.Bearer, error) {
	return NewMockBearer(), m.CreateBearerError
}

func (m *MockModem) DeleteBearer(bearer mm.Bearer) error {
	return m.DeleteBearerError
}

func (m *MockModem) Reset() error {
	m.StateValue = mm.MmModemStateDisabled
	return m.ResetError
}

func (m *MockModem) FactoryReset(code string) error {
	return m.FactoryResetError
}

func (m *MockModem) SetPowerState(state mm.MMModemPowerState) error {
	m.PowerStateValue = state
	return m.SetPowerStateError
}

func (m *MockModem) SetCurrentCapabilities(capabilities []mm.MMModemCapability) error {
	m.CurrentCapabilitiesValue = capabilities
	return m.SetCapabilitiesError
}

func (m *MockModem) SetCurrentModes(property mm.Mode) error {
	m.CurrentModesValue = property
	return m.SetModesError
}

func (m *MockModem) SetCurrentBands(bands []mm.MMModemBand) error {
	m.CurrentBandsValue = bands
	return m.SetBandsError
}

func (m *MockModem) Command(cmd string, timeout uint32) (string, error) {
	return "OK", m.CommandError
}

func (m *MockModem) GetSim() (mm.Sim, error) {
	if m.GetSimError != nil {
		return nil, m.GetSimError
	}
	return NewMockSim(), nil
}

func (m *MockModem) GetProperties() (mm.ModemProperty, error) {
	return mm.ModemProperty{
		Manufacturer:        m.ManufacturerValue,
		Model:               m.ModelValue,
		Revision:            m.RevisionValue,
		EquipmentIdentifier: m.EquipmentIdentifierValue,
		DeviceIdentifier:    m.DeviceIdentifierValue,
	}, m.GetPropertiesError
}

func (m *MockModem) GetState() (mm.MMModemState, error) {
	return m.StateValue, m.GetStateError
}

func (m *MockModem) GetSignalQuality() (mm.SignalQuality, error) {
	return m.SignalQualityValue, nil
}

func (m *MockModem) GetAccessTechnologies() ([]mm.MMModemAccessTechnology, error) {
	return m.AccessTechnologiesValue, nil
}

func (m *MockModem) GetUnlockRequired() (mm.MMModemLock, error) {
	return m.UnlockRequiredValue, nil
}

func (m *MockModem) GetPowerState() (mm.MMModemPowerState, error) {
	return m.PowerStateValue, nil
}

func (m *MockModem) GetSupportedCapabilities() ([][]mm.MMModemCapability, error) {
	return m.SupportedCapabilitiesValue, nil
}

func (m *MockModem) GetCurrentCapabilities() ([]mm.MMModemCapability, error) {
	return m.CurrentCapabilitiesValue, nil
}

func (m *MockModem) GetMaxBearers() (uint32, error) {
	return 1, m.GetMaxBearsError
}

func (m *MockModem) GetMaxActiveBearers() (uint32, error) {
	return 1, m.GetMaxActiveBearsError
}

func (m *MockModem) GetManufacturer() (string, error) {
	return m.ManufacturerValue, nil
}

func (m *MockModem) GetModel() (string, error) {
	return m.ModelValue, nil
}

func (m *MockModem) GetRevision() (string, error) {
	return m.RevisionValue, nil
}

func (m *MockModem) GetEquipmentIdentifier() (string, error) {
	return m.EquipmentIdentifierValue, nil
}

func (m *MockModem) GetDeviceIdentifier() (string, error) {
	return m.DeviceIdentifierValue, nil
}

func (m *MockModem) GetOwnNumbers() ([]string, error) {
	return []string{"+1234567890"}, nil
}

func (m *MockModem) GetSupportedModes() ([]mm.Mode, error) {
	return m.SupportedModesValue, nil
}

func (m *MockModem) GetCurrentModes() (mm.Mode, error) {
	return m.CurrentModesValue, nil
}

func (m *MockModem) GetSupportedBands() ([]mm.MMModemBand, error) {
	return m.SupportedBandsValue, nil
}

func (m *MockModem) GetCurrentBands() ([]mm.MMModemBand, error) {
	return m.CurrentBandsValue, nil
}

func (m *MockModem) GetSupportedIpFamilies() (mm.MMBearerIpFamily, error) {
	return mm.MmBearerIpFamilyIpv4 | mm.MmBearerIpFamilyIpv6, nil
}

func (m *MockModem) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Manufacturer":        m.ManufacturerValue,
		"Model":               m.ModelValue,
		"Revision":            m.RevisionValue,
		"EquipmentIdentifier": m.EquipmentIdentifierValue,
		"State":               m.StateValue.String(),
	})
}

func (m *MockModem) SubscribeStateChanged() <-chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)
	return ch
}

func (m *MockModem) ParseStateChanged(v *dbus.Signal) (old mm.MMModemState, new mm.MMModemState, reason mm.MMModemStateChangeReason, err error) {
	return mm.MmModemStateDisabled, mm.MmModemStateEnabled, mm.MmModemStateChangeReasonUserRequested, nil
}

func (m *MockModem) SubscribePropertiesChanged() <-chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)
	return ch
}

func (m *MockModem) ParsePropertiesChanged(v *dbus.Signal) (interfaceName string, changedProperties map[string]dbus.Variant, invalidatedProperties []string, err error) {
	return "", nil, nil, nil
}

func (m *MockModem) Unsubscribe() {}

// MockModemSimple is a mock implementation of ModemSimple interface
type MockModemSimple struct {
	ConnectError    error
	DisconnectError error
	GetStatusError  error
	StatusValue     mm.SimpleProperty
	BearerPathValue dbus.ObjectPath
	ObjectPathValue dbus.ObjectPath
}

func NewMockModemSimple() *MockModemSimple {
	return &MockModemSimple{
		StatusValue:     mm.SimpleProperty{},
		BearerPathValue: "/org/freedesktop/ModemManager1/Bearer/0",
		ObjectPathValue: "/org/freedesktop/ModemManager1/Modem/0",
	}
}

func (m *MockModemSimple) GetObjectPath() dbus.ObjectPath {
	return m.ObjectPathValue
}

func (m *MockModemSimple) Connect(property mm.SimpleProperty) (mm.Bearer, error) {
	if m.ConnectError != nil {
		return nil, m.ConnectError
	}
	return NewMockBearer(), nil
}

func (m *MockModemSimple) Disconnect(bearerPath dbus.ObjectPath) error {
	return m.DisconnectError
}

func (m *MockModemSimple) GetStatus() (mm.SimpleProperty, error) {
	return m.StatusValue, m.GetStatusError
}

// MockModem3gpp is a mock implementation of Modem3gpp interface
type MockModem3gpp struct {
	ObjectPathValue        dbus.ObjectPath
	ImeiValue              string
	RegistrationStateValue mm.MMModem3gppRegistrationState
	OperatorCodeValue      string
	OperatorNameValue      string
	RegisterError          error
	ScanError              error
}

func NewMockModem3gpp() *MockModem3gpp {
	return &MockModem3gpp{
		ObjectPathValue:        "/org/freedesktop/ModemManager1/Modem/0",
		ImeiValue:              "123456789012345",
		RegistrationStateValue: mm.MmModem3gppRegistrationStateHome,
		OperatorCodeValue:      "310260",
		OperatorNameValue:      "T-Mobile",
	}
}

func (m *MockModem3gpp) GetObjectPath() dbus.ObjectPath {
	return m.ObjectPathValue
}

func (m *MockModem3gpp) GetUssd() (mm.Ussd, error) {
	return nil, nil
}

func (m *MockModem3gpp) Register(operatorId string) error {
	return m.RegisterError
}

func (m *MockModem3gpp) Scan() ([]mm.Modem3gppNetwork, error) {
	return []mm.Modem3gppNetwork{
		{
			OperatorLong:  "T-Mobile",
			OperatorShort: "TMO",
			OperatorCode:  "310260",
		},
	}, m.ScanError
}

func (m *MockModem3gpp) GetImei() (string, error) {
	return m.ImeiValue, nil
}

func (m *MockModem3gpp) GetRegistrationState() (mm.MMModem3gppRegistrationState, error) {
	return m.RegistrationStateValue, nil
}

func (m *MockModem3gpp) GetOperatorCode() (string, error) {
	return m.OperatorCodeValue, nil
}

func (m *MockModem3gpp) GetOperatorName() (string, error) {
	return m.OperatorNameValue, nil
}

func (m *MockModem3gpp) SetEpsUeModeOperation(mode mm.MMModem3gppEpsUeModeOperation) error {
	return nil
}

func (m *MockModem3gpp) SetInitialEpsBearerSettings(property mm.BearerProperty) error {
	return nil
}

func (m *MockModem3gpp) GetEnabledFacilityLocks() ([]mm.MMModem3gppFacility, error) {
	return []mm.MMModem3gppFacility{}, nil
}

func (m *MockModem3gpp) GetEpsUeModeOperation() (mm.MMModem3gppEpsUeModeOperation, error) {
	return mm.MmModem3gppEpsUeModeOperationPsMode2, nil
}

func (m *MockModem3gpp) GetPco() ([]mm.RawPcoData, error) {
	return []mm.RawPcoData{}, nil
}

func (m *MockModem3gpp) GetInitialEpsBearer() (mm.Bearer, error) {
	return NewMockBearer(), nil
}

func (m *MockModem3gpp) GetInitialEpsBearerSettings() (mm.BearerProperty, error) {
	return mm.BearerProperty{}, nil
}

func (m *MockModem3gpp) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Imei":              m.ImeiValue,
		"RegistrationState": m.RegistrationStateValue.String(),
		"OperatorCode":      m.OperatorCodeValue,
		"OperatorName":      m.OperatorNameValue,
	})
}

func (m *MockModem3gpp) SubscribePropertiesChanged() <-chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)
	return ch
}

func (m *MockModem3gpp) ParsePropertiesChanged(v *dbus.Signal) (interfaceName string, changedProperties map[string]dbus.Variant, invalidatedProperties []string, err error) {
	return "", nil, nil, nil
}

func (m *MockModem3gpp) Unsubscribe() {}

// MockBearer is a mock implementation of Bearer interface
type MockBearer struct {
	ObjectPathValue dbus.ObjectPath
	ConnectedValue  bool
	InterfaceValue  string
	Ipv4ConfigValue mm.IpConfig
	Ipv6ConfigValue mm.IpConfig
	ConnectError    error
	DisconnectError error
}

func NewMockBearer() *MockBearer {
	return &MockBearer{
		ObjectPathValue: "/org/freedesktop/ModemManager1/Bearer/0",
		ConnectedValue:  false,
		InterfaceValue:  "wwan0",
		Ipv4ConfigValue: mm.IpConfig{
			Method:  mm.MmBearerIpMethodStatic,
			Address: "192.168.1.100",
			Prefix:  24,
			Gateway: "192.168.1.1",
			Dns:     []string{"8.8.8.8", "8.8.4.4"},
		},
	}
}

func (b *MockBearer) GetObjectPath() dbus.ObjectPath {
	return b.ObjectPathValue
}

func (b *MockBearer) Connect() error {
	b.ConnectedValue = true
	return b.ConnectError
}

func (b *MockBearer) Disconnect() error {
	b.ConnectedValue = false
	return b.DisconnectError
}

func (b *MockBearer) GetInterface() (string, error) {
	return b.InterfaceValue, nil
}

func (b *MockBearer) GetConnected() (bool, error) {
	return b.ConnectedValue, nil
}

func (b *MockBearer) GetSuspended() (bool, error) {
	return false, nil
}

func (b *MockBearer) GetIp4Config() (mm.IpConfig, error) {
	return b.Ipv4ConfigValue, nil
}

func (b *MockBearer) GetIp6Config() (mm.IpConfig, error) {
	return b.Ipv6ConfigValue, nil
}

func (b *MockBearer) GetIpTimeout() (uint32, error) {
	return 20, nil
}

func (b *MockBearer) GetProperties() (mm.BearerProperty, error) {
	return mm.BearerProperty{
		Apn:          "internet",
		IpType:       mm.MmBearerIpFamilyIpv4,
		AllowRoaming: false,
	}, nil
}

func (b *MockBearer) GetStats() (mm.BearerStats, error) {
	return mm.BearerStats{
		StartDate: time.Now().Unix(),
		BytesRx:   1024000,
		BytesTx:   512000,
	}, nil
}

func (b *MockBearer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Connected": b.ConnectedValue,
		"Interface": b.InterfaceValue,
	})
}

func (b *MockBearer) SubscribePropertiesChanged() <-chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)
	return ch
}

func (b *MockBearer) ParsePropertiesChanged(v *dbus.Signal) (interfaceName string, changedProperties map[string]dbus.Variant, invalidatedProperties []string, err error) {
	return "", nil, nil, nil
}

func (b *MockBearer) Unsubscribe() {}

// MockSim is a mock implementation of Sim interface
type MockSim struct {
	ObjectPathValue         dbus.ObjectPath
	SimIdentifierValue      string
	ImsiValue               string
	OperatorIdentifierValue string
	OperatorNameValue       string
	SendPinError            error
	SendPukError            error
	EnablePinError          error
	ChangePinError          error
}

func NewMockSim() *MockSim {
	return &MockSim{
		ObjectPathValue:         "/org/freedesktop/ModemManager1/Sim/0",
		SimIdentifierValue:      "89012345678901234567",
		ImsiValue:               "310260123456789",
		OperatorIdentifierValue: "310260",
		OperatorNameValue:       "T-Mobile",
	}
}

func (s *MockSim) GetObjectPath() dbus.ObjectPath {
	return s.ObjectPathValue
}

func (s *MockSim) SendPin(pin string) error {
	return s.SendPinError
}

func (s *MockSim) SendPuk(puk, pin string) error {
	return s.SendPukError
}

func (s *MockSim) EnablePin(pin string, enabled bool) error {
	return s.EnablePinError
}

func (s *MockSim) ChangePin(oldPin, newPin string) error {
	return s.ChangePinError
}

func (s *MockSim) GetSimIdentifier() (string, error) {
	return s.SimIdentifierValue, nil
}

func (s *MockSim) GetImsi() (string, error) {
	return s.ImsiValue, nil
}

func (s *MockSim) GetOperatorIdentifier() (string, error) {
	return s.OperatorIdentifierValue, nil
}

func (s *MockSim) GetOperatorName() (string, error) {
	return s.OperatorNameValue, nil
}

func (s *MockSim) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"SimIdentifier":      s.SimIdentifierValue,
		"Imsi":               s.ImsiValue,
		"OperatorIdentifier": s.OperatorIdentifierValue,
		"OperatorName":       s.OperatorNameValue,
	})
}

func (s *MockSim) SubscribePropertiesChanged() <-chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)
	return ch
}

func (s *MockSim) ParsePropertiesChanged(v *dbus.Signal) (interfaceName string, changedProperties map[string]dbus.Variant, invalidatedProperties []string, err error) {
	return "", nil, nil, nil
}

func (s *MockSim) Unsubscribe() {}
