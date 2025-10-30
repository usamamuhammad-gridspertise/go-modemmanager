package exporter

import (
	"fmt"
	"log"
	"time"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "modemmanager"
)

// Exporter collects ModemManager metrics and exports them using
// the prometheus client library.
type Exporter struct {
	mm modemmanager.ModemManager

	// ModemManager info
	mmInfo *prometheus.Desc

	// Modem info
	modemInfo             *prometheus.Desc
	modemState            *prometheus.Desc
	modemPowerState       *prometheus.Desc
	modemSignalQuality    *prometheus.Desc
	modemAccessTech       *prometheus.Desc
	modemUnlockRequired   *prometheus.Desc
	modemMaxBearers       *prometheus.Desc
	modemMaxActiveBearers *prometheus.Desc

	// Signal metrics (LTE)
	signalLteRssi *prometheus.Desc
	signalLteRsrq *prometheus.Desc
	signalLteRsrp *prometheus.Desc
	signalLteSnr  *prometheus.Desc

	// Signal metrics (UMTS)
	signalUmtsRssi *prometheus.Desc
	signalUmtsEcio *prometheus.Desc
	signalUmtsRscp *prometheus.Desc

	// Signal metrics (GSM)
	signalGsmRssi *prometheus.Desc

	// Signal metrics (CDMA)
	signalCdmaRssi *prometheus.Desc
	signalCdmaEcio *prometheus.Desc

	// Signal metrics (EVDO)
	signalEvdoRssi *prometheus.Desc
	signalEvdoEcio *prometheus.Desc
	signalEvdoSinr *prometheus.Desc
	signalEvdoIo   *prometheus.Desc

	// Bearer metrics
	bearerInfo      *prometheus.Desc
	bearerConnected *prometheus.Desc

	// SIM metrics
	simInfo *prometheus.Desc

	// 3GPP metrics
	modem3gppRegistrationState *prometheus.Desc
	modem3gppOperatorCode      *prometheus.Desc
	modem3gppOperatorName      *prometheus.Desc

	// Messaging metrics
	messagingSupported *prometheus.Desc
	smsCount           *prometheus.Desc

	// Location metrics
	locationEnabled   *prometheus.Desc
	locationLatitude  *prometheus.Desc
	locationLongitude *prometheus.Desc
	locationAltitude  *prometheus.Desc

	// Scrape metrics
	scrapeDuration *prometheus.Desc
	scrapeSuccess  *prometheus.Desc
	scrapeErrors   *prometheus.Desc
}

