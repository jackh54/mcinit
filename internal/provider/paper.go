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

// PaperProvider implements Provider for PaperMC
type PaperProvider struct {
	cache       *cache.Cache
	client      *http.Client
	projectName string
}

// NewPaperProvider creates a new PaperProvider
func NewPaperProvider() *PaperProvider {
	c, _ := cache.New()
	return &PaperProvider{
		cache:       c,
		projectName: "paper",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName returns the provider name
func (p *PaperProvider) GetName() string {
	return p.projectName
}

// GetAvailableVersions returns all available Minecraft versions
func (p *PaperProvider) GetAvailableVersions(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", p.projectName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var project PaperProject
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, err
	}

	return project.Versions, nil
}

// GetLatestBuild returns the latest build for a version
func (p *PaperProvider) GetLatestBuild(ctx context.Context, version string) (string, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s", p.projectName, version)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versionInfo PaperVersion
	if err := json.Unmarshal(body, &versionInfo); err != nil {
		return "", err
	}

	if len(versionInfo.Builds) == 0 {
		return "", fmt.Errorf("no builds available for version %s", version)
	}

	// Return the latest build (last in the array)
	return fmt.Sprintf("%d", versionInfo.Builds[len(versionInfo.Builds)-1]), nil
}

// DownloadJar downloads the Paper server jar
func (p *PaperProvider) DownloadJar(ctx context.Context, version, build string) (string, string, string, error) {
	// Resolve "latest" build if needed
	if build == "" || build == "latest" {
		latestBuild, err := p.GetLatestBuild(ctx, version)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to get latest build: %w", err)
		}
		build = latestBuild
	}

	// Get build information
	buildInfo, err := p.getBuildInfo(ctx, version, build)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get build info: %w", err)
	}

	downloadURL := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%s/downloads/%s",
		p.projectName, version, build, buildInfo.Downloads.Application.Name)

	checksum := buildInfo.Downloads.Application.SHA256

	// Download using cache/downloader
	downloader := cache.NewDownloader(p.cache)
	localPath, err := downloader.DownloadJar(downloadURL, p.projectName, version, build, checksum, "sha256")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to download jar: %w", err)
	}

	return localPath, downloadURL, checksum, nil
}

// GetDownloadURL returns the download URL
func (p *PaperProvider) GetDownloadURL(ctx context.Context, version, build string) (string, error) {
	if build == "" || build == "latest" {
		latestBuild, err := p.GetLatestBuild(ctx, version)
		if err != nil {
			return "", err
		}
		build = latestBuild
	}

	buildInfo, err := p.getBuildInfo(ctx, version, build)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%s/downloads/%s",
		p.projectName, version, build, buildInfo.Downloads.Application.Name), nil
}

// GetChecksum returns the checksum
func (p *PaperProvider) GetChecksum(ctx context.Context, version, build string) (string, string, error) {
	if build == "" || build == "latest" {
		latestBuild, err := p.GetLatestBuild(ctx, version)
		if err != nil {
			return "", "", err
		}
		build = latestBuild
	}

	buildInfo, err := p.getBuildInfo(ctx, version, build)
	if err != nil {
		return "", "", err
	}

	return "sha256", buildInfo.Downloads.Application.SHA256, nil
}

// getBuildInfo fetches build information from the API
func (p *PaperProvider) getBuildInfo(ctx context.Context, version, build string) (*PaperBuild, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%s",
		p.projectName, version, build)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var buildInfo PaperBuild
	if err := json.Unmarshal(body, &buildInfo); err != nil {
		return nil, err
	}

	return &buildInfo, nil
}

// PaperMC API structures

type PaperProject struct {
	ProjectID   string   `json:"project_id"`
	ProjectName string   `json:"project_name"`
	Versions    []string `json:"versions"`
}

type PaperVersion struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

type PaperBuild struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Build       int    `json:"build"`
	Time        string `json:"time"`
	Channel     string `json:"channel"`
	Downloads   struct {
		Application struct {
			Name   string `json:"name"`
			SHA256 string `json:"sha256"`
		} `json:"application"`
	} `json:"downloads"`
}
