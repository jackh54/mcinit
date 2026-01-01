package provider

import (
	"context"
)

// Provider defines the interface for server jar providers
type Provider interface {
	// GetAvailableVersions returns all available Minecraft versions
	GetAvailableVersions(ctx context.Context) ([]string, error)

	// GetLatestBuild returns the latest build number for a given version
	// Returns empty string for vanilla (no builds)
	GetLatestBuild(ctx context.Context, version string) (string, error)

	// DownloadJar downloads the server jar for the given version/build
	// Returns local file path, download URL, and checksum
	DownloadJar(ctx context.Context, version, build string) (localPath, downloadURL, checksum string, err error)

	// GetDownloadURL returns the direct download URL without downloading
	GetDownloadURL(ctx context.Context, version, build string) (string, error)

	// GetChecksum returns the expected checksum for verification
	GetChecksum(ctx context.Context, version, build string) (algorithm, checksum string, err error)

	// GetName returns the provider name
	GetName() string
}

// BuildInfo represents build information for a server
type BuildInfo struct {
	Version     string
	Build       string
	DownloadURL string
	Checksum    string
	Algorithm   string // "sha1" or "sha256"
}