// NewExporter returns a new ModemManager exporter.
func NewExporter(mm modemmanager.ModemManager) *Exporter {
	return &Exporter{
		mm: mm,

		// ModemManager info
		mmInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "info"),
			"ModemManager daemon version information",
			[]string{"version"},
			nil,
		),

		// Modem info
		modemInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "info"),
			"Modem device information",
			[]string{"device_id", "manufacturer", "model", "revision", "equipment_id", "device", "plugin", "primary_port"},
			nil,
		),
		modemState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "state"),
			"Current modem state (enumeration)",
			[]string{"device_id", "state"},
			nil,
		),
		modemPowerState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "power_state"),
			"Current modem power state (enumeration)",
			[]string{"device_id", "state"},
			nil,
		),
		modemSignalQuality: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "signal_quality_percent"),
			"Signal quality as a percentage (0-100)",
			[]string{"device_id"},
			nil,
		),
		modemAccessTech: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "access_technology"),
			"Current access technology (enumeration)",
			[]string{"device_id", "technology"},
			nil,
		),
		modemUnlockRequired: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "unlock_required"),
			"Type of unlock required (0 = none)",
			[]string{"device_id"},
			nil,
		),
		modemMaxBearers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "max_bearers"),
			"Maximum number of bearers supported",
			[]string{"device_id"},
			nil,
		),
		modemMaxActiveBearers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem", "max_active_bearers"),
			"Maximum number of active bearers supported",
			[]string{"device_id"},
			nil,
		),

		// Signal metrics (LTE)
		signalLteRssi: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "lte_rssi_dbm"),
			"LTE RSSI (Received Signal Strength Indication) in dBm",
			[]string{"device_id"},
			nil,
		),
		signalLteRsrq: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "lte_rsrq_db"),
			"LTE RSRQ (Reference Signal Received Quality) in dB",
			[]string{"device_id"},
			nil,
		),
		signalLteRsrp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "lte_rsrp_dbm"),
			"LTE RSRP (Reference Signal Received Power) in dBm",
			[]string{"device_id"},
			nil,
		),
		signalLteSnr: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "lte_snr_db"),
			"LTE SNR (Signal-to-Noise Ratio) in dB",
			[]string{"device_id"},
			nil,
		),

		// Signal metrics (UMTS)
		signalUmtsRssi: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "umts_rssi_dbm"),
			"UMTS RSSI in dBm",
			[]string{"device_id"},
			nil,
		),
		signalUmtsEcio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "umts_ecio_db"),
			"UMTS Ec/Io in dB",
			[]string{"device_id"},
			nil,
		),
		signalUmtsRscp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "umts_rscp_dbm"),
			"UMTS RSCP (Received Signal Code Power) in dBm",
			[]string{"device_id"},
			nil,
		),

		// Signal metrics (GSM)
		signalGsmRssi: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "gsm_rssi_dbm"),
			"GSM RSSI in dBm",
			[]string{"device_id"},
			nil,
		),

		// Signal metrics (CDMA)
		signalCdmaRssi: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "cdma_rssi_dbm"),
			"CDMA RSSI in dBm",
			[]string{"device_id"},
			nil,
		),
		signalCdmaEcio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "cdma_ecio_db"),
			"CDMA Ec/Io in dB",
			[]string{"device_id"},
			nil,
		),

		// Signal metrics (EVDO)
		signalEvdoRssi: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "evdo_rssi_dbm"),
			"EVDO RSSI in dBm",
			[]string{"device_id"},
			nil,
		),
		signalEvdoEcio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "evdo_ecio_db"),
			"EVDO Ec/Io in dB",
			[]string{"device_id"},
			nil,
		),
		signalEvdoSinr: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "evdo_sinr_db"),
			"EVDO SINR in dB",
			[]string{"device_id"},
			nil,
		),
		signalEvdoIo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "evdo_io_dbm"),
			"EVDO Io in dBm",
			[]string{"device_id"},
			nil,
		),

		// Bearer metrics
		bearerInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bearer", "info"),
			"Bearer information",
			[]string{"device_id", "bearer_path", "interface", "ip_method", "ip_address"},
			nil,
		),
		bearerConnected: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "bearer", "connected"),
			"Bearer connection status (1 = connected, 0 = disconnected)",
			[]string{"device_id", "bearer_path"},
			nil,
		),

		// SIM metrics
		simInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sim", "info"),
			"SIM card information",
			[]string{"device_id", "sim_path", "imsi", "operator_name"},
			nil,
		),

		// 3GPP metrics
		modem3gppRegistrationState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem_3gpp", "registration_state"),
			"3GPP registration state (enumeration)",
			[]string{"device_id", "state"},
			nil,
		),
		modem3gppOperatorCode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem_3gpp", "operator_code"),
			"3GPP operator code (MCC+MNC)",
			[]string{"device_id", "operator_code"},
			nil,
		),
		modem3gppOperatorName: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "modem_3gpp", "operator_name"),
			"3GPP operator name",
			[]string{"device_id", "operator_name"},
			nil,
		),

		// Messaging metrics
		messagingSupported: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "messaging", "supported"),
			"Whether messaging is supported (1 = yes, 0 = no)",
			[]string{"device_id"},
			nil,
		),
		smsCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "messaging", "sms_count"),
			"Number of SMS messages stored on the modem",
			[]string{"device_id"},
			nil,
		),

		// Location metrics
		locationEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "location", "enabled"),
			"Whether location services are enabled (1 = yes, 0 = no)",
			[]string{"device_id"},
			nil,
		),
		locationLatitude: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "location", "latitude_degrees"),
			"Current latitude in degrees",
			[]string{"device_id"},
			nil,
		),
		locationLongitude: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "location", "longitude_degrees"),
			"Current longitude in degrees",
			[]string{"device_id"},
			nil,
		),
		locationAltitude: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "location", "altitude_meters"),
			"Current altitude in meters",
			[]string{"device_id"},
			nil,
		),

		// Scrape metrics
		scrapeDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "duration_seconds"),
			"Duration of the scrape in seconds",
			nil,
			nil,
		),
		scrapeSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "success"),
			"Whether the scrape was successful (1 = yes, 0 = no)",
			nil,
			nil,
		),
		scrapeErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "errors_total"),
			"Total number of errors during scrape",
			nil,
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.mmInfo
	ch <- e.modemInfo
	ch <- e.modemState
	ch <- e.modemPowerState
	ch <- e.modemSignalQuality
	ch <- e.modemAccessTech
	ch <- e.modemUnlockRequired
	ch <- e.modemMaxBearers
	ch <- e.modemMaxActiveBearers
	ch <- e.signalLteRssi
	ch <- e.signalLteRsrq
	ch <- e.signalLteRsrp
	ch <- e.signalLteSnr
	ch <- e.signalUmtsRssi
	ch <- e.signalUmtsEcio
	ch <- e.signalUmtsRscp
	ch <- e.signalGsmRssi
	ch <- e.signalCdmaRssi
	ch <- e.signalCdmaEcio
	ch <- e.signalEvdoRssi
	ch <- e.signalEvdoEcio
	ch <- e.signalEvdoSinr
	ch <- e.signalEvdoIo
	ch <- e.bearerInfo
	ch <- e.bearerConnected
	ch <- e.simInfo
	ch <- e.modem3gppRegistrationState
	ch <- e.modem3gppOperatorCode
	ch <- e.modem3gppOperatorName
	ch <- e.messagingSupported
	ch <- e.smsCount
	ch <- e.locationEnabled
	ch <- e.locationLatitude
	ch <- e.locationLongitude
	ch <- e.locationAltitude
	ch <- e.scrapeDuration
	ch <- e.scrapeSuccess
	ch <- e.scrapeErrors
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	errorCount := 0
	success := 1.0

	// Collect ModemManager version
	if version, err := e.mm.GetVersion(); err == nil {
		ch <- prometheus.MustNewConstMetric(e.mmInfo, prometheus.GaugeValue, 1.0, version)
	} else {
		log.Printf("Error getting ModemManager version: %v", err)
		errorCount++
	}

	// Collect modem metrics
	modems, err := e.mm.GetModems()
	if err != nil {
		log.Printf("Error getting modems: %v", err)
		errorCount++
		success = 0.0
	} else {
		for _, modem := range modems {
			if err := e.collectModemMetrics(ch, modem); err != nil {
				log.Printf("Error collecting metrics for modem: %v", err)
				errorCount++
			}
		}
	}

	// Export scrape metrics
	duration := time.Since(start).Seconds()
	ch <- prometheus.MustNewConstMetric(e.scrapeDuration, prometheus.GaugeValue, duration)
	ch <- prometheus.MustNewConstMetric(e.scrapeSuccess, prometheus.GaugeValue, success)
	ch <- prometheus.MustNewConstMetric(e.scrapeErrors, prometheus.CounterValue, float64(errorCount))
}

