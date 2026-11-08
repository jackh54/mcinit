package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Show server status (running/stopped, PID, uptime)",
	Long:    `Display the current status of the Minecraft server.`,
	Example: `  mcinit status`,
	RunE:    runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Resolve the server from --config (or ./mcinit.json)
	mgr, err := newManager()
	if err != nil {
		return err
	}

	// Check status
	if !mgr.IsRunning() {
		printf("Server status: STOPPED\n")
		return nil
	}

	// Get detailed status
	state, err := mgr.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get server status: %w", err)
	}

	printf("Server status: RUNNING\n")
	printf("PID: %d\n", state.PID)
	printf("Uptime: %s\n", state.FormatUptime())
	printf("Started: %s\n", state.StartTime.Format("2006-01-02 15:04:05"))

	return nil
}
