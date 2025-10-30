package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/maltegrosse/go-modemmanager/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	version = "1.0.0"
)

var (
	listenAddress = flag.String("listen-address", ":9539", "Address on which to expose metrics and web interface")
	metricsPath   = flag.String("metrics-path", "/metrics", "Path under which to expose metrics")
	signalRate    = flag.Duration("signal-rate", 5*time.Second, "How frequently ModemManager should poll each modem for extended signal strength data (0 to disable)")
	showVersion   = flag.Bool("version", false, "Show version information and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("mm-exporter version %s\n", version)
		os.Exit(0)
	}

	log.Printf("Starting ModemManager Exporter v%s", version)
	log.Printf("Listening on %s", *listenAddress)
	log.Printf("Metrics path: %s", *metricsPath)
	log.Printf("Signal refresh rate: %s", *signalRate)

	// Connect to ModemManager
	mm, err := modemmanager.NewModemManager()
	if err != nil {
		log.Fatalf("Failed to connect to ModemManager: %v", err)
	}
	log.Println("Successfully connected to ModemManager")

	// Get ModemManager version
	mmVersion, err := mm.GetVersion()
	if err != nil {
		log.Printf("Warning: Failed to get ModemManager version: %v", err)
	} else {
		log.Printf("ModemManager version: %s", mmVersion)
	}

	// Setup signal monitoring for each modem
	if *signalRate > 0 {
		if err := setupSignalMonitoring(mm, *signalRate); err != nil {
			log.Printf("Warning: Failed to setup signal monitoring: %v", err)
		}
	}

	// Create Prometheus registry
	registry := prometheus.NewRegistry()

	// Register standard collectors
	registry.MustRegister(
		collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Register ModemManager exporter
	mmExporter := exporter.NewExporter(mm)
	registry.MustRegister(mmExporter)

	log.Println("Registered all collectors")

	// Setup HTTP handlers
	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, "", log.LstdFlags),
		ErrorHandling: promhttp.ContinueOnError,
	}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>ModemManager Exporter</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 40px; }
		h1 { color: #333; }
		.info { background: #f0f0f0; padding: 15px; border-radius: 5px; }
		.links { margin-top: 20px; }
		a { color: #0066cc; text-decoration: none; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<h1>ModemManager Exporter</h1>
	<div class="info">
		<p><strong>Version:</strong> %s</p>
		<p><strong>ModemManager Version:</strong> %s</p>
		<p><strong>Signal Refresh Rate:</strong> %s</p>
	</div>
	<div class="links">
		<p><a href="%s">Metrics</a></p>
	</div>
</body>
</html>
`, version, mmVersion, *signalRate, *metricsPath)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK\n")
	})

	// Setup graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:         *listenAddress,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-quit
		log.Println("Shutting down server...")

		if err := server.Close(); err != nil {
			log.Printf("Error closing server: %v", err)
		}
		close(done)
	}()

	log.Printf("Server is ready to handle requests at %s", *listenAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	<-done
	log.Println("Server stopped")
}

func setupSignalMonitoring(mm modemmanager.ModemManager, rate time.Duration) error {
	modems, err := mm.GetModems()
	if err != nil {
		return fmt.Errorf("failed to get modems: %w", err)
	}

	if len(modems) == 0 {
		log.Println("No modems found")
		return nil
	}

	log.Printf("Setting up signal monitoring for %d modem(s)", len(modems))

	for _, modem := range modems {
		deviceID, err := modem.GetDeviceIdentifier()
		if err != nil {
			log.Printf("Warning: Failed to get device identifier: %v", err)
			continue
		}

		model, err := modem.GetModel()
		if err != nil {
			model = "unknown"
		}

		log.Printf("Configuring modem %s (%s)", deviceID, model)

		// Get signal interface
		signal, err := modem.GetSignal()
		if err != nil {
			log.Printf("Warning: Signal interface not available for modem %s: %v", deviceID, err)
			continue
		}

		// Setup signal refresh rate
		rateSeconds := uint32(rate.Seconds())
		if err := signal.Setup(rateSeconds); err != nil {
			log.Printf("Warning: Failed to setup signal monitoring for modem %s: %v", deviceID, err)
			continue
		}

		log.Printf("Signal monitoring enabled for modem %s (refresh rate: %s)", deviceID, rate)
	}

	return nil
}
