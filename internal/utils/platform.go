package utils

import (
	"runtime"
)

// Platform represents the current operating system
type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformMacOS   Platform = "darwin"
	PlatformLinux   Platform = "linux"
	PlatformUnknown Platform = "unknown"
)

// GetPlatform returns the current platform
func GetPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return PlatformWindows
	case "darwin":
		return PlatformMacOS
	case "linux":
		return PlatformLinux
	default:
		return PlatformUnknown
	}
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return GetPlatform() == PlatformWindows
}

// IsUnix returns true if running on Unix-like system
func IsUnix() bool {
	p := GetPlatform()
	return p == PlatformMacOS || p == PlatformLinux
}

