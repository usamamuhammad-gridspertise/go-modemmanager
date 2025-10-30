#!/usr/bin/env python3
"""
Mock ModemManager D-Bus Service
Simulates ModemManager for testing go-modemmanager without hardware

This mock service provides:
- org.freedesktop.ModemManager1 interface
- Virtual modem devices
- Basic modem operations (enable, disable, connect, etc.)
- SMS and location capabilities
- Signal updates

Usage:
    python3 mock_modemmanager.py

Requirements:
    pip3 install dbus-python pygobject
"""

import dbus
import dbus.service
import dbus.mainloop.glib
from gi.repository import GLib
import logging
import sys

# Configure logging
logging.basicConfig(
    level=logging.DEBUG, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger("MockModemManager")

# D-Bus constants
DBUS_BUS_NAME = "org.freedesktop.ModemManager1"
DBUS_OBJECT_PATH = "/org/freedesktop/ModemManager1"
MM_INTERFACE = "org.freedesktop.ModemManager1"
MODEM_INTERFACE = "org.freedesktop.ModemManager1.Modem"
MODEM_3GPP_INTERFACE = "org.freedesktop.ModemManager1.Modem.Modem3gpp"
MODEM_SIMPLE_INTERFACE = "org.freedesktop.ModemManager1.Modem.Simple"
MODEM_LOCATION_INTERFACE = "org.freedesktop.ModemManager1.Modem.Location"
MODEM_MESSAGING_INTERFACE = "org.freedesktop.ModemManager1.Modem.Messaging"
SIM_INTERFACE = "org.freedesktop.ModemManager1.Sim"
BEARER_INTERFACE = "org.freedesktop.ModemManager1.Bearer"
SMS_INTERFACE = "org.freedesktop.ModemManager1.Sms"

# Modem State Enums
MM_MODEM_STATE_FAILED = -1
MM_MODEM_STATE_UNKNOWN = 0
MM_MODEM_STATE_INITIALIZING = 1
MM_MODEM_STATE_LOCKED = 2
MM_MODEM_STATE_DISABLED = 3
MM_MODEM_STATE_DISABLING = 4
MM_MODEM_STATE_ENABLING = 5
MM_MODEM_STATE_ENABLED = 6
MM_MODEM_STATE_SEARCHING = 7
MM_MODEM_STATE_REGISTERED = 8
MM_MODEM_STATE_DISCONNECTING = 9
MM_MODEM_STATE_CONNECTING = 10
MM_MODEM_STATE_CONNECTED = 11

# Modem Capability Enums
MM_MODEM_CAPABILITY_NONE = 0
MM_MODEM_CAPABILITY_GSM_UMTS = 4
MM_MODEM_CAPABILITY_LTE = 8
MM_MODEM_CAPABILITY_LTE_ADVANCED = 16

# Access Technology Enums
MM_MODEM_ACCESS_TECHNOLOGY_LTE = 16384


class MockModemManagerService(dbus.service.Object):
    """Main ModemManager service"""

    def __init__(self, bus, object_path):
        super().__init__(bus, object_path)
        self.modems = []
        self.version = "1.12.8-mock"
        logger.info(f"MockModemManager service started at {object_path}")

        # Create a mock modem by default
        self._create_mock_modem()

    def _create_mock_modem(self):
        """Create a virtual modem"""
        modem_path = f"{DBUS_OBJECT_PATH}/Modem/0"
        modem = MockModem(self.connection, modem_path, index=0)
        self.modems.append(modem_path)
        logger.info(f"Created mock modem at {modem_path}")

    @dbus.service.method(MM_INTERFACE, out_signature="s")
    def GetVersion(self):
        """Get ModemManager version"""
        logger.debug(f"GetVersion() called, returning: {self.version}")
        return self.version

    @dbus.service.method(MM_INTERFACE, in_signature="", out_signature="")
    def ScanDevices(self):
        """Scan for new devices"""
        logger.info("ScanDevices() called")
        # In real MM, this scans for hardware
        return

    @dbus.service.method(MM_INTERFACE, in_signature="s", out_signature="")
    def SetLogging(self, level):
        """Set logging level"""
        logger.info(f"SetLogging({level}) called")
        return

    @dbus.service.method(MM_INTERFACE, in_signature="sb", out_signature="")
    def InhibitDevice(self, uid, inhibit):
        """Inhibit a device"""
        logger.info(f"InhibitDevice({uid}, {inhibit}) called")
        return

    @dbus.service.method(dbus.PROPERTIES_IFACE, in_signature="ss", out_signature="v")
    def Get(self, interface_name, property_name):
        """Get property value"""
        logger.debug(f"Get({interface_name}, {property_name}) called")
        if property_name == "Version":
            return self.version
        raise dbus.exceptions.DBusException(
            f"Property {property_name} not found",
            name="org.freedesktop.DBus.Error.InvalidArgs",
        )

    @dbus.service.method(dbus.PROPERTIES_IFACE, in_signature="s", out_signature="a{sv}")
    def GetAll(self, interface_name):
        """Get all properties"""
        logger.debug(f"GetAll({interface_name}) called")
        return {"Version": self.version}


class MockModem(dbus.service.Object):
    """Mock Modem object"""

    def __init__(self, bus, object_path, index=0):
        super().__init__(bus, object_path)
        self.index = index
        self.state = MM_MODEM_STATE_DISABLED
        self.enabled = False

        # Modem properties
        self.manufacturer = "MockModem Inc."
        self.model = "MockModem X1000"
        self.revision = "1.0.0"
        self.device_identifier = f"mock-{index:04d}"
        self.equipment_identifier = f"IMEI{123456789012345 + index}"
        self.own_numbers = ["+1234567890"]
        self.unlock_required = 1  # MM_MODEM_LOCK_NONE
        self.power_state = 3  # MM_MODEM_POWER_STATE_ON
        self.supported_capabilities = [
            MM_MODEM_CAPABILITY_GSM_UMTS | MM_MODEM_CAPABILITY_LTE
        ]
        self.current_capabilities = MM_MODEM_CAPABILITY_LTE
        self.max_bearers = 1
        self.max_active_bearers = 1
        self.signal_quality = (75, True)  # (quality, recent)
        self.access_technologies = MM_MODEM_ACCESS_TECHNOLOGY_LTE
        self.supported_modes = [(519, 0)]  # All modes, no preferred
        self.current_modes = (519, 0)
        self.supported_bands = [0, 1, 2, 3]  # Mock bands
        self.current_bands = [0, 1, 2]
        self.supported_ip_families = 11  # IPv4 | IPv6 | IPv4v6

        # SIM
        self.sim_path = f"{object_path}/Sim/0"
        self.sim = MockSim(bus, self.sim_path)

        # Bearers
        self.bearers = []

        # 3GPP properties
        self.operator_code = "310260"
        self.operator_name = "T-Mobile"
        self.registration_state = 1  # MM_MODEM_3GPP_REGISTRATION_STATE_HOME

        logger.info(f"MockModem created at {object_path}")

    @dbus.service.method(MODEM_INTERFACE, in_signature="b", out_signature="")
    def Enable(self, enable):
        """Enable or disable the modem"""
        logger.info(f"Modem.Enable({enable}) called")
        if enable:
            self.state = MM_MODEM_STATE_ENABLING
            GLib.timeout_add(1000, self._finish_enable)
        else:
            self.state = MM_MODEM_STATE_DISABLING
            GLib.timeout_add(1000, self._finish_disable)
        self.enabled = enable
        self._emit_state_changed()

    def _finish_enable(self):
        """Finish enabling modem"""
        self.state = MM_MODEM_STATE_ENABLED
        self._emit_state_changed()
        # Transition to registered
        GLib.timeout_add(500, self._register)
        return False

    def _finish_disable(self):
        """Finish disabling modem"""
        self.state = MM_MODEM_STATE_DISABLED
        self._emit_state_changed()
        return False

    def _register(self):
        """Register with network"""
        self.state = MM_MODEM_STATE_SEARCHING
        self._emit_state_changed()
        GLib.timeout_add(1000, self._finish_register)
        return False

    def _finish_register(self):
        """Finish registration"""
        self.state = MM_MODEM_STATE_REGISTERED
        self._emit_state_changed()
        return False

    def _emit_state_changed(self):
        """Emit StateChanged signal"""
        logger.debug(f"Emitting StateChanged signal: {self.state}")
        self.StateChanged(self.state, self.state, 1)

    @dbus.service.method(MODEM_INTERFACE, in_signature="a{sv}", out_signature="o")
    def CreateBearer(self, properties):
        """Create a new bearer"""
        logger.info(f"Modem.CreateBearer({properties}) called")
        bearer_index = len(self.bearers)
        bearer_path = f"{self._object_path}/Bearer/{bearer_index}"
        bearer = MockBearer(self.connection, bearer_path, properties)
        self.bearers.append(bearer_path)
        return bearer_path

    @dbus.service.method(MODEM_INTERFACE, in_signature="o", out_signature="")
    def DeleteBearer(self, bearer_path):
        """Delete a bearer"""
        logger.info(f"Modem.DeleteBearer({bearer_path}) called")
        if bearer_path in self.bearers:
            self.bearers.remove(bearer_path)

    @dbus.service.method(MODEM_INTERFACE, in_signature="", out_signature="")
    def Reset(self):
        """Reset the modem"""
        logger.info("Modem.Reset() called")
        self.state = MM_MODEM_STATE_DISABLED
        self._emit_state_changed()

    @dbus.service.method(MODEM_INTERFACE, in_signature="s", out_signature="s")
    def Command(self, cmd):
        """Send AT command"""
        logger.info(f"Modem.Command({cmd}) called")
        return "OK"

    @dbus.service.signal(MODEM_INTERFACE, signature="iiu")
    def StateChanged(self, old, new, reason):
        """StateChanged signal"""
        pass

    @dbus.service.method(dbus.PROPERTIES_IFACE, in_signature="ss", out_signature="v")
    def Get(self, interface_name, property_name):
        """Get property value"""
        logger.debug(f"Modem.Get({interface_name}, {property_name}) called")

        property_map = {
            "State": self.state,
            "Manufacturer": self.manufacturer,
            "Model": self.model,
            "Revision": self.revision,
            "DeviceIdentifier": self.device_identifier,
            "EquipmentIdentifier": self.equipment_identifier,
            "Device": "/sys/devices/mock",
            "Drivers": ["mock_driver"],
            "Plugin": "Mock",
            "PrimaryPort": "ttyUSB0",
            "Ports": [("ttyUSB0", dbus.UInt32(1)), ("ttyUSB1", dbus.UInt32(2))],
            "Sim": dbus.ObjectPath(self.sim_path),
            "Bearers": [dbus.ObjectPath(b) for b in self.bearers],
            "SupportedCapabilities": dbus.Array(
                self.supported_capabilities, signature="u"
            ),
            "CurrentCapabilities": dbus.UInt32(self.current_capabilities),
            "MaxBearers": dbus.UInt32(self.max_bearers),
            "MaxActiveBearers": dbus.UInt32(self.max_active_bearers),
            "OwnNumbers": self.own_numbers,
            "UnlockRequired": dbus.UInt32(self.unlock_required),
            "PowerState": dbus.UInt32(self.power_state),
            "SignalQuality": dbus.Struct(self.signal_quality, signature="ub"),
            "AccessTechnologies": dbus.UInt32(self.access_technologies),
            "SupportedModes": dbus.Array(
                [
                    dbus.Struct((dbus.UInt32(m[0]), dbus.UInt32(m[1])), signature="uu")
                    for m in self.supported_modes
                ],
                signature="(uu)",
            ),
            "CurrentModes": dbus.Struct(
                (
                    dbus.UInt32(self.current_modes[0]),
                    dbus.UInt32(self.current_modes[1]),
                ),
                signature="uu",
            ),
            "SupportedBands": dbus.Array(
                [dbus.UInt32(b) for b in self.supported_bands], signature="u"
            ),
            "CurrentBands": dbus.Array(
                [dbus.UInt32(b) for b in self.current_bands], signature="u"
            ),
            "SupportedIpFamilies": dbus.UInt32(self.supported_ip_families),
        }

        if property_name in property_map:
            return property_map[property_name]

        raise dbus.exceptions.DBusException(
            f"Property {property_name} not found",
            name="org.freedesktop.DBus.Error.InvalidArgs",
        )

    @dbus.service.method(dbus.PROPERTIES_IFACE, in_signature="s", out_signature="a{sv}")
    def GetAll(self, interface_name):
        """Get all properties"""
        logger.debug(f"Modem.GetAll({interface_name}) called")

        props = {
            "State": dbus.Int32(self.state),
            "Manufacturer": self.manufacturer,
            "Model": self.model,
            "Revision": self.revision,
            "DeviceIdentifier": self.device_identifier,
            "EquipmentIdentifier": self.equipment_identifier,
        }
        return props


class MockSim(dbus.service.Object):
    """Mock SIM card object"""

    def __init__(self, bus, object_path):
        super().__init__(bus, object_path)
        self.sim_identifier = "89012345678901234567"
        self.imsi = "310260123456789"
        self.operator_identifier = "310260"
        self.operator_name = "T-Mobile"
        logger.info(f"MockSim created at {object_path}")

    @dbus.service.method(SIM_INTERFACE, in_signature="s", out_signature="")
    def SendPin(self, pin):
        """Send PIN"""
        logger.info(f"Sim.SendPin({pin}) called")
        return

    @dbus.service.method(SIM_INTERFACE, in_signature="ss", out_signature="")
    def SendPuk(self, puk, pin):
        """Send PUK"""
        logger.info(f"Sim.SendPuk({puk}, {pin}) called")
        return

    @dbus.service.method(dbus.PROPERTIES_IFACE, in_signature="ss", out_signature="v")
    def Get(self, interface_name, property_name):
        """Get property value"""
        logger.debug(f"Sim.Get({interface_name}, {property_name}) called")

        property_map = {
            "SimIdentifier": self.sim_identifier,
            "Imsi": self.imsi,
            "OperatorIdentifier": self.operator_identifier,
            "OperatorName": self.operator_name,
        }

        if property_name in property_map:
            return property_map[property_name]

        raise dbus.exceptions.DBusException(
            f"Property {property_name} not found",
            name="org.freedesktop.DBus.Error.InvalidArgs",
        )


class MockBearer(dbus.service.Object):
    """Mock Bearer object"""

    def __init__(self, bus, object_path, properties):
        super().__init__(bus, object_path)
        self.properties = properties
        self.connected = False
        self.interface = "wwan0"
        self.ip_type = properties.get("ip-type", "ipv4")
        self.apn = properties.get("apn", "internet")
        logger.info(f"MockBearer created at {object_path}")

    @dbus.service.method(BEARER_INTERFACE, in_signature="", out_signature="")
    def Connect(self):
        """Connect the bearer"""
        logger.info("Bearer.Connect() called")
        self.connected = True
        return

    @dbus.service.method(BEARER_INTERFACE, in_signature="", out_signature="")
    def Disconnect(self):
        """Disconnect the bearer"""
        logger.info("Bearer.Disconnect() called")
        self.connected = False
        return

    @dbus.service.method(dbus.PROPERTIES_IFACE, in_signature="ss", out_signature="v")
    def Get(self, interface_name, property_name):
        """Get property value"""
        logger.debug(f"Bearer.Get({interface_name}, {property_name}) called")

        property_map = {
            "Connected": self.connected,
            "Interface": self.interface,
            "IpType": self.ip_type,
        }

        if property_name in property_map:
            return property_map[property_name]

        raise dbus.exceptions.DBusException(
            f"Property {property_name} not found",
            name="org.freedesktop.DBus.Error.InvalidArgs",
        )


def main():
    """Main entry point"""
    logger.info("Starting Mock ModemManager D-Bus Service...")

    # Set up D-Bus main loop
    dbus.mainloop.glib.DBusGMainLoop(set_as_default=True)

    # Connect to system bus
    try:
        bus = dbus.SystemBus()
        logger.info("Connected to system bus")
    except Exception as e:
        logger.error(f"Failed to connect to system bus: {e}")
        logger.info("Trying session bus instead...")
        try:
            bus = dbus.SessionBus()
            logger.info("Connected to session bus")
        except Exception as e:
            logger.error(f"Failed to connect to session bus: {e}")
            sys.exit(1)

    # Request bus name
    try:
        name = dbus.service.BusName(DBUS_BUS_NAME, bus)
        logger.info(f"Acquired bus name: {DBUS_BUS_NAME}")
    except Exception as e:
        logger.error(f"Failed to acquire bus name: {e}")
        logger.error(
            "Is ModemManager already running? Stop it with: systemctl stop ModemManager"
        )
        sys.exit(1)

    # Create the service
    service = MockModemManagerService(bus, DBUS_OBJECT_PATH)

    logger.info("=" * 60)
    logger.info("Mock ModemManager is now running!")
    logger.info(f"Service: {DBUS_BUS_NAME}")
    logger.info(f"Object Path: {DBUS_OBJECT_PATH}")
    logger.info(f"Mock Modem: {DBUS_OBJECT_PATH}/Modem/0")
    logger.info("=" * 60)
    logger.info("Test with: mmcli -L")
    logger.info("Stop with: Ctrl+C")
    logger.info("=" * 60)

    # Run the main loop
    loop = GLib.MainLoop()
    try:
        loop.run()
    except KeyboardInterrupt:
        logger.info("Shutting down Mock ModemManager...")
        sys.exit(0)


if __name__ == "__main__":
    main()
