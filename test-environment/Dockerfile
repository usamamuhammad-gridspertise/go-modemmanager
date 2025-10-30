# ModemManager Testing Environment
FROM ubuntu:22.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install required packages
RUN apt-get update && apt-get install -y \
    modemmanager \
    libmm-glib-dev \
    dbus \
    usbutils \
    usb-modeswitch \
    libqmi-utils \
    libmbim-utils \
    golang-1.21 \
    git \
    vim \
    curl \
    python3 \
    python3-pip \
    policykit-1 \
    socat \
    && rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV PATH="/usr/lib/go-1.21/bin:${PATH}"
ENV GOPATH="/root/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Create directories
RUN mkdir -p /var/run/dbus \
    && mkdir -p /workspace \
    && mkdir -p /root/go/src

# Copy D-Bus configuration for ModemManager
RUN mkdir -p /etc/dbus-1/system.d

# Create a custom D-Bus policy to allow root access (for testing)
RUN echo '<!DOCTYPE busconfig PUBLIC "-//freedesktop//DTD D-BUS Bus Configuration 1.0//EN" \
 "http://www.freedesktop.org/standards/dbus/1.0/busconfig.dtd"> \
<busconfig> \
  <policy user="root"> \
    <allow own="org.freedesktop.ModemManager1"/> \
    <allow send_destination="org.freedesktop.ModemManager1"/> \
    <allow send_interface="org.freedesktop.ModemManager1"/> \
  </policy> \
  <policy context="default"> \
    <allow send_destination="org.freedesktop.ModemManager1"/> \
  </policy> \
</busconfig>' > /etc/dbus-1/system.d/org.freedesktop.ModemManager1.conf

# Create startup script
RUN echo '#!/bin/bash\n\
set -e\n\
\n\
echo "Starting D-Bus system daemon..."\n\
mkdir -p /var/run/dbus\n\
rm -f /var/run/dbus/pid\n\
dbus-daemon --system --fork\n\
\n\
echo "Waiting for D-Bus to be ready..."\n\
sleep 2\n\
\n\
echo "Starting ModemManager..."\n\
/usr/sbin/ModemManager --debug &\n\
\n\
echo "Waiting for ModemManager to be ready..."\n\
sleep 3\n\
\n\
echo ""\n\
echo "================================"\n\
echo "ModemManager Test Environment"\n\
echo "================================"\n\
echo ""\n\
echo "ModemManager version:"\n\
mmcli --version || echo "mmcli not available"\n\
echo ""\n\
echo "Available modems:"\n\
mmcli -L || echo "No modems detected"\n\
echo ""\n\
echo "To test go-modemmanager:"\n\
echo "  cd /workspace"\n\
echo "  go test -v"\n\
echo ""\n\
echo "To run examples:"\n\
echo "  cd /workspace/examples"\n\
echo "  go run test_run.go"\n\
echo ""\n\
echo "Note: Without a physical modem, some tests will fail."\n\
echo "This environment is useful for:"\n\
echo "  - Code compilation testing"\n\
echo "  - D-Bus interface validation"\n\
echo "  - Mock-based unit tests"\n\
echo ""\n\
\n\
exec "$@"\n\
' > /usr/local/bin/start-modemmanager.sh && chmod +x /usr/local/bin/start-modemmanager.sh

WORKDIR /workspace

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/start-modemmanager.sh"]
CMD ["/bin/bash"]
