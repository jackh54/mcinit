package server

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

// Process manages a server process
type Process struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	pid       int
	serverDir string
	stateFile *StateFile
}

// NewProcess creates a new Process instance
func NewProcess(serverDir string) *Process {
	return &Process{
		serverDir: serverDir,
		stateFile: NewStateFile(serverDir),
	}
}

// Start starts the server process
func (p *Process) Start(javaPath, jarPath string, jvmArgs []string, background bool) error {
	// Build command arguments
	args := append(jvmArgs, "-jar", jarPath, "nogui")

	// Create command
	p.cmd = exec.Command(javaPath, args...)
	p.cmd.Dir = p.serverDir

	// Set up stdin
	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	p.stdin = stdin

	// Set up stdout
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	p.stdout = stdout

	// Set up stderr
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	p.stderr = stderr

	// Set up process group for proper cleanup
	p.setupProcessGroup()

	// Start the process
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	p.pid = p.cmd.Process.Pid

	// Save state
	state := &State{
		PID:       p.pid,
		StartTime: time.Now(),
		Status:    "running",
	}
	if err := p.stateFile.Save(state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	if background {
		// In background mode, start goroutines to handle output
		go p.pipeOutput(p.stdout, os.Stdout)
		go p.pipeOutput(p.stderr, os.Stderr)
		
		// Start goroutine to wait for process and update state
		go func() {
			p.cmd.Wait()
			p.stateFile.Clear()
		}()
	} else {
		// In foreground mode, pipe output synchronously
		go p.pipeOutput(p.stdout, os.Stdout)
		go p.pipeOutput(p.stderr, os.Stderr)
		
		// Wait for process to complete
		if err := p.cmd.Wait(); err != nil {
			p.stateFile.Clear()
			return fmt.Errorf("server process exited with error: %w", err)
		}
		p.stateFile.Clear()
	}

	return nil
}

// Stop stops the server process gracefully
func (p *Process) Stop(force bool) error {
	// Read PID from state
	pid, err := p.stateFile.ReadPID()
	if err != nil {
		return fmt.Errorf("server not running or PID file not found")
	}

	// Find process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if force {
		// Force kill
		return p.kill(process)
	}

	// Try graceful shutdown first
	if err := p.gracefulStop(process, pid); err != nil {
		fmt.Printf("Graceful shutdown failed, forcing...\n")
		return p.kill(process)
	}

	return nil
}

// gracefulStop attempts a graceful shutdown via STDIN "stop" command
func (p *Process) gracefulStop(process *os.Process, pid int) error {
	// Try to send "stop" command via stdin if we have access
	// This is a simplified version - in production, you'd need to handle this more carefully
	
	// Send SIGTERM (Unix) or appropriate signal (Windows)
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM: %w", err)
	}

	// Wait for process to exit (with timeout)
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("graceful shutdown timed out")
		case <-ticker.C:
			// Check if process still exists
			if err := process.Signal(syscall.Signal(0)); err != nil {
				// Process is gone
				p.stateFile.Clear()
				return nil
			}
		}
	}
}

// kill forcefully kills the process
func (p *Process) kill(process *os.Process) error {
	var err error
	
	if runtime.GOOS == "windows" {
		err = process.Kill()
	} else {
		err = process.Signal(syscall.SIGKILL)
	}

	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	// Wait a bit for process to die
	time.Sleep(time.Second)
	p.stateFile.Clear()
	
	return nil
}

// IsRunning checks if the server is currently running
func (p *Process) IsRunning() bool {
	pid, err := p.stateFile.ReadPID()
	if err != nil {
		return false
	}

	// Try to find the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Check if process is alive by sending signal 0 (doesn't actually send a signal)
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// GetPID returns the process PID
func (p *Process) GetPID() (int, error) {
	return p.stateFile.ReadPID()
}

// GetState returns the current state
func (p *Process) GetState() (*State, error) {
	return p.stateFile.Load()
}

// pipeOutput pipes output from one reader to a writer
func (p *Process) pipeOutput(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Fprintln(writer, scanner.Text())
	}
}

// setupProcessGroup sets up the process group for proper cleanup
func (p *Process) setupProcessGroup() {
	if runtime.GOOS == "windows" {
		// Windows: Set process group flag
		p.cmd.SysProcAttr = &syscall.SysProcAttr{
			// Windows-specific flags would go here
			// For now, we keep it simple
		}
	} else {
		// Unix: Create new process group
		p.cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
	}
}

