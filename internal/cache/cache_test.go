package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if cache.GetBaseDir() == "" {
		t.Error("Cache base directory is empty")
	}

	// Check that base dir exists or can be created
	if err := cache.EnsureJarsDir(); err != nil {
		t.Errorf("EnsureJarsDir() error = %v", err)
	}
}

func TestGetJarPath(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	path := cache.GetJarPath("paper", "1.21.4", "123")
	
	if path == "" {
		t.Error("GetJarPath() returned empty string")
	}

	if !filepath.IsAbs(path) {
		t.Error("GetJarPath() should return absolute path")
	}

	// Check that path contains expected filename pattern
	filename := filepath.Base(path)
	if filename == "" {
		t.Error("GetJarPath() returned path with empty filename")
	}

	// Filename should include server type and version in some form
	// The actual format may vary, so we just check it's not empty
	if len(filename) < 5 {
		t.Errorf("GetJarPath() filename seems too short: %s", filename)
	}
}

func TestGetMetadataPath(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	metaPath := cache.GetMetadataPath("paper", "1.21.4", "123")
	jarPath := cache.GetJarPath("paper", "1.21.4", "123")
	
	if metaPath == "" {
		t.Error("GetMetadataPath() returned empty string")
	}

	// Should be jar path + .meta.json
	expected := jarPath + ".meta.json"
	if metaPath != expected {
		t.Errorf("GetMetadataPath() = %s, want %s", metaPath, expected)
	}
}

func TestHasJar(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Non-existent jar should return false
	if cache.HasJar("nonexistent", "1.0.0", "1") {
		t.Error("HasJar() returned true for non-existent jar")
	}
}

func TestSaveAndGetMetadata(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create metadata
	meta := &CacheMetadata{
		Version:     "1.21.4",
		Build:       "123",
		ServerType:  "paper",
		DownloadURL: "https://example.com/paper.jar",
		Checksum:    "abc123",
		Algorithm:   "sha256",
		CachedAt:    "2024-01-01T00:00:00Z",
	}

	// Save metadata
	err = cache.SaveMetadata("test-server", "1.0.0", "1", meta)
	if err != nil {
		t.Errorf("SaveMetadata() error = %v", err)
	}

	// Get metadata
	loaded, err := cache.GetMetadata("test-server", "1.0.0", "1")
	if err != nil {
		t.Errorf("GetMetadata() error = %v", err)
	}

	// Verify loaded metadata
	if loaded.Version != meta.Version {
		t.Errorf("Loaded version = %s, want %s", loaded.Version, meta.Version)
	}

	if loaded.ServerType != meta.ServerType {
		t.Errorf("Loaded server type = %s, want %s", loaded.ServerType, meta.ServerType)
	}

	// Clean up
	metaPath := cache.GetMetadataPath("test-server", "1.0.0", "1")
	_ = os.Remove(metaPath)
}

func TestEnsureJarsDir(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := cache.EnsureJarsDir(); err != nil {
		t.Errorf("EnsureJarsDir() error = %v", err)
	}

	jarsDir := filepath.Join(cache.GetBaseDir(), "jars")
	if _, err := os.Stat(jarsDir); os.IsNotExist(err) {
		t.Error("Jars directory was not created")
	}
}

