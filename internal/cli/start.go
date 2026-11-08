package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	background bool
	extraArgs  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server process",
	Long:  `Start the Minecraft server in the current directory, or the one holding --config.`,
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
	// Resolve the server from --config (or ./mcinit.json)
	mgr, err := newManager()
	if err != nil {
		return err
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
