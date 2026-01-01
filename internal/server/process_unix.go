//go:build !windows

package server

import "syscall"

// setupProcessGroup sets up the process group for Unix systems
func (p *Process) setupProcessGroup() {
	p.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

