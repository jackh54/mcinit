package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackh54/mcinit/internal/logs"
	"github.com/spf13/cobra"
)

var (
	follow bool
	lines  int
	grep   string
	since  string
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Display server logs",
	Long:  `View and filter server logs with various options.`,
	Example: `  mcinit logs
  mcinit logs --follow
  mcinit logs --lines 200 --grep "error"
  mcinit logs --since 30m --grep "warn"`,
	RunE: runLogs,
}

func init() {
	logsCmd.Flags().BoolVar(&follow, "follow", false, "Follow log output (like tail -f)")
	logsCmd.Flags().IntVar(&lines, "lines", 50, "Number of lines to show")
	logsCmd.Flags().StringVar(&grep, "grep", "", "Filter lines matching pattern")
	logsCmd.Flags().StringVar(&since, "since", "", "Show logs since duration (e.g., 10m, 1h)")
}

func runLogs(cmd *cobra.Command, args []string) error {
	// Get server directory
	serverDir := "."

	// Create log reader
	reader := logs.NewReader(serverDir)

	if !reader.Exists() {
		return fmt.Errorf("log file not found - server may not have been started yet")
	}

	// Parse since duration if provided
	var sinceTime time.Time
	if since != "" {
		st, err := logs.GetSinceTime(since)
		if err != nil {
			return fmt.Errorf("invalid since duration: %w", err)
		}
		sinceTime = st
	}

	if follow {
		// Follow mode
		printf("Following logs (press Ctrl+C to stop)...\n")

		output := make(chan string, 100)
		stop := make(chan struct{})
		errChan := make(chan error, 1)

		// Start following
		go func() {
			errChan <- reader.Follow(grep, output, stop)
		}()

		// Handle interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Print output
		for {
			select {
			case line := <-output:
				fmt.Println(line)
			case err := <-errChan:
				if err != nil {
					return fmt.Errorf("error following logs: %w", err)
				}
				return nil
			case <-sigChan:
				close(stop)
				return nil
			}
		}
	} else {
		// Read mode
		logLines, err := reader.Read(lines, grep, sinceTime)
		if err != nil {
			return fmt.Errorf("failed to read logs: %w", err)
		}

		for _, line := range logLines {
			fmt.Println(line)
		}
	}

	return nil
}