func (e *Exporter) collectModemMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem) error {
	deviceID, err := modem.GetDeviceIdentifier()
	if err != nil {
		return fmt.Errorf("failed to get device identifier: %w", err)
	}

	// Collect basic modem info
	e.collectModemInfo(ch, modem, deviceID)

	// Collect modem state
	e.collectModemState(ch, modem, deviceID)

	// Collect signal metrics
	e.collectSignalMetrics(ch, modem, deviceID)

	// Collect bearer metrics
	e.collectBearerMetrics(ch, modem, deviceID)

	// Collect SIM metrics
	e.collectSIMMetrics(ch, modem, deviceID)

	// Collect 3GPP metrics
	e.collect3GPPMetrics(ch, modem, deviceID)

	// Collect messaging metrics
	e.collectMessagingMetrics(ch, modem, deviceID)

	// Collect location metrics
	e.collectLocationMetrics(ch, modem, deviceID)

	return nil
}

func (e *Exporter) collectModemInfo(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	manufacturer, _ := modem.GetManufacturer()
	model, _ := modem.GetModel()
	revision, _ := modem.GetRevision()
	equipmentID, _ := modem.GetEquipmentIdentifier()
	device, _ := modem.GetDevice()
	plugin, _ := modem.GetPlugin()
	primaryPort, _ := modem.GetPrimaryPort()

	ch <- prometheus.MustNewConstMetric(
		e.modemInfo,
		prometheus.GaugeValue,
		1.0,
		deviceID, manufacturer, model, revision, equipmentID, device, plugin, primaryPort,
	)

	// Max bearers
	if maxBearers, err := modem.GetMaxBearers(); err == nil {
		ch <- prometheus.MustNewConstMetric(e.modemMaxBearers, prometheus.GaugeValue, float64(maxBearers), deviceID)
	}

	if maxActiveBearers, err := modem.GetMaxActiveBearers(); err == nil {
		ch <- prometheus.MustNewConstMetric(e.modemMaxActiveBearers, prometheus.GaugeValue, float64(maxActiveBearers), deviceID)
	}
}

