package cli

import (
	"github.com/spf13/cobra"
)

func newMailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mail",
		Short: "Agent Mail commands",
		Long:  "Manage Agent Mail messages, inbox, and file reservations.",
	}

	cmd.AddCommand(newMailInboxCmd())
	
	// Future commands: ack, read, lock, unlock (if moved here)

	return cmd
}
