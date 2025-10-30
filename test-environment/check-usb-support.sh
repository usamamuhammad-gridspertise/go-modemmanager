#!/bin/bash
# Check USB Modem Support Script
# Tests if USB modems can be detected in the current environment

set -e

echo "================================================"
echo "USB Modem Detection Script"
echo "================================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if running in container
if [ -f /.dockerenv ]; then
    echo -e "${CYAN}Running inside Docker container${NC}"
    IN_CONTAINER=true
else
    echo -e "${CYAN}Running on host system${NC}"
    IN_CONTAINER=false
fi
echo ""

# Check operating system
echo -e "${CYAN}Operating System:${NC}"
uname -a
echo ""

# Check for lsusb
echo -e "${CYAN}Checking for USB utilities...${NC}"
if command -v lsusb &> /dev/null; then
    echo -e "${GREEN}✓ lsusb found${NC}"
    HAS_LSUSB=true
else
    echo -e "${RED}✗ lsusb not found${NC}"
    echo "  Install with: apt-get install usbutils"
    HAS_LSUSB=false
fi
echo ""

# List USB devices
if [ "$HAS_LSUSB" = true ]; then
    echo -e "${CYAN}USB Devices:${NC}"
    lsusb || echo -e "${YELLOW}Unable to list USB devices${NC}"
    echo ""
fi

# Check for common modem vendor IDs
echo -e "${CYAN}Checking for known modem vendors...${NC}"
if [ "$HAS_LSUSB" = true ]; then
    MODEM_FOUND=false

    # Common modem vendor IDs
    declare -A VENDORS=(
        ["12d1"]="Huawei"
        ["2c7c"]="Quectel"
        ["1199"]="Sierra Wireless"
        ["19d2"]="ZTE"
        ["05c6"]="Qualcomm"
        ["0bdb"]="Ericsson"
        ["413c"]="Dell/Novatel"
        ["0930"]="Toshiba"
        ["1508"]="Fibocom"
    )

    for vid in "${!VENDORS[@]}"; do
        if lsusb | grep -i "$vid" > /dev/null; then
            echo -e "${GREEN}✓ ${VENDORS[$vid]} modem detected (VID: $vid)${NC}"
            MODEM_FOUND=true
        fi
    done

    if [ "$MODEM_FOUND" = false ]; then
        echo -e "${YELLOW}⚠ No known modem vendors detected${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Cannot check for modems without lsusb${NC}"
fi
echo ""

# Check for serial devices
echo -e "${CYAN}Checking for serial devices...${NC}"
if ls /dev/ttyUSB* 2>/dev/null; then
    echo -e "${GREEN}✓ USB serial devices found${NC}"
elif ls /dev/ttyACM* 2>/dev/null; then
    echo -e "${GREEN}✓ ACM serial devices found${NC}"
elif ls /dev/cdc-wdm* 2>/dev/null; then
    echo -e "${GREEN}✓ CDC-WDM devices found (QMI modems)${NC}"
else
    echo -e "${YELLOW}⚠ No serial devices found${NC}"
    echo "  Modems typically create /dev/ttyUSB*, /dev/ttyACM*, or /dev/cdc-wdm* devices"
fi
echo ""

# Check for ModemManager
echo -e "${CYAN}Checking ModemManager...${NC}"
if command -v mmcli &> /dev/null; then
    echo -e "${GREEN}✓ mmcli found${NC}"
    echo "  Version:"
    mmcli --version | head -n 1
    HAS_MM=true
else
    echo -e "${RED}✗ ModemManager not installed${NC}"
    echo "  Install with: apt-get install modemmanager"
    HAS_MM=false
fi
echo ""

# Check if ModemManager is running
if [ "$HAS_MM" = true ]; then
    echo -e "${CYAN}ModemManager Status:${NC}"
    if pgrep -x "ModemManager" > /dev/null; then
        echo -e "${GREEN}✓ ModemManager is running${NC}"

        echo ""
        echo -e "${CYAN}Detected Modems:${NC}"
        mmcli -L || echo -e "${YELLOW}No modems detected by ModemManager${NC}"
    else
        echo -e "${RED}✗ ModemManager is not running${NC}"
        echo "  Start with: /usr/sbin/ModemManager --debug &"
    fi