func (e *Exporter) collectModemState(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	// Modem state
	if state, err := modem.GetState(); err == nil {
		stateStr := stateToString(state)
		ch <- prometheus.MustNewConstMetric(e.modemState, prometheus.GaugeValue, 1.0, deviceID, stateStr)
	}

	// Power state
	if powerState, err := modem.GetPowerState(); err == nil {
		powerStateStr := powerStateToString(powerState)
		ch <- prometheus.MustNewConstMetric(e.modemPowerState, prometheus.GaugeValue, 1.0, deviceID, powerStateStr)
	}

	// Signal quality
	if quality, _, err := modem.GetSignalQuality(); err == nil {
		ch <- prometheus.MustNewConstMetric(e.modemSignalQuality, prometheus.GaugeValue, float64(quality), deviceID)
	}

	// Access technology
	if accessTechs, err := modem.GetAccessTechnologies(); err == nil {
		// Use the first technology in the list
		if len(accessTechs) > 0 {
			techStr := accessTechToString(accessTechs[0])
			ch <- prometheus.MustNewConstMetric(e.modemAccessTech, prometheus.GaugeValue, 1.0, deviceID, techStr)
		}
	}

	// Unlock required
	if unlockRequired, err := modem.GetUnlockRequired(); err == nil {
		ch <- prometheus.MustNewConstMetric(e.modemUnlockRequired, prometheus.GaugeValue, float64(unlockRequired), deviceID)
	}
}

func (e *Exporter) collectSignalMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	signal, err := modem.GetSignal()
	if err != nil {
		// Signal interface might not be available
		return
	}

	// LTE signal
	if lte, err := signal.GetLte(); err == nil && lte.Rssi != 0 {
		ch <- prometheus.MustNewConstMetric(e.signalLteRssi, prometheus.GaugeValue, lte.Rssi, deviceID)
		if lte.Rsrq != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalLteRsrq, prometheus.GaugeValue, lte.Rsrq, deviceID)
		}
		if lte.Rsrp != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalLteRsrp, prometheus.GaugeValue, lte.Rsrp, deviceID)
		}
		if lte.Snr != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalLteSnr, prometheus.GaugeValue, lte.Snr, deviceID)
		}
	}

	// UMTS signal
	if umts, err := signal.GetUmts(); err == nil && umts.Rssi != 0 {
		ch <- prometheus.MustNewConstMetric(e.signalUmtsRssi, prometheus.GaugeValue, umts.Rssi, deviceID)
		if umts.Ecio != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalUmtsEcio, prometheus.GaugeValue, umts.Ecio, deviceID)
		}
		if umts.Rscp != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalUmtsRscp, prometheus.GaugeValue, umts.Rscp, deviceID)
		}
	}

	// GSM signal
	if gsm, err := signal.GetGsm(); err == nil && gsm.Rssi != 0 {
		ch <- prometheus.MustNewConstMetric(e.signalGsmRssi, prometheus.GaugeValue, gsm.Rssi, deviceID)
	}

	// CDMA signal
	if cdma, err := signal.GetCdma(); err == nil && cdma.Rssi != 0 {
		ch <- prometheus.MustNewConstMetric(e.signalCdmaRssi, prometheus.GaugeValue, cdma.Rssi, deviceID)
		if cdma.Ecio != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalCdmaEcio, prometheus.GaugeValue, cdma.Ecio, deviceID)
		}
	}

	// EVDO signal
	if evdo, err := signal.GetEvdo(); err == nil && evdo.Rssi != 0 {
		ch <- prometheus.MustNewConstMetric(e.signalEvdoRssi, prometheus.GaugeValue, evdo.Rssi, deviceID)
		if evdo.Ecio != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalEvdoEcio, prometheus.GaugeValue, evdo.Ecio, deviceID)
		}
		if evdo.Sinr != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalEvdoSinr, prometheus.GaugeValue, evdo.Sinr, deviceID)
		}
		if evdo.Io != 0 {
			ch <- prometheus.MustNewConstMetric(e.signalEvdoIo, prometheus.GaugeValue, evdo.Io, deviceID)
		}
	}
}

