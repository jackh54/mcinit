package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackh54/mcinit/internal/utils"
)

// State represents the server runtime state
type State struct {
	PID       int       `json:"pid"`
	StartTime time.Time `json:"startTime"`
	Status    string    `json:"status"` // "running", "stopped", "crashed"
}

// StateFile manages the server state file
type StateFile struct {
	serverDir string
	statePath string
	pidPath   string
}

// NewStateFile creates a new StateFile instance
func NewStateFile(serverDir string) *StateFile {
	mcinitDir := filepath.Join(serverDir, ".mcinit")
	return &StateFile{
		serverDir: serverDir,
		statePath: filepath.Join(mcinitDir, "state.json"),
		pidPath:   filepath.Join(mcinitDir, "pid"),
	}
}

// Load loads the state from disk
func (s *StateFile) Load() (*State, error) {
	if !utils.PathExists(s.statePath) {
		return nil, fmt.Errorf("state file not found")
	}

	data, err := os.ReadFile(s.statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// Save saves the state to disk
func (s *StateFile) Save(state *State) error {
	// Ensure .mcinit directory exists
	mcinitDir := filepath.Dir(s.statePath)
	if err := utils.EnsureDir(mcinitDir); err != nil {
		return fmt.Errorf("failed to create .mcinit directory: %w", err)
	}

	// Marshal state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write state file
	if err := os.WriteFile(s.statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	// Write PID file
	pidData := fmt.Sprintf("%d\n", state.PID)
	if err := os.WriteFile(s.pidPath, []byte(pidData), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	return nil
}

// ReadPID reads the PID from the PID file
func (s *StateFile) ReadPID() (int, error) {
	if !utils.PathExists(s.pidPath) {
		return 0, fmt.Errorf("PID file not found")
	}

	data, err := os.ReadFile(s.pidPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read PID file: %w", err)
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		return 0, fmt.Errorf("failed to parse PID: %w", err)
	}

	return pid, nil
}

// Clear removes the state and PID files
func (s *StateFile) Clear() error {
	var lastErr error

	if utils.PathExists(s.statePath) {
		if err := os.Remove(s.statePath); err != nil {
			lastErr = err
		}
	}

	if utils.PathExists(s.pidPath) {
		if err := os.Remove(s.pidPath); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// Exists checks if state files exist
func (s *StateFile) Exists() bool {
	return utils.PathExists(s.statePath) || utils.PathExists(s.pidPath)
}

// GetUptime calculates the server uptime
func (s *State) GetUptime() time.Duration {
	if s.StartTime.IsZero() {
		return 0
	}
	return time.Since(s.StartTime)
}

// FormatUptime returns a human-readable uptime string
func (s *State) FormatUptime() string {
	uptime := s.GetUptime()
	if uptime == 0 {
		return "0s"
	}

	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

