package cache

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackh54/mcinit/internal/utils"
)

// Downloader handles downloading jar files
type Downloader struct {
	cache  *Cache
	client *http.Client
	// UserAgent is sent with download requests when set.
	UserAgent string
}

// NewDownloader creates a new Downloader instance
func NewDownloader(cache *Cache) *Downloader {
	return &Downloader{
		cache: cache,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// DownloadJar downloads a jar file to the cache
func (d *Downloader) DownloadJar(url, serverType, version, build, expectedChecksum, algorithm string) (string, error) {
	return d.DownloadJarContext(context.Background(), url, serverType, version, build, expectedChecksum, algorithm)
}

// DownloadJarContext is DownloadJar with a context for cancellation.
func (d *Downloader) DownloadJarContext(ctx context.Context, url, serverType, version, build, expectedChecksum, algorithm string) (string, error) {
	// Check if already cached and valid
	if d.cache.HasJar(serverType, version, build) {
		jarPath := d.cache.GetJarPath(serverType, version, build)

		// Verify checksum if provided
		if expectedChecksum != "" {
			valid, err := d.cache.VerifyChecksum(jarPath, expectedChecksum, algorithm)
			if err == nil && valid {
				return jarPath, nil
			}
			// If verification fails, re-download
			fmt.Printf("Cached jar checksum mismatch, re-downloading...\n")
		} else {
			return jarPath, nil
		}
	}

	// Ensure cache directory exists
	if err := d.cache.EnsureJarsDir(); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	jarPath := d.cache.GetJarPath(serverType, version, build)

	// Download to temporary file first
	tempPath := jarPath + ".tmp"
	if err := d.downloadFile(ctx, url, tempPath, expectedChecksum, algorithm); err != nil {
		_ = os.Remove(tempPath)
		return "", fmt.Errorf("failed to download jar: %w", err)
	}

	// Move temp file to final location
	if err := os.Rename(tempPath, jarPath); err != nil {
		_ = os.Remove(tempPath)
		return "", fmt.Errorf("failed to move downloaded jar: %w", err)
	}

	// Save metadata
	meta := &CacheMetadata{
		Version:     version,
		Build:       build,
		ServerType:  serverType,
		DownloadURL: url,
		Checksum:    expectedChecksum,
		Algorithm:   algorithm,
		CachedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := d.cache.SaveMetadata(serverType, version, build, meta); err != nil {
		return jarPath, fmt.Errorf("warning: failed to save metadata: %w", err)
	}

	return jarPath, nil
}

// downloadFile downloads a file from a URL with optional checksum verification
func (d *Downloader) downloadFile(ctx context.Context, url, destPath, expectedChecksum, algorithm string) error {
	// Ensure parent directory exists
	if err := utils.EnsureDir(filepath.Dir(destPath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = out.Close() }()

	// Download file
	fmt.Printf("Downloading from %s...\n", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	if d.UserAgent != "" {
		req.Header.Set("User-Agent", d.UserAgent)
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Set up checksum verification if needed
	var writer io.Writer = out
	var hasher hash.Hash

	if expectedChecksum != "" {
		switch algorithm {
		case "sha256":
			hasher = sha256.New()
		case "sha1":
			hasher = sha1.New()
		default:
			return fmt.Errorf("unsupported checksum algorithm: %s", algorithm)
		}
		writer = io.MultiWriter(out, hasher)
	}

	// Copy with progress indication
	written, err := io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("Downloaded %d bytes\n", written)

	// Verify checksum if provided
	if expectedChecksum != "" && hasher != nil {
		actualChecksum := hex.EncodeToString(hasher.Sum(nil))
		if !strings.EqualFold(actualChecksum, expectedChecksum) {
			return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
		}
		fmt.Printf("Checksum verified (%s)\n", algorithm)
	}

	return nil
}
