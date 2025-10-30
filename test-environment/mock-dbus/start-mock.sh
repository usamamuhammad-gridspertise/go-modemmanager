#!/bin/bash
# Start Mock ModemManager D-Bus Service
# This script stops the real ModemManager and starts the mock version

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}================================================${NC}"
echo -e "${CYAN}Starting Mock ModemManager${NC}"
echo -e "${CYAN}================================================${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${YELLOW}Note: May need root privileges for system bus access${NC}"
    echo -e "${YELLOW}If this fails, try: sudo $0${NC}"
    echo ""
fi

# Check if Python and dependencies are available
echo -e "${CYAN}Checking dependencies...${NC}"

if ! command -v python3 &> /dev/null; then
    echo -e "${RED}✗ python3 not found${NC}"
    echo "Install with: apt-get install python3"
    exit 1
fi
echo -e "${GREEN}✓ python3 found${NC}"

# Check for required Python modules
python3 -c "import dbus" 2>/dev/null
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}⚠ dbus-python not found${NC}"
    echo "Installing python3-dbus..."
    apt-get update && apt-get install -y python3-dbus || pip3 install dbus-python
fi

python3 -c "from gi.repository import GLib" 2>/dev/null
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}⚠ pygobject not found${NC}"
    echo "Installing python3-gi..."
    apt-get update && apt-get install -y python3-gi || pip3 install pygobject
fi

echo -e "${GREEN}✓ All dependencies available${NC}"
echo ""

# Check if D-Bus is running
echo -e "${CYAN}Checking D-Bus...${NC}"
if ! pgrep -x "dbus-daemon" > /dev/null; then
    echo -e "${YELLOW}⚠ D-Bus daemon not running${NC}"
    echo "Starting D-Bus..."
    mkdir -p /var/run/dbus
    rm -f /var/run/dbus/pid
    dbus-daemon --system --fork
    sleep 2
fi
echo -e "${GREEN}✓ D-Bus is running${NC}"
echo ""

# Stop real ModemManager if running
echo -e "${CYAN}Checking for real ModemManager...${NC}"
if pgrep -x "ModemManager" > /dev/null; then
    echo -e "${YELLOW}⚠ Real ModemManager is running${NC}"
    echo "Stopping ModemManager..."
    systemctl stop ModemManager 2>/dev/null || killall ModemManager 2>/dev/null || true
    sleep 2

    if pgrep -x "ModemManager" > /dev/null; then
        echo -e "${RED}✗ Failed to stop ModemManager${NC}"
        echo "Try manually: systemctl stop ModemManager"
        exit 1
    fi
    echo -e "${GREEN}✓ ModemManager stopped${NC}"
else
    echo -e "${GREEN}✓ ModemManager not running${NC}"
fi
echo ""

# Start mock ModemManager
echo -e "${CYAN}Starting Mock ModemManager...${NC}"
echo ""

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Run the mock service
echo -e "${GREEN}Mock ModemManager is starting...${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""

python3 mock_modemmanager.py

# Cleanup on exit
echo ""
echo -e "${CYAN}Cleaning up...${NC}"
echo -e "${YELLOW}To restart real ModemManager:${NC}"
echo "  systemctl start ModemManager"
echo ""
