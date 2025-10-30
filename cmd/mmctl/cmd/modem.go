package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/maltegrosse/go-modemmanager"
	"github.com/spf13/cobra"
)

var (
	modemCmd = &cobra.Command{
		Use:   "modem",
		Short: "Manage modem devices",
		Long: `Perform operations on modem devices.

This command group provides subcommands to:
  - Get detailed modem information
  - Enable or disable modems
  - Reset modems
  - Get signal quality
  - Manage modem power state
  - Send AT commands

Use a subcommand to perform a specific operation.`,
		Example: `  # Get modem info
  mmctl modem info -m 0

  # Enable a modem
  mmctl modem enable -m 0

  # Get signal quality
  mmctl modem signal -m 0

  # Send AT command
  mmctl modem command -m 0 "ATI"`,
	}

	modemInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Display detailed modem information",
		Long: `Show comprehensive information about a specific modem including:
  - Hardware details (manufacturer, model, revision)
  - Current state and capabilities
  - Network information
  - SIM information
  - Signal quality
  - Supported technologies`,
		Example: `  # Get info for modem 0
  mmctl modem info -m 0

  # Get info in JSON format
  mmctl modem info -m 0 --json`,
		RunE: runModemInfo,
	}

	modemEnableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enable a modem",
		Long:  `Enable a modem device, powering it on and making it ready for operations.`,
		Example: `  # Enable modem 0
  mmctl modem enable -m 0`,
		RunE: runModemEnable,
	}

	modemDisableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disable a modem",
		Long:  `Disable a modem device, powering it down and stopping operations.`,
		Example: `  # Disable modem 0
  mmctl modem disable -m 0`,
		RunE: runModemDisable,
	}

	modemResetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset a modem",
		Long:  `Reset a modem device to its initial state.`,
		Example: `  # Reset modem 0
  mmctl modem reset -m 0`,
		RunE: runModemReset,
	}

	modemSignalCmd = &cobra.Command{
		Use:   "signal",
		Short: "Get signal quality information",
		Long:  `Display signal quality and strength information for a modem.`,
		Example: `  # Get signal for modem 0
  mmctl modem signal -m 0

  # Get signal in JSON format
  mmctl modem signal -m 0 --json`,
		RunE: runModemSignal,
	}

	modemCommandCmd = &cobra.Command{
		Use:   "command [AT_COMMAND]",
		Short: "Send AT command to modem",
		Long: `Send a raw AT command to the modem and display the response.

Warning: Sending incorrect AT commands can disrupt modem operation.`,
		Example: `  # Get modem information
  mmctl modem command -m 0 "ATI"

  # Get signal quality
  mmctl modem command -m 0 "AT+CSQ"`,
		Args: cobra.ExactArgs(1),
		RunE: runModemCommand,
	}

	// Flags
	commandTimeout uint32
)

func init() {
	rootCmd.AddCommand(modemCmd)

	// Add subcommands
	modemCmd.AddCommand(modemInfoCmd)
	modemCmd.AddCommand(modemEnableCmd)
	modemCmd.AddCommand(modemDisableCmd)
	modemCmd.AddCommand(modemResetCmd)
	modemCmd.AddCommand(modemSignalCmd)
	modemCmd.AddCommand(modemCommandCmd)

	// Command-specific flags
	modemCommandCmd.Flags().Uint32VarP(&commandTimeout, "timeout", "t", 10, "Command timeout in seconds")
}

func getModem() (modemmanager.Modem, error) {
	mm, err := modemmanager.NewModemManager()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ModemManager: %w", err)
	}

	modems, err := mm.GetModems()
	if err != nil {
		return nil, fmt.Errorf("failed to get modems: %w", err)
	}

	if len(modems) == 0 {
		return nil, fmt.Errorf("no modems found")
	}

	if modemIndex < 0 {
		modemIndex = 0
	}

	if modemIndex >= len(modems) {
		return nil, fmt.Errorf("modem index %d out of range (0-%d)", modemIndex, len(modems)-1)
	}

	return modems[modemIndex], nil
}

