package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"simple path", "/tmp/test", false},
		{"relative path", "./test", false},
		{"home dir", "~/test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expanded, err := ExpandPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && expanded == "" {
				t.Error("ExpandPath() returned empty string")
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	path := "some/path/to/file"
	normalized := NormalizePath(path)
	
	if normalized == "" {
		t.Error("NormalizePath() returned empty string")
	}

	// On Windows, should contain backslashes
	if runtime.GOOS == "windows" {
		if !filepath.IsAbs(normalized) && normalized == path {
			// Should be different on Windows
		}
	}
}

func TestJoinPath(t *testing.T) {
	result := JoinPath("dir1", "dir2", "file.txt")
	
	if result == "" {
		t.Error("JoinPath() returned empty string")
	}

	expected := filepath.Join("dir1", "dir2", "file.txt")
	if result != expected {
		t.Errorf("JoinPath() = %s, want %s", result, expected)
	}
}

func TestEnsureDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mcinit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "test", "nested", "dir")
	
	if err := EnsureDir(testDir); err != nil {
		t.Errorf("EnsureDir() error = %v", err)
	}

	if !PathExists(testDir) {
		t.Error("Directory was not created")
	}

	if !IsDirectory(testDir) {
		t.Error("Path exists but is not a directory")
	}
}

func TestPathExists(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	if !PathExists(tempFile.Name()) {
		t.Error("PathExists() returned false for existing file")
	}

	if PathExists("/nonexistent/path/file.txt") {
		t.Error("PathExists() returned true for non-existent path")
	}
}

func TestIsDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mcinit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if !IsDirectory(tempDir) {
		t.Error("IsDirectory() returned false for directory")
	}

	tempFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)

	if IsDirectory(tempFile) {
		t.Error("IsDirectory() returned true for file")
	}
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "/tmp/test"},
		{"with dots", "/tmp/../test"},
		{"with double slash", "/tmp//test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned := CleanPath(tt.input)
			if cleaned == "" {
				t.Error("CleanPath() returned empty string")
			}
		})
	}
}

func TestGetExecutableName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"unix binary", "mcinit", "mcinit"},
	}

	if runtime.GOOS == "windows" {
		tests = append(tests, struct {
			name     string
			input    string
			expected string
		}{"windows exe", "mcinit", "mcinit.exe"})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetExecutableName(tt.input)
			if runtime.GOOS == "windows" {
				if result != "mcinit.exe" {
					t.Errorf("GetExecutableName() on Windows = %s, want mcinit.exe", result)
				}
			} else {
				if result != "mcinit" {
					t.Errorf("GetExecutableName() on Unix = %s, want mcinit", result)
				}
			}
		})
	}
}

