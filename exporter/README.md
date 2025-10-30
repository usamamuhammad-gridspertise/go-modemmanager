# ModemManager Prometheus Exporter

A Prometheus exporter for ModemManager that exposes detailed metrics about cellular modems and their connections.

## Overview

This exporter connects to ModemManager via D-Bus and exports comprehensive metrics about:
- Modem status and configuration
- Signal strength (LTE, UMTS, GSM, CDMA, EVDO)
- Bearer/connection information
- SIM card details
- 3GPP network registration
- Location data (if enabled)
- SMS messaging statistics

## Features

- **Comprehensive Metrics**: Exports 40+ different metric types covering all aspects of modem operation
- **Multi-Technology Support**: LTE, UMTS, GSM, CDMA, and EVDO signal metrics
- **Auto-Discovery**: Automatically discovers and monitors all available modems
- **Signal Polling**: Configurable refresh rate for extended signal quality information
- **Production Ready**: Includes health checks, graceful shutdown, and proper error handling
- **Standard Prometheus Integration**: Compatible with all Prometheus-based monitoring stacks

## Installation

### From Source

```bash
cd cmd/mm-exporter
go build -o mm-exporter
sudo cp mm-exporter /usr/local/bin/
```

### As a Service

Create a systemd service file at `/etc/systemd/system/mm-exporter.service`:

```ini
[Unit]
Description=ModemManager Prometheus Exporter
After=network.target ModemManager.service
Requires=ModemManager.service

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/mm-exporter
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable mm-exporter
sudo systemctl start mm-exporter
```

## Usage

### Basic Usage

```bash
# Start with default settings (port 9539, signal refresh every 5 seconds)
./mm-exporter

# Custom port and metrics path
./mm-exporter -listen-address=":9090" -metrics-path="/modem-metrics"

# Adjust signal refresh rate (or disable with 0s)
./mm-exporter -signal-rate=10s

# Show version
./mm-exporter -version
```

### Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-listen-address` | `:9539` | Address on which to expose metrics and web interface |
| `-metrics-path` | `/metrics` | Path under which to expose metrics |
| `-signal-rate` | `5s` | How frequently to poll modems for extended signal data (0 to disable) |
| `-version` | `false` | Show version information and exit |

### Endpoints

- `/` - Landing page with exporter information
- `/metrics` - Prometheus metrics endpoint
- `/health` - Health check endpoint (returns 200 OK)

## Exported Metrics

### ModemManager Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_info` | Gauge | `version` | ModemManager daemon version |

### Modem Information Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_modem_info` | Gauge | `device_id`, `manufacturer`, `model`, `revision`, `equipment_id`, `device`, `plugin`, `primary_port` | Modem device information |
| `modemmanager_modem_state` | Gauge | `device_id`, `state` | Current modem state (1 = active state) |
| `modemmanager_modem_power_state` | Gauge | `device_id`, `state` | Current power state (1 = active state) |
| `modemmanager_modem_signal_quality_percent` | Gauge | `device_id` | Signal quality percentage (0-100) |
| `modemmanager_modem_access_technology` | Gauge | `device_id`, `technology` | Current access technology (1 = active) |
| `modemmanager_modem_unlock_required` | Gauge | `device_id` | Unlock requirement type (0 = none) |
| `modemmanager_modem_max_bearers` | Gauge | `device_id` | Maximum bearers supported |
| `modemmanager_modem_max_active_bearers` | Gauge | `device_id` | Maximum active bearers supported |

### Signal Strength Metrics

#### LTE Signals
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_signal_lte_rssi_dbm` | Gauge | `device_id` | LTE RSSI in dBm |
| `modemmanager_signal_lte_rsrq_db` | Gauge | `device_id` | LTE RSRQ in dB |
| `modemmanager_signal_lte_rsrp_dbm` | Gauge | `device_id` | LTE RSRP in dBm |
| `modemmanager_signal_lte_snr_db` | Gauge | `device_id` | LTE SNR in dB |

#### UMTS Signals
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_signal_umts_rssi_dbm` | Gauge | `device_id` | UMTS RSSI in dBm |
| `modemmanager_signal_umts_ecio_db` | Gauge | `device_id` | UMTS Ec/Io in dB |
| `modemmanager_signal_umts_rscp_dbm` | Gauge | `device_id` | UMTS RSCP in dBm |

#### GSM Signals
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_signal_gsm_rssi_dbm` | Gauge | `device_id` | GSM RSSI in dBm |

#### CDMA Signals
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_signal_cdma_rssi_dbm` | Gauge | `device_id` | CDMA RSSI in dBm |
| `modemmanager_signal_cdma_ecio_db` | Gauge | `device_id` | CDMA Ec/Io in dB |

#### EVDO Signals
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_signal_evdo_rssi_dbm` | Gauge | `device_id` | EVDO RSSI in dBm |
| `modemmanager_signal_evdo_ecio_db` | Gauge | `device_id` | EVDO Ec/Io in dB |
| `modemmanager_signal_evdo_sinr_db` | Gauge | `device_id` | EVDO SINR in dB |
| `modemmanager_signal_evdo_io_dbm` | Gauge | `device_id` | EVDO Io in dBm |

### Bearer Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_bearer_info` | Gauge | `device_id`, `bearer_path`, `interface`, `ip_method`, `ip_address` | Bearer information |
| `modemmanager_bearer_connected` | Gauge | `device_id`, `bearer_path` | Bearer connection status |

