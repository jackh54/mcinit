package cli

import (
	"fmt"

	"github.com/jackh54/mcinit/internal/server"
	"github.com/spf13/cobra"
)

var (
	background bool
	extraArgs  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server process",
	Long:  `Start the Minecraft server in the current directory.`,
	Example: `  mcinit start
  mcinit start --background
  mcinit start --args "-XX:+UseG1GC"`,
	RunE: runStart,
}

func init() {
	startCmd.Flags().BoolVar(&background, "background", false, "Run server in background")
	startCmd.Flags().StringVar(&extraArgs, "args", "", "Additional JVM arguments")
}

func runStart(cmd *cobra.Command, args []string) error {
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

	// Check if already running
	if mgr.IsRunning() {
		return fmt.Errorf("server is already running")
	}

	// Start server
	printf("Starting server...\n")
	if err := mgr.Start(background, extraArgs); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	if background {
		printf("Server started in background\n")
	} else {
		printf("Server stopped\n")
	}

	return nil
}
