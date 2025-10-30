package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available modems",
	Long: `List all modems detected by ModemManager.

This command displays information about all available modems including:
  - Modem index
  - Manufacturer
  - Model
  - State
  - Signal quality
  - Equipment identifier (IMEI)

Use --json flag for machine-readable output.`,
	Example: `  # List all modems
  mmctl list

  # List modems in JSON format
  mmctl list --json

  # List modems with verbose output
  mmctl list --verbose`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

type modemInfo struct {
	Index               int    `json:"index"`
	Path                string `json:"path"`
	Manufacturer        string `json:"manufacturer"`
	Model               string `json:"model"`
	State               string `json:"state"`
	SignalQuality       uint32 `json:"signal_quality"`
	EquipmentIdentifier string `json:"equipment_identifier"`
	Device              string `json:"device"`
	PrimaryPort         string `json:"primary_port"`
}

func runList(cmd *cobra.Command, args []string) error {
	// Connect to ModemManager
	mm, err := modemmanager.NewModemManager()
	if err != nil {
		return fmt.Errorf("failed to connect to ModemManager: %w", err)
	}

	if verbose {
		version, err := mm.GetVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not get ModemManager version: %v\n", err)
		} else {
			fmt.Printf("ModemManager version: %s\n\n", version)
		}
	}

	// Get list of modems
	modems, err := mm.GetModems()
	if err != nil {
		return fmt.Errorf("failed to get modems: %w", err)
	}

	if len(modems) == 0 {
		fmt.Println("No modems found")
		return nil
	}

	// Collect modem information
	var modemInfos []modemInfo
	for i, modem := range modems {
		info := modemInfo{
			Index: i,
			Path:  string(modem.GetObjectPath()),
		}

		// Get manufacturer
		if manufacturer, err := modem.GetManufacturer(); err == nil {
			info.Manufacturer = manufacturer
		}

		// Get model
		if model, err := modem.GetModel(); err == nil {
			info.Model = model
		}

		// Get state
		if state, err := modem.GetState(); err == nil {
			info.State = state.String()
		}

		// Get signal quality
		if signalPercent, _, err := modem.GetSignalQuality(); err == nil {
			info.SignalQuality = signalPercent
		}

		// Get equipment identifier (IMEI)
		if imei, err := modem.GetEquipmentIdentifier(); err == nil {
			info.EquipmentIdentifier = imei
		}

		// Get device identifier
		if device, err := modem.GetDeviceIdentifier(); err == nil {
			info.Device = device
		}

		// Get primary port - Not available in current API
		// Using device identifier as fallback
		info.PrimaryPort = ""

		modemInfos = append(modemInfos, info)
	}

	// Output results
	if jsonOutput {
		return outputJSON(modemInfos)
	}

	return outputTable(modemInfos)
}

func outputJSON(modems []modemInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(modems); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

func outputTable(modems []modemInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Header
	fmt.Fprintln(w, "INDEX\tMANUFACTURER\tMODEL\tSTATE\tSIGNAL\tIMEI\tPORT")
	fmt.Fprintln(w, "-----\t------------\t-----\t-----\t------\t----\t----")

	// Rows
	for _, modem := range modems {
		signal := fmt.Sprintf("%d%%", modem.SignalQuality)
		if modem.SignalQuality == 0 {
			signal = "N/A"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			modem.Index,
			truncate(modem.Manufacturer, 20),
			truncate(modem.Model, 20),
			formatState(modem.State),
			signal,
			truncate(modem.EquipmentIdentifier, 15),
			truncate(modem.PrimaryPort, 15),
		)
	}

	if verbose {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Total modems: %d\n", len(modems))
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatState(state string) string {
	// Remove "MmModemState" prefix for cleaner output
	if len(state) > 12 && state[:12] == "MmModemState" {
		return state[12:]
	}
	return state
}