func runModemInfo(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	info := make(map[string]interface{})

	// Basic information
	if manufacturer, err := modem.GetManufacturer(); err == nil {
		info["manufacturer"] = manufacturer
	}
	if model, err := modem.GetModel(); err == nil {
		info["model"] = model
	}
	if revision, err := modem.GetRevision(); err == nil {
		info["revision"] = revision
	}
	if imei, err := modem.GetEquipmentIdentifier(); err == nil {
		info["equipment_identifier"] = imei
	}
	if deviceId, err := modem.GetDeviceIdentifier(); err == nil {
		info["device_identifier"] = deviceId
	}

	// State information
	if state, err := modem.GetState(); err == nil {
		info["state"] = state.String()
	}
	if powerState, err := modem.GetPowerState(); err == nil {
		info["power_state"] = powerState.String()
	}
	if unlockRequired, err := modem.GetUnlockRequired(); err == nil {
		info["unlock_required"] = unlockRequired.String()
	}

	// Signal quality
	if signal, err := modem.GetSignalQuality(); err == nil {
		info["signal_quality"] = map[string]interface{}{
			"quality": signal.Quality,
			"recent":  signal.Recent,
		}
	}

	// Access technologies
	if accessTech, err := modem.GetAccessTechnologies(); err == nil {
		techs := make([]string, len(accessTech))
		for i, tech := range accessTech {
			techs[i] = tech.String()
		}
		info["access_technologies"] = techs
	}

	// Capabilities
	if caps, err := modem.GetCurrentCapabilities(); err == nil {
		capStrs := make([]string, len(caps))
		for i, cap := range caps {
			capStrs[i] = cap.String()
		}
		info["current_capabilities"] = capStrs
	}

	// Modes
	if modes, err := modem.GetCurrentModes(); err == nil {
		modeStrs := make([]string, len(modes.AllowedModes))
		for i, mode := range modes.AllowedModes {
			modeStrs[i] = mode.String()
		}
		info["current_modes"] = map[string]interface{}{
			"allowed":   modeStrs,
			"preferred": modes.PreferredMode.String(),
		}
	}

	// Bands
	if bands, err := modem.GetCurrentBands(); err == nil {
		bandStrs := make([]string, len(bands))
		for i, band := range bands {
			bandStrs[i] = band.String()
		}
		info["current_bands"] = bandStrs
	}

	// Own numbers
	if numbers, err := modem.GetOwnNumbers(); err == nil {
		info["own_numbers"] = numbers
	}

	// SIM information
	if sim, err := modem.GetSim(); err == nil {
		simInfo := make(map[string]interface{})
		if imsi, err := sim.GetImsi(); err == nil {
			simInfo["imsi"] = imsi
		}
		if iccid, err := sim.GetSimIdentifier(); err == nil {
			simInfo["iccid"] = iccid
		}
		if opId, err := sim.GetOperatorIdentifier(); err == nil {
			simInfo["operator_id"] = opId
		}
		if opName, err := sim.GetOperatorName(); err == nil {
			simInfo["operator_name"] = opName
		}
		info["sim"] = simInfo
	}

	// 3GPP information
	if modem3gpp, err := modem.Get3gpp(); err == nil {
		gppInfo := make(map[string]interface{})
		if imei, err := modem3gpp.GetImei(); err == nil {
			gppInfo["imei"] = imei
		}
		if regState, err := modem3gpp.GetRegistrationState(); err == nil {
			gppInfo["registration_state"] = regState.String()
		}
		if opCode, err := modem3gpp.GetOperatorCode(); err == nil {
			gppInfo["operator_code"] = opCode
		}
		if opName, err := modem3gpp.GetOperatorName(); err == nil {
			gppInfo["operator_name"] = opName
		}
		info["3gpp"] = gppInfo
	}

	// Output
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(info)
	}

	// Table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Property\tValue\n")
	fmt.Fprintf(w, "--------\t-----\n")

	printInfo := func(key string, value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Fprintf(w, "%s:\t\n", key)
			for k, val := range v {
				fmt.Fprintf(w, "  %s\t%v\n", k, val)
			}
		case []string:
			fmt.Fprintf(w, "%s\t%s\n", key, strings.Join(v, ", "))
		default:
			fmt.Fprintf(w, "%s\t%v\n", key, v)
		}
	}

	// Print in order
	keys := []string{
		"manufacturer", "model", "revision", "equipment_identifier",
		"device_identifier", "state", "power_state", "unlock_required",
		"signal_quality", "access_technologies", "current_capabilities",
		"current_modes", "current_bands", "own_numbers", "sim", "3gpp",
	}

	for _, key := range keys {
		if value, ok := info[key]; ok {
			printInfo(key, value)
		}
	}

	return nil
}

func runModemEnable(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Enabling modem %d...\n", modemIndex)
	}

	if err := modem.Enable(true); err != nil {
		return fmt.Errorf("failed to enable modem: %w", err)
	}

	fmt.Println("Modem enabled successfully")
	return nil
}

func runModemDisable(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Disabling modem %d...\n", modemIndex)
	}

	if err := modem.Enable(false); err != nil {
		return fmt.Errorf("failed to disable modem: %w", err)
	}

	fmt.Println("Modem disabled successfully")
	return nil
}

func runModemReset(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Resetting modem %d...\n", modemIndex)
	}

	if err := modem.Reset(); err != nil {
		return fmt.Errorf("failed to reset modem: %w", err)
	}

	fmt.Println("Modem reset successfully")
	return nil
}

func runModemSignal(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	signal, err := modem.GetSignalQuality()
	if err != nil {
		return fmt.Errorf("failed to get signal quality: %w", err)
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(map[string]interface{}{
			"quality": signal.Quality,
			"recent":  signal.Recent,
		})
	}

	fmt.Printf("Signal Quality: %d%%", signal.Quality)
	if signal.Recent {
		fmt.Print(" (recent)")
	}
	fmt.Println()

	// Signal bar representation
	bars := signal.Quality / 20
	fmt.Printf("Signal Bars:    [")
	for i := uint32(0); i < 5; i++ {
		if i < bars {
			fmt.Print("█")
		} else {
			fmt.Print("░")
		}
	}
	fmt.Println("]")

	return nil
}

func runModemCommand(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	atCommand := args[0]

	if verbose {
		fmt.Printf("Sending command: %s\n", atCommand)
	}

	response, err := modem.Command(atCommand, commandTimeout)
	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(map[string]interface{}{
			"command":  atCommand,
			"response": response,
		})
	}

	fmt.Println(response)
	return nil
}
