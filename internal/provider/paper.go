package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jackh54/mcinit/internal/buildinfo"
	"github.com/jackh54/mcinit/internal/cache"
)

// fillBaseURL is the PaperMC Fill v3 API. The v2 API (api.papermc.io/v2) returns 410 Gone.
const fillBaseURL = "https://fill.papermc.io/v3"

// fillServerDownload is the downloads key for the server jar in a Fill build.
const fillServerDownload = "server:default"

// PaperProvider implements Provider for PaperMC projects (paper, folia, velocity, waterfall)
// using the Fill v3 API.
type PaperProvider struct {
	cache       *cache.Cache
	client      *http.Client
	projectName string
	baseURL     string
}

// NewPaperProvider creates a new PaperProvider
func NewPaperProvider() *PaperProvider {
	return newFillProvider("paper")
}

func newFillProvider(project string) *PaperProvider {
	c, _ := cache.New()
	return &PaperProvider{
		cache:       c,
		projectName: project,
		baseURL:     fillBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName returns the provider name
func (p *PaperProvider) GetName() string {
	return p.projectName
}

// GetAvailableVersions returns all available versions, newest first, in the order Fill lists them.
func (p *PaperProvider) GetAvailableVersions(ctx context.Context) ([]string, error) {
	var project fillProject
	if err := p.getJSON(ctx, p.projectURL(), &project); err != nil {
		return nil, err
	}
	return flattenVersionGroups(project.Versions)
}

// GetLatestBuild returns the build mcinit would pick for "latest": the newest STABLE build,
// else the newest BETA, else the newest ALPHA.
func (p *PaperProvider) GetLatestBuild(ctx context.Context, version string) (string, error) {
	b, err := p.resolveBuild(ctx, version, "latest")
	if err != nil {
		return "", err
	}
	return strconv.Itoa(b.ID), nil
}

// DownloadJar downloads the server jar and verifies its sha256 against the Fill response.
func (p *PaperProvider) DownloadJar(ctx context.Context, version, build string) (string, string, string, error) {
	b, err := p.resolveBuild(ctx, version, build)
	if err != nil {
		return "", "", "", err
	}
	dl, err := b.serverDownload()
	if err != nil {
		return "", "", "", fmt.Errorf("%s %s build %d: %w", p.projectName, version, b.ID, err)
	}

	if p.cache == nil {
		return "", "", "", errors.New("cache directory unavailable")
	}
	downloader := cache.NewDownloader(p.cache)
	downloader.UserAgent = buildinfo.UserAgent()
	localPath, err := downloader.DownloadJarContext(ctx, dl.URL, p.projectName, version, strconv.Itoa(b.ID), dl.Checksums.SHA256, "sha256")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to download jar: %w", err)
	}

	return localPath, dl.URL, dl.Checksums.SHA256, nil
}

// GetDownloadURL returns the server jar URL for a build
func (p *PaperProvider) GetDownloadURL(ctx context.Context, version, build string) (string, error) {
	b, err := p.resolveBuild(ctx, version, build)
	if err != nil {
		return "", err
	}
	dl, err := b.serverDownload()
	if err != nil {
		return "", err
	}
	return dl.URL, nil
}

// GetChecksum returns the sha256 of the server jar for a build
func (p *PaperProvider) GetChecksum(ctx context.Context, version, build string) (string, string, error) {
	b, err := p.resolveBuild(ctx, version, build)
	if err != nil {
		return "", "", err
	}
	dl, err := b.serverDownload()
	if err != nil {
		return "", "", err
	}
	return "sha256", dl.Checksums.SHA256, nil
}

// resolveBuild fetches the builds for a version and picks one.
func (p *PaperProvider) resolveBuild(ctx context.Context, version, build string) (*fillBuild, error) {
	var builds []fillBuild
	err := p.getJSON(ctx, p.projectURL()+"/versions/"+url.PathEscape(version)+"/builds", &builds)
	if err != nil {
		var apiErr *fillAPIError
		if errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound {
			return nil, fmt.Errorf("%s has no version %q (see https://fill.papermc.io/v3/projects/%s for available versions)",
				p.projectName, version, p.projectName)
		}
		return nil, fmt.Errorf("failed to list %s %s builds: %w", p.projectName, version, err)
	}
	b, err := selectBuild(builds, build)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", p.projectName, version, err)
	}
	return b, nil
}

