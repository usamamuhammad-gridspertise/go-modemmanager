# ModemManager Exporter - Quick Start Guide

Get up and running with the ModemManager Prometheus Exporter in minutes!

## Prerequisites

Before you begin, ensure you have:

1. **ModemManager installed and running**
   ```bash
   # Check if ModemManager is installed
   which mmcli
   
   # Check if ModemManager service is running
   systemctl status ModemManager
   
   # If not running, start it
   sudo systemctl start ModemManager
   ```

2. **A modem detected by ModemManager**
   ```bash
   # List available modems
   mmcli -L
   
   # Get details of first modem
   mmcli -m 0
   ```

3. **Go 1.13 or later** (for building from source)
   ```bash
   go version
   ```

4. **Root or appropriate D-Bus permissions**

## Option 1: Quick Test (5 minutes)

The fastest way to try the exporter:

```bash
# Navigate to the exporter directory
cd cmd/mm-exporter

# Build the exporter
go build -o mm-exporter

# Run it
./mm-exporter
```

The exporter will start on port 9539. Open your browser to:
- http://localhost:9539 - Landing page
- http://localhost:9539/metrics - Prometheus metrics

Test it with curl:
```bash
curl http://localhost:9539/metrics
```

## Option 2: Using Make (Recommended)

```bash
# From the exporter directory
cd exporter

# Build
make build

# Run locally
make run

# Or install system-wide
make install
```

## Option 3: System Service Installation (Production)

For production deployments:

```bash
cd exporter

# Build and install
make install

# Install systemd service
make install-service

# Enable and start
sudo systemctl enable mm-exporter
sudo systemctl start mm-exporter

# Check status
make status

# View logs
make logs
```

## Option 4: Docker Deployment

```bash
cd exporter/examples

# Using docker-compose (includes Prometheus and Grafana)
docker-compose up -d

# Check logs
docker-compose logs -f mm-exporter

# Access services:
# - Exporter: http://localhost:9539
# - Prometheus: http://localhost:9090
# - Grafana: http://localhost:3000 (admin/admin)
```

## Verify Installation

Once running, verify the exporter is working:

```bash
# Check health
curl http://localhost:9539/health

# View metrics
curl http://localhost:9539/metrics | grep modemmanager

# Check specific metrics
curl -s http://localhost:9539/metrics | grep "modemmanager_modem_info"
curl -s http://localhost:9539/metrics | grep "modemmanager_signal_lte"
```

## Configure Prometheus

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'modemmanager'
    static_configs:
      - targets: ['localhost:9539']
    scrape_interval: 15s
```

Reload Prometheus:
```bash
# If using systemd
sudo systemctl reload prometheus

# Or send SIGHUP
kill -HUP $(pidof prometheus)

# Or use API (if --web.enable-lifecycle is set)
curl -X POST http://localhost:9090/-/reload
```

## View Metrics in Prometheus

1. Open Prometheus UI: http://localhost:9090
2. Go to "Graph" tab
3. Try these queries:

```promql
# All ModemManager metrics
{__name__=~"modemmanager_.*"}

# Signal strength
modemmanager_signal_lte_rsrp_dbm

# Connection status
modemmanager_modem_state{state="connected"}

# Signal quality percentage
modemmanager_modem_signal_quality_percent
```

## Common Configuration Options

### Change Port
```bash
./mm-exporter -listen-address=":9090"
```

### Adjust Signal Polling Rate
```bash
# Poll every 10 seconds
./mm-exporter -signal-rate=10s

# Disable signal polling (saves CPU)
./mm-exporter -signal-rate=0s
```

### Custom Metrics Path
```bash
./mm-exporter -metrics-path="/modem-metrics"
```

### All Options Together
```bash
./mm-exporter \
  -listen-address=":9090" \
  -metrics-path="/metrics" \
  -signal-rate=10s
```

## Example Metrics Output

Here's what you'll see:

```
# HELP modemmanager_info ModemManager daemon version information
# TYPE modemmanager_info gauge
modemmanager_info{version="1.18.6"} 1

