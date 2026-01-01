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

// PurpurProvider implements Provider for Purpur
type PurpurProvider struct {
	cache  *cache.Cache
	client *http.Client
}

// NewPurpurProvider creates a new PurpurProvider
func NewPurpurProvider() *PurpurProvider {
	c, _ := cache.New()
	return &PurpurProvider{
		cache: c,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName returns the provider name
func (p *PurpurProvider) GetName() string {
	return "purpur"
}

// GetAvailableVersions returns all available Minecraft versions
func (p *PurpurProvider) GetAvailableVersions(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.purpurmc.org/v2/purpur", nil)
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

	var result struct {
		Versions []string `json:"versions"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Versions, nil
}

// GetLatestBuild returns the latest build for a version
func (p *PurpurProvider) GetLatestBuild(ctx context.Context, version string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s", version), nil)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Builds struct {
			Latest string `json:"latest"`
		} `json:"builds"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Builds.Latest, nil
}

// DownloadJar downloads the Purpur server jar
func (p *PurpurProvider) DownloadJar(ctx context.Context, version, build string) (string, string, string, error) {
	if build == "" || build == "latest" {
		latestBuild, err := p.GetLatestBuild(ctx, version)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to get latest build: %w", err)
		}
		build = latestBuild
	}

	downloadURL := fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s/%s/download", version, build)

	// Purpur doesn't provide checksums in the API, so we'll skip verification
	downloader := cache.NewDownloader(p.cache)
	localPath, err := downloader.DownloadJar(downloadURL, "purpur", version, build, "", "")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to download jar: %w", err)
	}

	return localPath, downloadURL, "", nil
}

// GetDownloadURL returns the download URL
func (p *PurpurProvider) GetDownloadURL(ctx context.Context, version, build string) (string, error) {
	if build == "" || build == "latest" {
		latestBuild, err := p.GetLatestBuild(ctx, version)
		if err != nil {
			return "", err
		}
		build = latestBuild
	}

	return fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s/%s/download", version, build), nil
}

// GetChecksum returns empty (Purpur doesn't provide checksums)
func (p *PurpurProvider) GetChecksum(ctx context.Context, version, build string) (string, string, error) {
	return "", "", nil
}
