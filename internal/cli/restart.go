package cli

import (
	"fmt"
	"time"

	"github.com/jackh54/mcinit/internal/server"
	"github.com/spf13/cobra"
)

var (
	waitSeconds int
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the server (stop + start)",
	Long:  `Restart the Minecraft server by stopping and starting it.`,
	Example: `  mcinit restart
  mcinit restart --background --wait 5`,
	RunE: runRestart,
}

func init() {
	restartCmd.Flags().BoolVar(&background, "background", false, "Run in background after restart")
	restartCmd.Flags().IntVar(&waitSeconds, "wait", 2, "Wait time before restart (seconds)")
}

func runRestart(cmd *cobra.Command, args []string) error {
	// Get server directory
	serverDir := "."

	// Create server manager
	mgr, err := server.NewManager(serverDir)
	if err != nil {
		return fmt.Errorf("failed to create server manager: %w", err)
	}

	// Load configuration
	if err := mgr.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Restart server
	printf("Restarting server...\n")

	// Stop if running
	if mgr.IsRunning() {
		printf("Stopping server...\n")
		if err := mgr.Stop(false); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	// Wait
	if waitSeconds > 0 {
		printf("Waiting %d seconds before restart...\n", waitSeconds)
		time.Sleep(time.Duration(waitSeconds) * time.Second)
	}

	// Start
	printf("Starting server...\n")
	if err := mgr.Start(background, ""); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	if background {
		printf("Server restarted in background\n")
	} else {
		printf("Server stopped\n")
	}

	return nil
}