# HELP modemmanager_modem_info Modem device information
# TYPE modemmanager_modem_info gauge
modemmanager_modem_info{device="cdc-wdm0",device_id="...",equipment_id="...",manufacturer="Quectel",model="EC25",plugin="quectel",primary_port="cdc-wdm0",revision="EC25EFAR06A03M4G"} 1

# HELP modemmanager_modem_signal_quality_percent Signal quality as a percentage (0-100)
# TYPE modemmanager_modem_signal_quality_percent gauge
modemmanager_modem_signal_quality_percent{device_id="..."} 80

# HELP modemmanager_signal_lte_rsrp_dbm LTE RSRP (Reference Signal Received Power) in dBm
# TYPE modemmanager_signal_lte_rsrp_dbm gauge
modemmanager_signal_lte_rsrp_dbm{device_id="..."} -85

# HELP modemmanager_modem_state Current modem state (enumeration)
# TYPE modemmanager_modem_state gauge
modemmanager_modem_state{device_id="...",state="connected"} 1
```

## Troubleshooting

### No modems found
```bash
# Check ModemManager
systemctl status ModemManager
mmcli -L

# Scan for new devices
sudo mmcli -S
```

### Permission denied errors
```bash
# The exporter needs D-Bus access, typically requires root
sudo ./mm-exporter

# Or add user to appropriate groups (distribution-dependent)
```

### Signal metrics missing
```bash
# Enable signal interface on modem first
mmcli -m 0 --signal-setup=5

# Or let the exporter do it automatically (default)
./mm-exporter -signal-rate=5s
```

### High CPU usage
```bash
# Reduce signal polling frequency
./mm-exporter -signal-rate=30s

# Or disable it completely
./mm-exporter -signal-rate=0s
```

### Exporter not responding
```bash
# Check if it's running
ps aux | grep mm-exporter

# Check system logs
journalctl -u mm-exporter -f

# Check if port is in use
netstat -tlnp | grep 9539
```

## Next Steps

1. **Set up alerting**: Create alert rules for signal degradation, disconnections, etc.
2. **Create Grafana dashboard**: Visualize signal strength, connection status, etc.
3. **Monitor multiple devices**: The exporter automatically handles multiple modems
4. **Integrate with your monitoring stack**: Export to your existing Prometheus/Grafana setup

## Useful Commands

```bash
# Build
make build

# Run locally for testing
make run

# Install system-wide
make install

# Install as service
make install-service

# Check ModemManager status
make check-modemmanager

# Test metrics endpoint
make test-metrics

# View logs
make logs

# Uninstall everything
make uninstall

# Get help
make help
```

## Performance Tips

1. **Signal polling**: Higher rates = more CPU. Start with 15-30s for production.
2. **Scrape interval**: Match your Prometheus scrape interval to signal rate.
3. **Multiple modems**: The exporter handles them efficiently in a single process.
4. **Resource limits**: Use systemd limits if needed (see service file).

## Security Considerations

1. **Firewall**: Consider restricting access to port 9539
   ```bash
   sudo ufw allow from 10.0.0.0/8 to any port 9539
   ```

2. **TLS/Authentication**: Use a reverse proxy (nginx, traefik) for production
   ```nginx
   location /metrics {
       proxy_pass http://localhost:9539;
       allow 10.0.0.0/8;
       deny all;
   }
   ```

3. **Read-only access**: The exporter only reads data, never modifies modem config

## Getting Help

- Check logs: `journalctl -u mm-exporter -f`
- ModemManager status: `mmcli -m 0`
- Test manually: `curl http://localhost:9539/metrics`
- See full README: [README.md](README.md)
- Open an issue on GitHub

## Summary

You should now have:
- ✅ Exporter built and running
- ✅ Metrics exposed on port 9539
- ✅ Prometheus scraping the metrics
- ✅ Basic understanding of available metrics

**Congratulations!** Your ModemManager exporter is ready to use.

For more advanced configuration and all available metrics, see [README.md](README.md).