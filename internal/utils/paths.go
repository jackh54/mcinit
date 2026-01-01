package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// ExpandPath expands ~ and environment variables in a path
func ExpandPath(path string) (string, error) {
	// Expand home directory
	expanded, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}

	// Expand environment variables
	expanded = os.ExpandEnv(expanded)

	return expanded, nil
}

// NormalizePath converts a path to the native format
func NormalizePath(path string) string {
	return filepath.FromSlash(path)
}

// JoinPath joins path elements using the OS-specific separator
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// AbsolutePath returns the absolute path
func AbsolutePath(path string) (string, error) {
	expanded, err := ExpandPath(path)
	if err != nil {
		return "", err
	}

	return filepath.Abs(expanded)
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// PathExists checks if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory checks if a path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// RelativePath returns a relative path from base to target
func RelativePath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

// CleanPath cleans and normalizes a path
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// SplitPath splits a path into directory and file
func SplitPath(path string) (dir, file string) {
	return filepath.Split(path)
}

// GetExecutableName returns the executable name for the current platform
func GetExecutableName(name string) string {
	if IsWindows() {
		if !strings.HasSuffix(name, ".exe") {
			return name + ".exe"
		}
	}
	return name
}

