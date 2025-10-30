package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	jsonOutput bool
	verbose    bool
	modemIndex int
	modemPath  string
	version    = "0.1.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mmctl",
	Short: "A CLI tool for managing modems via ModemManager",
	Long: `mmctl is a command-line interface for interacting with ModemManager.

It provides a user-friendly way to manage cellular modems, including:
  - Listing and managing modems
  - Enabling/disabling modems
  - Creating and managing data connections
  - Sending and receiving SMS messages
  - Getting signal quality and network information
  - Managing SIM cards
  - And much more...

This tool uses the go-modemmanager library to communicate with ModemManager
via D-Bus.`,
	Version: version,
	Example: `  # List all modems
  mmctl list

  # Get modem information
  mmctl modem info -i 0

  # Enable a modem
  mmctl modem enable -i 0

  # Connect to network
  mmctl connect -i 0 --apn internet

  # Send SMS
  mmctl sms send -i 0 --number +1234567890 --text "Hello World"

  # Get signal quality
  mmctl modem signal -i 0`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().IntVarP(&modemIndex, "modem", "m", -1, "Modem index (alternative to --path)")
	rootCmd.PersistentFlags().StringVarP(&modemPath, "path", "p", "", "Modem D-Bus path")

	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Helper function to print version info
func printVersion() {
	fmt.Printf("mmctl version %s\n", version)
	fmt.Println("Built with go-modemmanager library")
}
