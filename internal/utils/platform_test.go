package utils

import (
	"runtime"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	platform := GetPlatform()
	
	// Should return one of the known platforms
	validPlatforms := []Platform{PlatformWindows, PlatformMacOS, PlatformLinux}
	valid := false
	for _, p := range validPlatforms {
		if platform == p {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("GetPlatform() returned unexpected platform: %s", platform)
	}

	// Check it matches runtime.GOOS
	switch runtime.GOOS {
	case "windows":
		if platform != PlatformWindows {
			t.Error("Expected PlatformWindows on Windows")
		}
	case "darwin":
		if platform != PlatformMacOS {
			t.Error("Expected PlatformMacOS on macOS")
		}
	case "linux":
		if platform != PlatformLinux {
			t.Error("Expected PlatformLinux on Linux")
		}
	}
}

func TestIsWindows(t *testing.T) {
	isWin := IsWindows()
	expected := runtime.GOOS == "windows"
	
	if isWin != expected {
		t.Errorf("IsWindows() = %v, want %v", isWin, expected)
	}
}

func TestIsUnix(t *testing.T) {
	isUnix := IsUnix()
	expected := runtime.GOOS == "darwin" || runtime.GOOS == "linux"
	
	if isUnix != expected {
		t.Errorf("IsUnix() = %v, want %v", isUnix, expected)
	}
}

