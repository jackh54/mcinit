package cli

import (
	"fmt"

	"github.com/jackh54/mcinit/internal/server"
	"github.com/spf13/cobra"
)

var (
	force bool
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running server gracefully",
	Long:  `Stop the Minecraft server gracefully using RCON or STDIN, with optional force kill.`,
	Example: `  mcinit stop
  mcinit stop --force`,
	RunE: runStop,
}

func init() {
	stopCmd.Flags().BoolVar(&force, "force", false, "Force kill if graceful shutdown fails")
}

func runStop(cmd *cobra.Command, args []string) error {
	// Get server directory
	serverDir := "."

	// Create server manager
	mgr, err := server.NewManager(serverDir)
	if err != nil {
		return fmt.Errorf("failed to create server manager: %w", err)
	}

	// Check if running
	if !mgr.IsRunning() {
		return fmt.Errorf("server is not running")
	}

	// Stop server
	if force {
		printf("Force stopping server...\n")
	} else {
		printf("Stopping server gracefully...\n")
	}

	if err := mgr.Stop(force); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	printf("Server stopped\n")
	return nil
}
