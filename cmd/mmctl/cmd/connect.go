package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/spf13/cobra"
)

var (
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to mobile network",
		Long: `Create a data connection to the mobile network.

This command creates a bearer connection and activates it. You can specify
connection parameters like APN, username, and password.`,
		Example: `  # Simple connect with APN
  mmctl connect -m 0 --apn internet

  # Connect with authentication
  mmctl connect -m 0 --apn internet --user myuser --password mypass

  # Connect with specific IP type
  mmctl connect -m 0 --apn internet --ip-type ipv4v6`,
		RunE: runConnect,
	}

	disconnectCmd = &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from mobile network",
		Long:  `Disconnect an active data connection.`,
		Example: `  # Disconnect modem 0
  mmctl disconnect -m 0`,
		RunE: runDisconnect,
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get connection status",
		Long:  `Display current connection status including IP configuration.`,
		Example: `  # Get connection status
  mmctl status -m 0

  # Get status in JSON format
  mmctl status -m 0 --json`,
		RunE: runStatus,
	}

	// Connect flags
	apn          string
	username     string
	password     string
	ipType       string
	allowRoaming bool
)

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(disconnectCmd)
	rootCmd.AddCommand(statusCmd)

	// Connect command flags
	connectCmd.Flags().StringVarP(&apn, "apn", "a", "", "Access Point Name (required)")
	connectCmd.MarkFlagRequired("apn")
	connectCmd.Flags().StringVarP(&username, "user", "u", "", "Username for authentication")
	connectCmd.Flags().StringVarP(&password, "password", "P", "", "Password for authentication")
	connectCmd.Flags().StringVar(&ipType, "ip-type", "ipv4", "IP type (ipv4, ipv6, ipv4v6)")
	connectCmd.Flags().BoolVar(&allowRoaming, "allow-roaming", false, "Allow connection while roaming")
}

func runConnect(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get the simple interface for easy connection
	simple, err := modem.GetSimpleModem()
	if err != nil {
		return fmt.Errorf("failed to get simple modem interface: %w", err)
	}

	if verbose {
		fmt.Printf("Connecting to network with APN: %s\n", apn)
		fmt.Printf("IP Type: %s\n", ipType)
		if username != "" {
			fmt.Printf("Username: %s\n", username)
		}
		if allowRoaming {
			fmt.Println("Roaming: allowed")
		}
	}

	// Parse IP type
	var ipFamily modemmanager.MMBearerIpFamily
	switch ipType {
	case "ipv4":
		ipFamily = modemmanager.MmBearerIpFamilyIpv4
	case "ipv6":
		ipFamily = modemmanager.MmBearerIpFamilyIpv6
	case "ipv4v6":
		ipFamily = modemmanager.MmBearerIpFamilyIpv4v6
	default:
		return fmt.Errorf("invalid IP type: %s (must be ipv4, ipv6, or ipv4v6)", ipType)
	}

	// Create connection properties
	props := modemmanager.SimpleProperties{
		Apn:            apn,
		User:           username,
		Password:       password,
		IpType:         ipFamily,
		AllowedRoaming: allowRoaming,
	}

	// Connect
	fmt.Println("Connecting...")
	bearer, err := simple.Connect(props)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Wait for connection to establish
	if verbose {
		fmt.Println("Waiting for connection to establish...")
	}
	time.Sleep(2 * time.Second)

	// Get connection status
	connected, err := bearer.GetConnected()
	if err != nil {
		return fmt.Errorf("failed to get connection status: %w", err)
	}

	if !connected {
		return fmt.Errorf("connection failed - bearer not connected")
	}

	fmt.Println("✓ Connected successfully!")

	// Get IP configuration
	if verbose {
		fmt.Println("\nConnection details:")

		if iface, err := bearer.GetInterface(); err == nil {
			fmt.Printf("Interface: %s\n", iface)
		}

		if ipv4Config, err := bearer.GetIp4Config(); err == nil {
			fmt.Printf("\nIPv4 Configuration:\n")
			fmt.Printf("  Address:  %s/%d\n", ipv4Config.Address, ipv4Config.Prefix)
			fmt.Printf("  Gateway:  %s\n", ipv4Config.Gateway)
			dns := []string{}
			if ipv4Config.Dns1 != "" {
				dns = append(dns, ipv4Config.Dns1)
			}
			if ipv4Config.Dns2 != "" {
				dns = append(dns, ipv4Config.Dns2)
			}
			if ipv4Config.Dns3 != "" {
				dns = append(dns, ipv4Config.Dns3)
			}
			if len(dns) > 0 {
				fmt.Printf("  DNS:      %v\n", dns)
			}
		}

		if ipv6Config, err := bearer.GetIp6Config(); err == nil && ipv6Config.Address != "" {
			fmt.Printf("\nIPv6 Configuration:\n")
			fmt.Printf("  Address:  %s/%d\n", ipv6Config.Address, ipv6Config.Prefix)
			fmt.Printf("  Gateway:  %s\n", ipv6Config.Gateway)
			dns := []string{}
			if ipv6Config.Dns1 != "" {
				dns = append(dns, ipv6Config.Dns1)
			}
			if ipv6Config.Dns2 != "" {
				dns = append(dns, ipv6Config.Dns2)
			}
			if ipv6Config.Dns3 != "" {
				dns = append(dns, ipv6Config.Dns3)
			}
			if len(dns) > 0 {
				fmt.Printf("  DNS:      %v\n", dns)
			}
		}
	}

	return nil
}

