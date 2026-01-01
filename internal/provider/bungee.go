package provider

import (
	"context"
	"fmt"
)

// BungeeProvider implements Provider for BungeeCord
// Note: BungeeCord doesn't have a proper API, this is a placeholder
type BungeeProvider struct{}

// NewBungeeProvider creates a new BungeeProvider
func NewBungeeProvider() *BungeeProvider {
	return &BungeeProvider{}
}

// GetName returns the provider name
func (p *BungeeProvider) GetName() string {
	return "bungee"
}

// GetAvailableVersions returns error (not supported for BungeeCord)
func (p *BungeeProvider) GetAvailableVersions(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("BungeeCord does not have version API support - manual download required")
}

// GetLatestBuild returns error (not supported)
func (p *BungeeProvider) GetLatestBuild(ctx context.Context, version string) (string, error) {
	return "", fmt.Errorf("BungeeCord does not have API support - manual download required")
}

// DownloadJar returns error (not supported)
func (p *BungeeProvider) DownloadJar(ctx context.Context, version, build string) (string, string, string, error) {
	return "", "", "", fmt.Errorf("BungeeCord does not have API support - manual download required")
}

// GetDownloadURL returns error (not supported)
func (p *BungeeProvider) GetDownloadURL(ctx context.Context, version, build string) (string, error) {
	return "", fmt.Errorf("BungeeCord does not have API support - manual download required")
}

// GetChecksum returns error (not supported)
func (p *BungeeProvider) GetChecksum(ctx context.Context, version, build string) (string, string, error) {
	return "", "", fmt.Errorf("BungeeCord does not have API support - manual download required")
}

