package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	smsCmd = &cobra.Command{
		Use:   "sms",
		Short: "Manage SMS messages",
		Long: `Send, receive, and manage SMS messages.

This command group provides operations for:
  - Sending SMS messages
  - Listing received messages
  - Reading messages
  - Deleting messages`,
		Example: `  # Send SMS
  mmctl sms send -m 0 --number +1234567890 --text "Hello"

  # List messages
  mmctl sms list -m 0

  # Read a message
  mmctl sms read -m 0 --sms-index 0

  # Delete a message
  mmctl sms delete -m 0 --sms-index 0`,
	}

	smsSendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send an SMS message",
		Long: `Send an SMS message to a phone number.

The message will be sent using the modem's messaging interface.`,
		Example: `  # Send simple SMS
  mmctl sms send -m 0 --number +1234567890 --text "Hello World"

  # Send SMS with verbose output
  mmctl sms send -m 0 --number +1234567890 --text "Test" --verbose`,
		RunE: runSmsSend,
	}

	smsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List SMS messages",
		Long: `List all SMS messages stored on the modem.

This includes received, sent, and draft messages.`,
		Example: `  # List all messages
  mmctl sms list -m 0

  # List in JSON format
  mmctl sms list -m 0 --json`,
		RunE: runSmsList,
	}

	smsReadCmd = &cobra.Command{
		Use:   "read",
		Short: "Read an SMS message",
		Long:  `Display the content of a specific SMS message.`,
		Example: `  # Read message at index 0
  mmctl sms read -m 0 --sms-index 0

  # Read message in JSON format
  mmctl sms read -m 0 --sms-index 0 --json`,
		RunE: runSmsRead,
	}

	smsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an SMS message",
		Long:  `Delete a specific SMS message from the modem.`,
		Example: `  # Delete message at index 0
  mmctl sms delete -m 0 --sms-index 0`,
		RunE: runSmsDelete,
	}

	// SMS flags
	smsNumber   string
	smsText     string
	smsIndex    int
	smsValidity int
)

func init() {
	rootCmd.AddCommand(smsCmd)

	// Add subcommands
	smsCmd.AddCommand(smsSendCmd)
	smsCmd.AddCommand(smsListCmd)
	smsCmd.AddCommand(smsReadCmd)
	smsCmd.AddCommand(smsDeleteCmd)

	// Send command flags
	smsSendCmd.Flags().StringVarP(&smsNumber, "number", "n", "", "Recipient phone number (required)")
	smsSendCmd.Flags().StringVarP(&smsText, "text", "t", "", "Message text (required)")
	smsSendCmd.Flags().IntVar(&smsValidity, "validity", 0, "Message validity period in minutes (0 = default)")
	smsSendCmd.MarkFlagRequired("number")
	smsSendCmd.MarkFlagRequired("text")

	// Read and delete command flags
	smsReadCmd.Flags().IntVarP(&smsIndex, "sms-index", "i", 0, "SMS message index")
	smsReadCmd.MarkFlagRequired("sms-index")
	smsDeleteCmd.Flags().IntVarP(&smsIndex, "sms-index", "i", 0, "SMS message index")
	smsDeleteCmd.MarkFlagRequired("sms-index")
}

func runSmsSend(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get messaging interface
	messaging, err := modem.GetMessaging()
	if err != nil {
		return fmt.Errorf("failed to get messaging interface: %w", err)
	}

	if verbose {
		fmt.Printf("Sending SMS to %s\n", smsNumber)
		fmt.Printf("Message: %s\n", smsText)
	}

	// Create SMS
	sms, err := messaging.Create(smsNumber, smsText)
	if err != nil {
		return fmt.Errorf("failed to create SMS: %w", err)
	}

	if verbose {
		fmt.Println("SMS created, sending...")
	}

	// Send SMS
	if err := sms.Send(); err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	fmt.Println("✓ SMS sent successfully")

	if verbose {
		// Get SMS state
		if state, err := sms.GetState(); err == nil {
			fmt.Printf("Final state: %s\n", state.String())
		}
	}

	return nil
}

