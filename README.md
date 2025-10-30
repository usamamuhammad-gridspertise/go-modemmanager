![Alt Go-ModemManager](./go-modemmanager.png)

[![GoDoc](https://godoc.org/github.com/maltegrosse/go-modemmanager?status.svg)](https://pkg.go.dev/github.com/maltegrosse/go-modemmanager)
[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)
![Go](https://github.com/maltegrosse/go-modemmanager/workflows/Go/badge.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/maltegrosse/go-modemmanager)](https://goreportcard.com/report/github.com/maltegrosse/go-modemmanager)

Go D-Bus bindings for ModemManager


Additional information: [ModemManager D-Bus Specs](https://www.freedesktop.org/software/ModemManager/api/1.12.0/ref-dbus.html)

Tested with [ModemManager - Version 1.12.8](https://gitlab.freedesktop.org/mobile-broadband/ModemManager), Go 1.13, on `Debian Buster (armv7)` with `Kernel 5.4.x` and `libqmi 1.24.6`.

Test hardware: [SolidRun Hummingboard Edge](https://www.solid-run.com/nxp-family/hummingboard/)   and a `Quectel EC25 - EC25EFA` mini pcie modem.

## Features

### Library
- Complete D-Bus bindings for ModemManager
- Support for all major modem interfaces
- Simple and complex operations
- Example code included

### Prometheus Exporter
A production-ready Prometheus exporter is now available! Monitor your cellular modems with 40+ metrics covering:
- Signal strength (LTE, UMTS, GSM, CDMA, EVDO)
- Connection status and bearer information
- Network registration and operator details
- SIM card information
- Location data (GPS)
- SMS messaging statistics

**Quick Start:**
```bash
cd exporter
make build
make run
# Visit http://localhost:9539/metrics
```

See [exporter/README.md](exporter/README.md) for full documentation and [EXPORTER_SUMMARY.md](EXPORTER_SUMMARY.md) for complete feature list.

## Notes
 ModemManager works great together with GeoClue. A dbus wrapper can be found [here](https://github.com/maltegrosse/go-geoclue2).

A NetworkManager dbus wrapper in golang can be found [here](https://github.com/Wifx/gonetworkmanager).

## Status
Some methods/properties are untested as they are not supported by my modem/lack of how to use them. See `todo` tags in the code.

## Installation

This packages requires Go 1.13 (for the dbus lib). If you installed it and set up your GOPATH, just run:

`go get -u github.com/maltegrosse/go-modemmanager`

## Usage

### Library Usage
You can find some examples in the [examples](examples) directory.

### Prometheus Exporter Usage
```bash
# Build the exporter
cd cmd/mm-exporter
go build -o mm-exporter

# Run with default settings
./mm-exporter

# Or use the Makefile
cd exporter
make install
make install-service
```

See [exporter/QUICKSTART.md](exporter/QUICKSTART.md) for detailed instructions.

## Limitations
Not all interfaces, methods and properties are supported in QMI or AT mode. In addition, not all methods and properties are supported by every modem.
A brief overview of the availability of each interface by using Quectel EC-25:

| Interface     | QMI   | AT    |
|---------------|-------|-------|
| ModemManager1 | true  | true  |
| Modem         | true  | true  |
| Simple        | true  | true  |
| Modem3gpp     | true  | true  |
| Ussd          | false | true  |
| ModemCdma     | false | false |
| Messaging     | true  | false |
| Location      | true  | true  |
| Time          | true  | true  |
| Firmware      | true  | true  |
| Signal        | true  | false |
| Oma           | false | false |
| Bearer        | true  | true  |
| Sim           | true  | true  |
| SMS           | true  | true  |
| Call          | true  | true  |

## License
**[MIT license](http://opensource.org/licenses/mit-license.php)**

Copyright 2020 © Malte Grosse.

Other:
- [ModemManager Logo under GPLv2+](https://gitlab.freedesktop.org/mobile-broadband/ModemManager/-/tree/master/data)

- [GoLang Logo under Creative Commons Attribution 3.0](https://blog.golang.org/go-brand)
