package logs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackh54/mcinit/internal/utils"
)

// Reader reads server log files
type Reader struct {
	serverDir string
	logPath   string
}

// NewReader creates a new Reader instance
func NewReader(serverDir string) *Reader {
	// Standard Minecraft server log location
	logPath := filepath.Join(serverDir, "logs", "latest.log")

	return &Reader{
		serverDir: serverDir,
		logPath:   logPath,
	}
}

// Read reads the log file with optional filtering
func (r *Reader) Read(lines int, grepPattern string, sinceTime time.Time) ([]string, error) {
	if !utils.PathExists(r.logPath) {
		return nil, fmt.Errorf("log file not found: %s", r.logPath)
	}

	file, err := os.Open(r.logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var allLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Apply grep filter
		if grepPattern != "" && !r.matchesPattern(line, grepPattern) {
			continue
		}

		// Apply time filter (simplified - actual log timestamp parsing would be more complex)
		if !sinceTime.IsZero() {
			// Skip time filtering for now - would need to parse log timestamps
			// This is a simplified implementation
		}

		allLines = append(allLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	// Return last N lines
	if lines > 0 && lines < len(allLines) {
		return allLines[len(allLines)-lines:], nil
	}

	return allLines, nil
}

// Follow tails the log file in real-time
func (r *Reader) Follow(grepPattern string, output chan<- string, stop <-chan struct{}) error {
	if !utils.PathExists(r.logPath) {
		return fmt.Errorf("log file not found: %s", r.logPath)
	}

	file, err := os.Open(r.logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Seek to end of file
	if _, err := file.Seek(0, 2); err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return nil
		case <-ticker.C:
			for scanner.Scan() {
				line := scanner.Text()

				// Apply grep filter
				if grepPattern != "" && !r.matchesPattern(line, grepPattern) {
					continue
				}

				select {
				case output <- line:
				case <-stop:
					return nil
				}
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading log file: %w", err)
			}
		}
	}
}

// matchesPattern checks if a line matches the grep pattern
func (r *Reader) matchesPattern(line, pattern string) bool {
	// Simple case-insensitive contains match
	return strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
}

// GetLogPath returns the log file path
func (r *Reader) GetLogPath() string {
	return r.logPath
}

// Exists checks if the log file exists
func (r *Reader) Exists() bool {
	return utils.PathExists(r.logPath)
}