func (e *Exporter) collectBearerMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	bearers, err := modem.GetBearers()
	if err != nil {
		return
	}

	for _, bearer := range bearers {
		// Bearer info
		iface, _ := bearer.GetInterface()
		connected, _ := bearer.GetConnected()
		bearerPath := bearer.GetObjectPath()

		ipConfig, err := bearer.GetIp4Config()
		ipMethod := ""
		ipAddress := ""
		if err == nil {
			ipMethod = fmt.Sprint(ipConfig.Method)
			ipAddress = ipConfig.Address
		}

		ch <- prometheus.MustNewConstMetric(
			e.bearerInfo,
			prometheus.GaugeValue,
			1.0,
			deviceID, string(bearerPath), iface, ipMethod, ipAddress,
		)

		// Bearer connected status
		connectedValue := 0.0
		if connected {
			connectedValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(e.bearerConnected, prometheus.GaugeValue, connectedValue, deviceID, string(bearerPath))
	}
}

func (e *Exporter) collectSIMMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	sim, err := modem.GetSim()
	if err != nil {
		return
	}

	simPath := sim.GetObjectPath()
	imsi, _ := sim.GetImsi()
	operatorName, _ := sim.GetOperatorName()

	ch <- prometheus.MustNewConstMetric(
		e.simInfo,
		prometheus.GaugeValue,
		1.0,
		deviceID, string(simPath), imsi, operatorName,
	)
}

func (e *Exporter) collect3GPPMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	modem3gpp, err := modem.Get3gpp()
	if err != nil {
		return
	}

	// Registration state
	if regState, err := modem3gpp.GetRegistrationState(); err == nil {
		regStateStr := registrationStateToString(regState)
		ch <- prometheus.MustNewConstMetric(e.modem3gppRegistrationState, prometheus.GaugeValue, 1.0, deviceID, regStateStr)
	}

	// Operator code
	if operatorCode, err := modem3gpp.GetOperatorCode(); err == nil && operatorCode != "" {
		ch <- prometheus.MustNewConstMetric(e.modem3gppOperatorCode, prometheus.GaugeValue, 1.0, deviceID, operatorCode)
	}

	// Operator name
	if operatorName, err := modem3gpp.GetOperatorName(); err == nil && operatorName != "" {
		ch <- prometheus.MustNewConstMetric(e.modem3gppOperatorName, prometheus.GaugeValue, 1.0, deviceID, operatorName)
	}
}

func (e *Exporter) collectMessagingMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	messaging, err := modem.GetMessaging()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.messagingSupported, prometheus.GaugeValue, 0.0, deviceID)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.messagingSupported, prometheus.GaugeValue, 1.0, deviceID)

	// Get SMS count
	if messages, err := messaging.GetMessages(); err == nil {
		ch <- prometheus.MustNewConstMetric(e.smsCount, prometheus.GaugeValue, float64(len(messages)), deviceID)
	}
}