fi
echo ""

# Check D-Bus
echo -e "${CYAN}Checking D-Bus...${NC}"
if command -v dbus-send &> /dev/null; then
    echo -e "${GREEN}✓ D-Bus tools found${NC}"

    if pgrep -x "dbus-daemon" > /dev/null; then
        echo -e "${GREEN}✓ D-Bus daemon is running${NC}"

        # Try to query ModemManager via D-Bus
        if dbus-send --system --print-reply --dest=org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.ListNames 2>/dev/null | grep -q "ModemManager1"; then
            echo -e "${GREEN}✓ ModemManager is registered on D-Bus${NC}"
        else
            echo -e "${YELLOW}⚠ ModemManager not found on D-Bus${NC}"
        fi
    else
        echo -e "${RED}✗ D-Bus daemon is not running${NC}"
    fi
else
    echo -e "${RED}✗ D-Bus tools not found${NC}"
fi
echo ""

# Check kernel modules
echo -e "${CYAN}Checking kernel modules...${NC}"
MODULES=("qmi_wwan" "cdc_wdm" "cdc_acm" "option" "usb_wwan" "qcserial")
MODULE_FOUND=false

for mod in "${MODULES[@]}"; do
    if lsmod | grep -q "^$mod"; then
        echo -e "${GREEN}✓ $mod loaded${NC}"
        MODULE_FOUND=true
    fi
done

if [ "$MODULE_FOUND" = false ]; then
    echo -e "${YELLOW}⚠ No modem-related kernel modules loaded${NC}"
    echo "  Common modules: ${MODULES[*]}"
fi
echo ""

# Platform-specific notes
echo -e "${CYAN}Platform Notes:${NC}"
if [ "$IN_CONTAINER" = true ]; then
    echo -e "${YELLOW}⚠ Running in container:${NC}"
    echo "  - USB device access may be limited"
    echo "  - Container needs --privileged flag or specific device mappings"
    echo "  - On macOS, Docker runs in a VM which limits USB passthrough"
else
    # Check if on macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo -e "${YELLOW}⚠ Running on macOS:${NC}"
        echo "  - Native USB modem access not available"
        echo "  - Consider using:"
        echo "    1. Linux VM with USB passthrough (UTM/VirtualBox)"
        echo "    2. Docker container (limited USB support)"
        echo "    3. Remote Linux machine with SSH"
    fi
fi
echo ""

# Summary
echo "================================================"
echo -e "${CYAN}Summary${NC}"
echo "================================================"

if [ "$HAS_MM" = true ] && pgrep -x "ModemManager" > /dev/null; then
    if mmcli -L 2>/dev/null | grep -q "No modems"; then
        echo -e "${YELLOW}⚠ ModemManager is running but no modems detected${NC}"
        echo ""
        echo "Possible reasons:"
        echo "  1. No physical modem connected"
        echo "  2. Modem not in correct USB mode (use usb_modeswitch)"
        echo "  3. Missing kernel modules or drivers"
        echo "  4. USB device not accessible (Docker/VM limitation)"
        echo ""
        echo "Troubleshooting:"
        echo "  - Check 'lsusb' output above for modem hardware"
        echo "  - Try: mmcli -S  (force modem scan)"
        echo "  - Check: journalctl -u ModemManager -f"
    else
        echo -e "${GREEN}✓ ModemManager is operational with modems detected!${NC}"
        echo ""
        echo "You can now test go-modemmanager:"
        echo "  cd /workspace"
        echo "  go test -v"
        echo "  cd examples && go run test_run.go"
    fi
else
    echo -e "${RED}✗ ModemManager is not operational${NC}"
    echo ""
    echo "Setup steps:"
    echo "  1. Install ModemManager: apt-get install modemmanager"
    echo "  2. Start D-Bus: dbus-daemon --system --fork"
    echo "  3. Start ModemManager: /usr/sbin/ModemManager --debug &"
    echo "  4. Run this script again to verify"
fi

echo ""
