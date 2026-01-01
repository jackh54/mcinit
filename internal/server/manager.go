package server

import (
	"fmt"
	"path/filepath"

	"github.com/jackh54/mcinit/internal/config"
	"github.com/jackh54/mcinit/internal/java"
	"github.com/jackh54/mcinit/internal/utils"
	"github.com/jackh54/mcinit/pkg/jvmflags"
)

// Manager manages server lifecycle
type Manager struct {
	serverDir string
	config    *config.Config
	process   *Process
}

// NewManager creates a new Manager instance
func NewManager(serverDir string) (*Manager, error) {
	absPath, err := utils.AbsolutePath(serverDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server directory: %w", err)
	}

	return &Manager{
		serverDir: absPath,
		process:   NewProcess(absPath),
	}, nil
}

// LoadConfig loads the server configuration
func (m *Manager) LoadConfig() error {
	cfgPath := filepath.Join(m.serverDir, "mcinit.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	m.config = cfg
	return nil
}

// Start starts the server
func (m *Manager) Start(background bool, extraArgs string) error {
	// Load config if not already loaded
	if m.config == nil {
		if err := m.LoadConfig(); err != nil {
			return err
		}
	}

	// Check if already running
	if m.process.IsRunning() {
		return fmt.Errorf("server is already running")
	}

	// Resolve Java path
	javaPath, err := m.resolveJavaPath()
	if err != nil {
		return fmt.Errorf("failed to resolve Java path: %w", err)
	}

	// Build JVM arguments
	jvmArgs, err := m.buildJVMArgs(extraArgs)
	if err != nil {
		return fmt.Errorf("failed to build JVM arguments: %w", err)
	}

	// Get jar path
	jarPath := filepath.Join(m.serverDir, m.config.Server.JarPath)
	if !utils.PathExists(jarPath) {
		return fmt.Errorf("server jar not found: %s", jarPath)
	}

	// Start the process
	return m.process.Start(javaPath, jarPath, jvmArgs, background)
}

// Stop stops the server
func (m *Manager) Stop(force bool) error {
	return m.process.Stop(force)
}

// Restart restarts the server
func (m *Manager) Restart(background bool, waitSeconds int) error {
	// Stop the server
	if m.process.IsRunning() {
		if err := m.Stop(false); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	// Wait before restarting
	if waitSeconds > 0 {
		fmt.Printf("Waiting %d seconds before restart...\n", waitSeconds)
		// In a real implementation, you'd use time.Sleep here
	}

	// Start the server
	return m.Start(background, "")
}

// IsRunning checks if the server is running
func (m *Manager) IsRunning() bool {
	return m.process.IsRunning()
}

// GetStatus returns the server status
func (m *Manager) GetStatus() (*State, error) {
	return m.process.GetState()
}

// resolveJavaPath resolves the Java executable path
func (m *Manager) resolveJavaPath() (string, error) {
	if m.config.Java.Path != "" && m.config.Java.Path != "auto" {
		// Use explicitly configured path
		return m.config.Java.Path, nil
	}

	// Auto-detect Java
	detector := java.NewDetector()
	
	if m.config.Java.Version != "" && m.config.Java.Version != "auto" {
		// Find specific version
		var majorVersion int
		fmt.Sscanf(m.config.Java.Version, "%d", &majorVersion)
		
		inst, err := detector.FindByMajorVersion(majorVersion)
		if err != nil {
			return "", fmt.Errorf("Java %d not found: %w", majorVersion, err)
		}
		return inst.Path, nil
	}

	// Find best available Java
	inst, err := detector.FindBest()
	if err != nil {
		return "", fmt.Errorf("no Java installation found: %w", err)
	}

	return inst.Path, nil
}

// buildJVMArgs builds the JVM arguments array
func (m *Manager) buildJVMArgs(extraArgs string) ([]string, error) {
	// Get base flags from preset
	flags, err := jvmflags.Get(
		m.config.JVM.Flags,
		m.config.JVM.Xms,
		m.config.JVM.Xmx,
		m.config.JVM.CustomFlags,
	)
	if err != nil {
		return nil, err
	}

	// Parse and add extra arguments
	if extraArgs != "" {
		extra := jvmflags.ParseCustomFlags(extraArgs)
		flags = append(flags, extra...)
	}

	return flags, nil
}

