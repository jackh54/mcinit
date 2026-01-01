package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackh54/mcinit/internal/cache"
)

// VanillaProvider implements Provider for Mojang's official server
type VanillaProvider struct {
	cache  *cache.Cache
	client *http.Client
}

// NewVanillaProvider creates a new VanillaProvider
func NewVanillaProvider() *VanillaProvider {
	c, _ := cache.New()
	return &VanillaProvider{
		cache: c,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName returns the provider name
func (p *VanillaProvider) GetName() string {
	return "vanilla"
}

// GetAvailableVersions returns all available Minecraft versions
func (p *VanillaProvider) GetAvailableVersions(ctx context.Context) ([]string, error) {
	manifest, err := p.fetchVersionManifest(ctx)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(manifest.Versions))
	for _, v := range manifest.Versions {
		// Only include release versions by default
		if v.Type == "release" {
			versions = append(versions, v.ID)
		}
	}

	return versions, nil
}

// GetLatestBuild returns empty string (vanilla has no builds)
func (p *VanillaProvider) GetLatestBuild(ctx context.Context, version string) (string, error) {
	return "", nil
}

// DownloadJar downloads the vanilla server jar
func (p *VanillaProvider) DownloadJar(ctx context.Context, version, build string) (string, string, string, error) {
	// Get version details
	manifest, err := p.fetchVersionManifest(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch version manifest: %w", err)
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return "", "", "", fmt.Errorf("version not found: %s", version)
	}

	// Fetch version details
	versionInfo, err := p.fetchVersionInfo(ctx, versionURL)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch version info: %w", err)
	}

	if versionInfo.Downloads.Server.URL == "" {
		return "", "", "", fmt.Errorf("server download not available for version %s", version)
	}

	downloadURL := versionInfo.Downloads.Server.URL
	checksum := versionInfo.Downloads.Server.SHA1

	// Download using cache/downloader
	downloader := cache.NewDownloader(p.cache)
	localPath, err := downloader.DownloadJar(downloadURL, "vanilla", version, "", checksum, "sha1")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to download jar: %w", err)
	}

	return localPath, downloadURL, checksum, nil
}

// GetDownloadURL returns the download URL for a vanilla server
func (p *VanillaProvider) GetDownloadURL(ctx context.Context, version, build string) (string, error) {
	manifest, err := p.fetchVersionManifest(ctx)
	if err != nil {
		return "", err
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return "", fmt.Errorf("version not found: %s", version)
	}

	versionInfo, err := p.fetchVersionInfo(ctx, versionURL)
	if err != nil {
		return "", err
	}

	return versionInfo.Downloads.Server.URL, nil
}

// GetChecksum returns the checksum for a vanilla server
func (p *VanillaProvider) GetChecksum(ctx context.Context, version, build string) (string, string, error) {
	manifest, err := p.fetchVersionManifest(ctx)
	if err != nil {
		return "", "", err
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return "", "", fmt.Errorf("version not found: %s", version)
	}

	versionInfo, err := p.fetchVersionInfo(ctx, versionURL)
	if err != nil {
		return "", "", err
	}

	return "sha1", versionInfo.Downloads.Server.SHA1, nil
}

// fetchVersionManifest fetches the Mojang version manifest
func (p *VanillaProvider) fetchVersionManifest(ctx context.Context) (*VersionManifest, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json", nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest VersionManifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// fetchVersionInfo fetches version-specific information
func (p *VanillaProvider) fetchVersionInfo(ctx context.Context, url string) (*VersionInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info VersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// Mojang API structures

type VersionManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []VersionEntry `json:"versions"`
}

type VersionEntry struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}

type VersionInfo struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Downloads struct {
		Server struct {
			SHA1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}