func runSmsList(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get messaging interface
	messaging, err := modem.GetMessaging()
	if err != nil {
		return fmt.Errorf("failed to get messaging interface: %w", err)
	}

	if verbose {
		fmt.Println("Retrieving SMS messages...")
	}

	// List messages
	messages, err := messaging.List()
	if err != nil {
		return fmt.Errorf("failed to list messages: %w", err)
	}

	if len(messages) == 0 {
		fmt.Println("No messages found")
		return nil
	}

	// Collect message information
	type smsInfo struct {
		Index     int       `json:"index"`
		Path      string    `json:"path"`
		Number    string    `json:"number"`
		Text      string    `json:"text"`
		State     string    `json:"state"`
		Timestamp time.Time `json:"timestamp,omitempty"`
		Storage   string    `json:"storage"`
	}

	var smsInfos []smsInfo
	for i, sms := range messages {
		info := smsInfo{
			Index: i,
			Path:  string(sms.GetObjectPath()),
		}

		// Get number
		if number, err := sms.GetNumber(); err == nil {
			info.Number = number
		}

		// Get text
		if text, err := sms.GetText(); err == nil {
			info.Text = text
		}

		// Get state
		if state, err := sms.GetState(); err == nil {
			info.State = state.String()
		}

		// Get timestamp
		if timestamp, err := sms.GetTimestamp(); err == nil {
			info.Timestamp = timestamp
		}

		// Get storage
		if storage, err := sms.GetStorage(); err == nil {
			info.Storage = storage.String()
		}

		smsInfos = append(smsInfos, info)
	}

	// Output
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(smsInfos)
	}

	// Table output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "INDEX\tNUMBER\tSTATE\tTIMESTAMP\tMESSAGE")
	fmt.Fprintln(w, "-----\t------\t-----\t---------\t-------")

	for _, msg := range smsInfos {
		timestamp := ""
		if !msg.Timestamp.IsZero() {
			timestamp = msg.Timestamp.Format("2006-01-02 15:04")
		}

		text := msg.Text
		if len(text) > 50 {
			text = text[:47] + "..."
		}

		state := msg.State
		if len(state) > 10 && state[:10] == "MmSmsState" {
			state = state[10:]
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			msg.Index,
			truncate(msg.Number, 15),
			state,
			timestamp,
			text,
		)
	}

	if verbose {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Total messages: %d\n", len(smsInfos))
	}

	return nil
}

func runSmsRead(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get messaging interface
	messaging, err := modem.GetMessaging()
	if err != nil {
		return fmt.Errorf("failed to get messaging interface: %w", err)
	}

	// List messages
	messages, err := messaging.List()
	if err != nil {
		return fmt.Errorf("failed to list messages: %w", err)
	}

	if smsIndex >= len(messages) {
		return fmt.Errorf("SMS index %d out of range (0-%d)", smsIndex, len(messages)-1)
	}

	sms := messages[smsIndex]

	// Collect SMS details
	info := make(map[string]interface{})
	info["index"] = smsIndex
	info["path"] = string(sms.GetObjectPath())

	if number, err := sms.GetNumber(); err == nil {
		info["number"] = number
	}

	if text, err := sms.GetText(); err == nil {
		info["text"] = text
	}

	if state, err := sms.GetState(); err == nil {
		info["state"] = state.String()
	}

	if pduType, err := sms.GetPduType(); err == nil {
		info["pdu_type"] = pduType.String()
	}

	if timestamp, err := sms.GetTimestamp(); err == nil {
		info["timestamp"] = timestamp.Format(time.RFC3339)
	}

	if dischargeTimestamp, err := sms.GetDischargeTimestamp(); err == nil {
		info["discharge_timestamp"] = dischargeTimestamp.Format(time.RFC3339)
	}

	if storage, err := sms.GetStorage(); err == nil {
		info["storage"] = storage.String()
	}

	if smsc, err := sms.GetSmsc(); err == nil {
		info["smsc"] = smsc
	}

	if deliveryState, err := sms.GetDeliveryState(); err == nil {
		info["delivery_state"] = deliveryState.String()
	}

	// Output
	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(info)
	}

	// Formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "SMS Message Details")
	fmt.Fprintln(w, "===================")
	fmt.Fprintln(w)

	if number, ok := info["number"].(string); ok {
		fmt.Fprintf(w, "From/To:\t%s\n", number)
	}

	if state, ok := info["state"].(string); ok {
		fmt.Fprintf(w, "State:\t%s\n", state)
	}

	if timestamp, ok := info["timestamp"].(string); ok {
		fmt.Fprintf(w, "Timestamp:\t%s\n", timestamp)
	}

	if storage, ok := info["storage"].(string); ok {
		fmt.Fprintf(w, "Storage:\t%s\n", storage)
	}

	if pduType, ok := info["pdu_type"].(string); ok {
		fmt.Fprintf(w, "Type:\t%s\n", pduType)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Message:")
	fmt.Fprintln(w, "--------")

	if text, ok := info["text"].(string); ok {
		fmt.Fprintln(w, text)
	}

	return nil
}

func runSmsDelete(cmd *cobra.Command, args []string) error {
	modem, err := getModem()
	if err != nil {
		return err
	}

	// Get messaging interface
	messaging, err := modem.GetMessaging()
	if err != nil {
		return fmt.Errorf("failed to get messaging interface: %w", err)
	}

	// List messages
	messages, err := messaging.List()
	if err != nil {
		return fmt.Errorf("failed to list messages: %w", err)
	}

	if smsIndex >= len(messages) {
		return fmt.Errorf("SMS index %d out of range (0-%d)", smsIndex, len(messages)-1)
	}

	sms := messages[smsIndex]

	if verbose {
		if number, err := sms.GetNumber(); err == nil {
			fmt.Printf("Deleting SMS from %s\n", number)
		}
	}

	// Delete the message
	if err := messaging.Delete(sms.GetObjectPath()); err != nil {
		return fmt.Errorf("failed to delete SMS: %w", err)
	}

	fmt.Println("✓ SMS deleted successfully")
	return nil
}