func (e *Exporter) collectLocationMetrics(ch chan<- prometheus.Metric, modem modemmanager.Modem, deviceID string) {
	location, err := modem.GetLocation()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.locationEnabled, prometheus.GaugeValue, 0.0, deviceID)
		return
	}

	// Check if location is enabled
	if signalsLocation, err := location.GetSignalsLocation(); err == nil {
		enabledValue := 0.0
		if signalsLocation {
			enabledValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(e.locationEnabled, prometheus.GaugeValue, enabledValue, deviceID)

		// Get location data if enabled
		if signalsLocation {
			if loc, err := location.GetLocation(); err == nil {
				// Export GPS location if available
				if loc.GpsRaw.Latitude != 0 || loc.GpsRaw.Longitude != 0 {
					ch <- prometheus.MustNewConstMetric(e.locationLatitude, prometheus.GaugeValue, loc.GpsRaw.Latitude, deviceID)
					ch <- prometheus.MustNewConstMetric(e.locationLongitude, prometheus.GaugeValue, loc.GpsRaw.Longitude, deviceID)
					if loc.GpsRaw.Altitude != 0 {
						ch <- prometheus.MustNewConstMetric(e.locationAltitude, prometheus.GaugeValue, loc.GpsRaw.Altitude, deviceID)
					}
				}
			}
		}
	}
}

// Helper functions to convert enums to strings
func stateToString(state modemmanager.MMModemState) string {
	switch state {
	case modemmanager.MmModemStateFailed:
		return "failed"
	case modemmanager.MmModemStateUnknown:
		return "unknown"
	case modemmanager.MmModemStateInitializing:
		return "initializing"
	case modemmanager.MmModemStateLocked:
		return "locked"
	case modemmanager.MmModemStateDisabled:
		return "disabled"
	case modemmanager.MmModemStateDisabling:
		return "disabling"
	case modemmanager.MmModemStateEnabling:
		return "enabling"
	case modemmanager.MmModemStateEnabled:
		return "enabled"
	case modemmanager.MmModemStateSearching:
		return "searching"
	case modemmanager.MmModemStateRegistered:
		return "registered"
	case modemmanager.MmModemStateDisconnecting:
		return "disconnecting"
	case modemmanager.MmModemStateConnecting:
		return "connecting"
	case modemmanager.MmModemStateConnected:
		return "connected"
	default:
		return "unknown"
	}
}

func powerStateToString(state modemmanager.MMModemPowerState) string {
	switch state {
	case modemmanager.MmModemPowerStateUnknown:
		return "unknown"
	case modemmanager.MmModemPowerStateOff:
		return "off"
	case modemmanager.MmModemPowerStateLow:
		return "low"
	case modemmanager.MmModemPowerStateOn:
		return "on"
	default:
		return "unknown"
	}
}

func accessTechToString(tech modemmanager.MMModemAccessTechnology) string {
	// This is a simplified version - you might want to handle multiple technologies
	switch {
	case tech&modemmanager.MmModemAccessTechnologyLte != 0:
		return "lte"
	case tech&modemmanager.MmModemAccessTechnologyHspaPlus != 0:
		return "hspa_plus"
	case tech&modemmanager.MmModemAccessTechnologyHspa != 0:
		return "hspa"
	case tech&modemmanager.MmModemAccessTechnologyUmts != 0:
		return "umts"
	case tech&modemmanager.MmModemAccessTechnologyEdge != 0:
		return "edge"
	case tech&modemmanager.MmModemAccessTechnologyGprs != 0:
		return "gprs"
	case tech&modemmanager.MmModemAccessTechnologyGsm != 0:
		return "gsm"
	default:
		return "unknown"
	}
}

func registrationStateToString(state modemmanager.MMModem3gppRegistrationState) string {
	switch state {
	case modemmanager.MmModem3gppRegistrationStateIdle:
		return "idle"
	case modemmanager.MmModem3gppRegistrationStateHome:
		return "home"
	case modemmanager.MmModem3gppRegistrationStateSearching:
		return "searching"
	case modemmanager.MmModem3gppRegistrationStateDenied:
		return "denied"
	case modemmanager.MmModem3gppRegistrationStateUnknown:
		return "unknown"
	case modemmanager.MmModem3gppRegistrationStateRoaming:
		return "roaming"
	default:
		return "unknown"
	}
}
