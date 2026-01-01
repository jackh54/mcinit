package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jackh54/mcinit/internal/utils"
)

// Cache manages the mcinit cache directory
type Cache struct {
	baseDir string
}

// CacheMetadata stores metadata about cached jars
type CacheMetadata struct {
	Version     string `json:"version"`
	Build       string `json:"build,omitempty"`
	ServerType  string `json:"serverType"`
	DownloadURL string `json:"downloadUrl"`
	Checksum    string `json:"checksum"`
	Algorithm   string `json:"algorithm"` // "sha1" or "sha256"
	CachedAt    string `json:"cachedAt"`
}

// New creates a new Cache instance
func New() (*Cache, error) {
	baseDir, err := getCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	cache := &Cache{
		baseDir: baseDir,
	}

	// Ensure cache directory exists
	if err := utils.EnsureDir(baseDir); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cache, nil
}

// getCacheDir returns the platform-specific cache directory
func getCacheDir() (string, error) {
	var cacheDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %LOCALAPPDATA%\mcinit\cache
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		cacheDir = filepath.Join(localAppData, "mcinit", "cache")

	case "darwin":
		// macOS: ~/Library/Caches/mcinit
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		cacheDir = filepath.Join(home, "Library", "Caches", "mcinit")

	case "linux":
		// Linux: ~/.cache/mcinit or $XDG_CACHE_HOME/mcinit
		xdgCache := os.Getenv("XDG_CACHE_HOME")
		if xdgCache != "" {
			cacheDir = filepath.Join(xdgCache, "mcinit")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			cacheDir = filepath.Join(home, ".cache", "mcinit")
		}

	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cacheDir, nil
}

// GetBaseDir returns the cache base directory
func (c *Cache) GetBaseDir() string {
	return c.baseDir
}

// GetJarPath returns the path where a jar file should be cached
func (c *Cache) GetJarPath(serverType, version, build string) string {
	filename := fmt.Sprintf("%s-%s", serverType, version)
	if build != "" {
		filename += fmt.Sprintf("-%s", build)
	}
	filename += ".jar"

	return filepath.Join(c.baseDir, "jars", filename)
}

// GetMetadataPath returns the path where jar metadata is stored
func (c *Cache) GetMetadataPath(serverType, version, build string) string {
	jarPath := c.GetJarPath(serverType, version, build)
	return jarPath + ".meta.json"
}

// HasJar checks if a jar is already cached
func (c *Cache) HasJar(serverType, version, build string) bool {
	jarPath := c.GetJarPath(serverType, version, build)
	metaPath := c.GetMetadataPath(serverType, version, build)

	return utils.PathExists(jarPath) && utils.PathExists(metaPath)
}

// GetMetadata reads the metadata for a cached jar
func (c *Cache) GetMetadata(serverType, version, build string) (*CacheMetadata, error) {
	metaPath := c.GetMetadataPath(serverType, version, build)

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var meta CacheMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &meta, nil
}

// SaveMetadata saves metadata for a cached jar
func (c *Cache) SaveMetadata(serverType, version, build string, meta *CacheMetadata) error {
	metaPath := c.GetMetadataPath(serverType, version, build)

	// Ensure parent directory exists
	if err := utils.EnsureDir(filepath.Dir(metaPath)); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// VerifyChecksum verifies the checksum of a cached jar
func (c *Cache) VerifyChecksum(jarPath, expectedChecksum, algorithm string) (bool, error) {
	file, err := os.Open(jarPath)
	if err != nil {
		return false, fmt.Errorf("failed to open jar file: %w", err)
	}
	defer file.Close()

	// For now, only implement SHA256 verification
	// SHA1 can be added later if needed
	if algorithm != "sha256" && algorithm != "sha1" {
		return false, fmt.Errorf("unsupported checksum algorithm: %s", algorithm)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false, fmt.Errorf("failed to compute checksum: %w", err)
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	return actualChecksum == expectedChecksum, nil
}

// EnsureJarsDir ensures the jars directory exists
func (c *Cache) EnsureJarsDir() error {
	jarsDir := filepath.Join(c.baseDir, "jars")
	return utils.EnsureDir(jarsDir)
}

// Clear removes all cached files
func (c *Cache) Clear() error {
	return os.RemoveAll(c.baseDir)
}

// GetSize returns the total size of the cache in bytes
func (c *Cache) GetSize() (int64, error) {
	var size int64

	err := filepath.Walk(c.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