// channelRank orders release channels; higher is preferred. Unknown channels rank lowest.
func channelRank(channel string) int {
	switch strings.ToUpper(channel) {
	case "RECOMMENDED", "STABLE":
		return 3
	case "BETA":
		return 2
	case "ALPHA":
		return 1
	default:
		return 0
	}
}

// selectBuild picks a build. For "" or "latest" it returns the newest build in the best
// available channel (STABLE, then BETA, then ALPHA). Otherwise build is a build number.
func selectBuild(builds []fillBuild, build string) (*fillBuild, error) {
	if len(builds) == 0 {
		return nil, errors.New("no builds available")
	}

	if build != "" && build != "latest" {
		id, err := strconv.Atoi(build)
		if err != nil {
			return nil, fmt.Errorf("invalid build %q: must be a number or \"latest\"", build)
		}
		for i := range builds {
			if builds[i].ID == id {
				return &builds[i], nil
			}
		}
		return nil, fmt.Errorf("build %d not found", id)
	}

	best := -1
	for i := range builds {
		if best < 0 {
			best = i
			continue
		}
		ri, rb := channelRank(builds[i].Channel), channelRank(builds[best].Channel)
		if ri > rb || (ri == rb && builds[i].ID > builds[best].ID) {
			best = i
		}
	}
	return &builds[best], nil
}

func (p *PaperProvider) projectURL() string {
	return strings.TrimRight(p.baseURL, "/") + "/projects/" + url.PathEscape(p.projectName)
}

// getJSON GETs url with mcinit's User-Agent and decodes a 200 response into v.
func (p *PaperProvider) getJSON(ctx context.Context, rawURL string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", buildinfo.UserAgent())
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		apiErr := &fillAPIError{Status: resp.StatusCode}
		_ = json.Unmarshal(body, apiErr)
		return apiErr
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse response from %s: %w", rawURL, err)
	}
	return nil
}

// flattenVersionGroups turns Fill's {"26.3": ["26.3"], "1.21": ["1.21.11", ...]} into a flat
// list while keeping the API's ordering (Go maps would lose it).
func flattenVersionGroups(raw json.RawMessage) ([]string, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to parse versions: %w", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, errors.New("failed to parse versions: expected an object")
	}
	var out []string
	for dec.More() {
		if _, err := dec.Token(); err != nil { // group key
			return nil, fmt.Errorf("failed to parse versions: %w", err)
		}
		var group []string
		if err := dec.Decode(&group); err != nil {
			return nil, fmt.Errorf("failed to parse versions: %w", err)
		}
		out = append(out, group...)
	}
	return out, nil
}

// Fill v3 API structures

type fillProject struct {
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Versions json.RawMessage `json:"versions"`
}

type fillBuild struct {
	ID        int                     `json:"id"`
	Time      string                  `json:"time"`
	Channel   string                  `json:"channel"`
	Downloads map[string]fillDownload `json:"downloads"`
}

type fillDownload struct {
	Name      string `json:"name"`
	Checksums struct {
		SHA256 string `json:"sha256"`
	} `json:"checksums"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func (b *fillBuild) serverDownload() (*fillDownload, error) {
	dl, ok := b.Downloads[fillServerDownload]
	if !ok || dl.URL == "" {
		return nil, fmt.Errorf("build %d has no %q download", b.ID, fillServerDownload)
	}
	if dl.Checksums.SHA256 == "" {
		return nil, fmt.Errorf("build %d has no sha256 checksum", b.ID)
	}
	return &dl, nil
}

// fillAPIError is a non-200 response from Fill, e.g. {"ok":false,"error":"version_not_found",...}.
type fillAPIError struct {
	Status  int    `json:"-"`
	Code    string `json:"error"`
	Message string `json:"message"`
}

func (e *fillAPIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("PaperMC API returned %d: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("PaperMC API returned %d", e.Status)
}