### SIM Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_sim_info` | Gauge | `device_id`, `sim_path`, `imsi`, `operator_name` | SIM card information |

### 3GPP Network Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_modem_3gpp_registration_state` | Gauge | `device_id`, `state` | 3GPP registration state (1 = active) |
| `modemmanager_modem_3gpp_operator_code` | Gauge | `device_id`, `operator_code` | Operator code (MCC+MNC) |
| `modemmanager_modem_3gpp_operator_name` | Gauge | `device_id`, `operator_name` | Operator name |

### Messaging Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_messaging_supported` | Gauge | `device_id` | Whether messaging is supported |
| `modemmanager_messaging_sms_count` | Gauge | `device_id` | Number of stored SMS messages |

### Location Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_location_enabled` | Gauge | `device_id` | Whether location services are enabled |
| `modemmanager_location_latitude_degrees` | Gauge | `device_id` | Current latitude |
| `modemmanager_location_longitude_degrees` | Gauge | `device_id` | Current longitude |
| `modemmanager_location_altitude_meters` | Gauge | `device_id` | Current altitude |

### Scrape Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `modemmanager_scrape_duration_seconds` | Gauge | - | Duration of the scrape |
| `modemmanager_scrape_success` | Gauge | - | Whether scrape was successful |
| `modemmanager_scrape_errors_total` | Counter | - | Total scrape errors |

## Prometheus Configuration

Add the exporter to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'modemmanager'
    static_configs:
      - targets: ['localhost:9539']
    scrape_interval: 15s
```

## Example Queries

### Signal Strength
```promql
# LTE RSRP (signal power)
modemmanager_signal_lte_rsrp_dbm

# Signal quality percentage
modemmanager_modem_signal_quality_percent
```

### Connection Status
```promql
# Modems in connected state
modemmanager_modem_state{state="connected"}

# Active bearers
modemmanager_bearer_connected == 1
```

### Network Registration
```promql
# Registration state
modemmanager_modem_3gpp_registration_state{state="home"}
```

### Alerting Examples

```yaml
groups:
  - name: modemmanager
    rules:
      - alert: ModemSignalWeak
        expr: modemmanager_signal_lte_rsrp_dbm < -110
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Weak LTE signal on {{ $labels.device_id }}"
          description: "RSRP is {{ $value }} dBm"

      - alert: ModemDisconnected
        expr: modemmanager_bearer_connected == 0
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Modem {{ $labels.device_id }} disconnected"

      - alert: ModemUnlockRequired
        expr: modemmanager_modem_unlock_required > 0
        labels:
          severity: warning
        annotations:
          summary: "Modem {{ $labels.device_id }} requires unlock"
```

## Grafana Dashboard

A sample Grafana dashboard JSON is available in the `examples/` directory (coming soon).

Key panels to include:
- Signal strength over time (RSRP, RSRQ, SNR)
- Connection status
- Network operator and technology
- Bearer information
- Data transfer metrics (if available from other sources)

## Requirements

- ModemManager 1.0 or later
- Go 1.13 or later (for building)
- D-Bus system bus access
- Root or appropriate permissions to access ModemManager

## Troubleshooting

### No Metrics Appearing

1. **Check ModemManager is running:**
   ```bash
   systemctl status ModemManager
   ```

2. **Verify modems are detected:**
   ```bash
   mmcli -L
   ```

3. **Check D-Bus permissions:**
   The exporter needs permission to access the ModemManager D-Bus interface.

### Signal Metrics Missing

Some signal metrics require the Signal interface to be available, which may depend on:
- Modem capabilities
- Current modem state (must be registered/connected)
- Protocol in use (QMI vs AT commands)

### High CPU Usage

If signal polling causes high CPU usage, increase the `-signal-rate` or disable it with `-signal-rate=0s`.

## Architecture

The exporter uses the Prometheus client library's collector pattern:

1. **Exporter struct**: Implements `prometheus.Collector` interface
2. **Describe()**: Registers metric descriptors
3. **Collect()**: Gathers metrics on each scrape
4. **Main loop**: HTTP server exposes metrics endpoint

The exporter connects to ModemManager via D-Bus and queries modem properties on each Prometheus scrape.

## Development

### Building

```bash
cd cmd/mm-exporter
go build -o mm-exporter
```

### Testing

```bash
# Run the exporter
./mm-exporter

# Query metrics
curl http://localhost:9539/metrics

# Health check
curl http://localhost:9539/health
```

### Adding New Metrics

1. Add metric descriptor to `Exporter` struct in `handler.go`
2. Initialize the descriptor in `NewExporter()`
3. Add to `Describe()` method
4. Add collection logic in `Collect()` or helper methods

## License

MIT License - See LICENSE.md in the root directory

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Related Projects

- [go-modemmanager](https://github.com/maltegrosse/go-modemmanager) - The underlying Go library
- [ModemManager](https://www.freedesktop.org/wiki/Software/ModemManager/) - The modem management daemon
- [node_exporter](https://github.com/prometheus/node_exporter) - Complementary system metrics

## Support

For issues and questions:
- Open an issue on GitHub
- Check ModemManager logs: `journalctl -u ModemManager`
- Verify modem status: `mmcli -m 0 --output-json`