func runDisconnect(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get the simple interface
	simple, err := modem.GetSimpleModem()
	if err != nil {
		return fmt.Errorf("failed to get simple modem interface: %w", err)
	}

	if verbose {
		fmt.Printf("Disconnecting modem %d...\n", modemIndex)
	}

	// Get bearers to disconnect
	bearers, err := modem.GetBearers()
	if err != nil {
		return fmt.Errorf("failed to get bearers: %w", err)
	}

	if len(bearers) == 0 {
		return fmt.Errorf("no active bearers found")
	}

	// Disconnect each bearer
	for _, bearer := range bearers {
		connected, err := bearer.GetConnected()
		if err != nil {
			continue
		}

		if connected {
			if err := simple.Disconnect(bearer); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to disconnect bearer: %v\n", err)
			} else {
				fmt.Println("✓ Disconnected successfully")
			}
		}
	}

	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get modem state
	state, err := modem.GetState()
	if err != nil {
		return fmt.Errorf("failed to get modem state: %w", err)
	}

	// Get bearers
	bearers, err := modem.GetBearers()
	if err != nil {
		return fmt.Errorf("failed to get bearers: %w", err)
	}

	// Build status information
	status := make(map[string]interface{})
	status["state"] = state.String()
	status["connected"] = state == modemmanager.MmModemStateConnected

	// Get signal quality
	if signalPercent, _, err := modem.GetSignalQuality(); err == nil {
		status["signal_quality"] = signalPercent
	}

	// Get 3GPP info
	if modem3gpp, err := modem.Get3gpp(); err == nil {
		if regState, err := modem3gpp.GetRegistrationState(); err == nil {
			status["registration_state"] = regState.String()
		}
		if opName, err := modem3gpp.GetOperatorName(); err == nil {
			status["operator"] = opName
		}
	}

	// Get bearer information
	if len(bearers) > 0 {
		bearerInfos := make([]map[string]interface{}, 0)
		for _, bearer := range bearers {
			info := make(map[string]interface{})

			connected, _ := bearer.GetConnected()
			info["connected"] = connected

			if iface, err := bearer.GetInterface(); err == nil {
				info["interface"] = iface
			}

			if props, err := bearer.GetProperties(); err == nil {
				info["apn"] = props.APN
				info["ip_type"] = props.IPType.String()
			}

			if connected {
				if ipv4, err := bearer.GetIp4Config(); err == nil {
					dns := []string{}
					if ipv4.Dns1 != "" {
						dns = append(dns, ipv4.Dns1)
					}
					if ipv4.Dns2 != "" {
						dns = append(dns, ipv4.Dns2)
					}
					if ipv4.Dns3 != "" {
						dns = append(dns, ipv4.Dns3)
					}
					info["ipv4"] = map[string]interface{}{
						"address": ipv4.Address,
						"prefix":  ipv4.Prefix,
						"gateway": ipv4.Gateway,
						"dns":     dns,
					}
				}

				if ipv6, err := bearer.GetIp6Config(); err == nil && ipv6.Address != "" {
					dns := []string{}
					if ipv6.Dns1 != "" {
						dns = append(dns, ipv6.Dns1)
					}
					if ipv6.Dns2 != "" {
						dns = append(dns, ipv6.Dns2)
					}
					if ipv6.Dns3 != "" {
						dns = append(dns, ipv6.Dns3)
					}
					info["ipv6"] = map[string]interface{}{
						"address": ipv6.Address,
						"prefix":  ipv6.Prefix,
						"gateway": ipv6.Gateway,
						"dns":     dns,
					}
				}

				if stats, err := bearer.GetStats(); err == nil {
					info["stats"] = map[string]interface{}{
						"bytes_rx": stats.RxBytes,
						"bytes_tx": stats.TxBytes,
						"duration": fmt.Sprintf("%ds", stats.Duration),
					}
				}
			}

			bearerInfos = append(bearerInfos, info)
		}
		status["bearers"] = bearerInfos
	}

	// Output
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(status)
	}

	// Table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Connection Status\n")
	fmt.Fprintf(w, "=================\n\n")
	fmt.Fprintf(w, "State:\t%s\n", formatState(status["state"].(string)))

	if signal, ok := status["signal_quality"].(uint32); ok {
		fmt.Fprintf(w, "Signal:\t%d%%\n", signal)
	}

	if regState, ok := status["registration_state"].(string); ok {
		fmt.Fprintf(w, "Registration:\t%s\n", regState)
	}

	if operator, ok := status["operator"].(string); ok {
		fmt.Fprintf(w, "Operator:\t%s\n", operator)
	}

	// Bearer information
	if bearers, ok := status["bearers"].([]map[string]interface{}); ok && len(bearers) > 0 {
		fmt.Fprintf(w, "\nData Connection:\n")
		for i, bearer := range bearers {
			fmt.Fprintf(w, "\nBearer %d:\n", i)

			if connected, ok := bearer["connected"].(bool); ok {
				if connected {
					fmt.Fprintf(w, "  Status:\t✓ Connected\n")
				} else {
					fmt.Fprintf(w, "  Status:\t✗ Disconnected\n")
				}
			}

			if iface, ok := bearer["interface"].(string); ok {
				fmt.Fprintf(w, "  Interface:\t%s\n", iface)
			}

			if apn, ok := bearer["apn"].(string); ok {
				fmt.Fprintf(w, "  APN:\t%s\n", apn)
			}

			if ipv4, ok := bearer["ipv4"].(map[string]interface{}); ok {
				fmt.Fprintf(w, "  IPv4:\t%s/%v\n", ipv4["address"], ipv4["prefix"])
				fmt.Fprintf(w, "  Gateway:\t%s\n", ipv4["gateway"])
				if dns, ok := ipv4["dns"].([]string); ok && len(dns) > 0 {
					fmt.Fprintf(w, "  DNS:\t%v\n", dns)
				}
			}

			if stats, ok := bearer["stats"].(map[string]interface{}); ok {
				if rxBytes, ok := stats["bytes_rx"].(uint64); ok {
					fmt.Fprintf(w, "  RX:\t%d bytes\n", rxBytes)
				}
				if txBytes, ok := stats["bytes_tx"].(uint64); ok {
					fmt.Fprintf(w, "  TX:\t%d bytes\n", txBytes)
				}
				if duration, ok := stats["duration"].(string); ok {
					fmt.Fprintf(w, "  Duration:\t%s\n", duration)
				}
			}
		}
	} else {
		fmt.Fprintf(w, "\nData Connection:\tNot connected\n")
	}

	return nil
}
