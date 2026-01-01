//go:build windows

package server

import "syscall"

// setupProcessGroup sets up the process group for Windows systems
func (p *Process) setupProcessGroup() {
	p.cmd.SysProcAttr = &syscall.SysProcAttr{
		// Windows uses different process management
		// Process groups are handled differently than Unix
	}
}

