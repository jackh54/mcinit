package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	dryRun   bool
	verbose  bool
	noColor  bool
)

var rootCmd = &cobra.Command{
	Use:   "mcinit",
	Short: "Cross-platform Minecraft developer server management tool",
	Long: `mcinit is a CLI tool for Minecraft plugin developers that creates and manages
local dev servers quickly and reproducibly across Windows, macOS, and Linux.`,
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./mcinit.json)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview changes without executing")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(statusCmd)
}

// getConfigPath returns the config file path
func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return "mcinit.json"
}

// printf prints formatted output if not in dry-run mode
func printf(format string, args ...interface{}) {
	if verbose || !dryRun {
		fmt.Printf(format, args...)
	}
}

// verboseLog prints verbose output
func verboseLog(format string, args ...interface{}) {
	if verbose {
		fmt.Printf("[DEBUG] "+format, args...)
	}
}

// errorLog prints error messages
func errorLog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format, args...)
}

