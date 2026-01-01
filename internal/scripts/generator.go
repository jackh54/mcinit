package scripts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackh54/mcinit/internal/utils"
)

// Generator generates startup scripts
type Generator struct {
	serverDir string
}

// NewGenerator creates a new Generator instance
func NewGenerator(serverDir string) *Generator {
	return &Generator{
		serverDir: serverDir,
	}
}

// GenerateAll generates all startup scripts (Unix + Windows)
func (g *Generator) GenerateAll(javaPath, jarPath, xms, xmx string, jvmFlags []string) error {
	if err := g.GenerateUnix(javaPath, jarPath, xms, xmx, jvmFlags); err != nil {
		return fmt.Errorf("failed to generate Unix script: %w", err)
	}

	if err := g.GenerateWindows(javaPath, jarPath, xms, xmx, jvmFlags); err != nil {
		return fmt.Errorf("failed to generate Windows scripts: %w", err)
	}

	return nil
}

// GenerateUnix generates Unix startup script (start.sh)
func (g *Generator) GenerateUnix(javaPath, jarPath, xms, xmx string, jvmFlags []string) error {
	script := generateUnixScript(javaPath, jarPath, xms, xmx, jvmFlags)

	scriptPath := filepath.Join(g.serverDir, "start.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to write start.sh: %w", err)
	}

	return nil
}

// GenerateWindows generates Windows startup scripts (start.ps1 and start.cmd)
func (g *Generator) GenerateWindows(javaPath, jarPath, xms, xmx string, jvmFlags []string) error {
	// Generate PowerShell script
	ps1Script := generateWindowsPowerShellScript(javaPath, jarPath, xms, xmx, jvmFlags)
	ps1Path := filepath.Join(g.serverDir, "start.ps1")
	if err := os.WriteFile(ps1Path, []byte(ps1Script), 0644); err != nil {
		return fmt.Errorf("failed to write start.ps1: %w", err)
	}

	// Generate CMD script
	cmdScript := generateWindowsCMDScript(javaPath, jarPath, xms, xmx, jvmFlags)
	cmdPath := filepath.Join(g.serverDir, "start.cmd")
	if err := os.WriteFile(cmdPath, []byte(cmdScript), 0644); err != nil {
		return fmt.Errorf("failed to write start.cmd: %w", err)
	}

	return nil
}

// EnsureExecutable ensures the Unix script is executable
func (g *Generator) EnsureExecutable(scriptName string) error {
	if utils.IsWindows() {
		return nil // Not needed on Windows
	}

	scriptPath := filepath.Join(g.serverDir, scriptName)
	return os.Chmod(scriptPath, 0755)
}
