package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackh54/mcinit/internal/cache"
)

func TestSelectBuild(t *testing.T) {
	b := func(id int, ch string) fillBuild { return fillBuild{ID: id, Channel: ch} }

	tests := []struct {
		name    string
		builds  []fillBuild
		build   string
		wantID  int
		wantErr string
	}{
		{
			name:   "newest stable wins over newer beta and alpha",
			builds: []fillBuild{b(12, "ALPHA"), b(11, "BETA"), b(10, "STABLE"), b(9, "STABLE")},
			build:  "latest",
			wantID: 10,
		},
		{
			name:   "newest stable regardless of response order",
			builds: []fillBuild{b(3, "STABLE"), b(7, "STABLE"), b(5, "STABLE")},
			build:  "",
			wantID: 7,
		},
		{
			name:   "recommended counts as stable",
			builds: []fillBuild{b(4, "BETA"), b(2, "RECOMMENDED")},
			build:  "latest",
			wantID: 2,
		},
		{
			name:   "beta when no stable (Folia 26.x)",
			builds: []fillBuild{b(6, "ALPHA"), b(5, "BETA"), b(4, "BETA"), b(3, "ALPHA")},
			build:  "latest",
			wantID: 5,
		},
		{
			name:   "alpha only",
			builds: []fillBuild{b(4, "ALPHA"), b(6, "ALPHA"), b(5, "ALPHA")},
			build:  "latest",
			wantID: 6,
		},
		{
			name:   "channel match is case-insensitive",
			builds: []fillBuild{b(2, "beta"), b(1, "stable")},
			build:  "latest",
			wantID: 1,
		},
		{
			name:   "unknown channel only used as last resort",
			builds: []fillBuild{b(9, "EXPERIMENTAL"), b(1, "ALPHA")},
			build:  "latest",
			wantID: 1,
		},
		{
			name:   "explicit build number",
			builds: []fillBuild{b(12, "STABLE"), b(11, "ALPHA")},
			build:  "11",
			wantID: 11,
		},
		{
			name:    "explicit build missing",
			builds:  []fillBuild{b(12, "STABLE")},
			build:   "99",
			wantErr: "build 99 not found",
		},
		{
			name:    "explicit build not a number",
			builds:  []fillBuild{b(12, "STABLE")},
			build:   "abc",
			wantErr: "invalid build",
		},
		{
			name:    "no builds",
			builds:  nil,
			build:   "latest",
			wantErr: "no builds available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectBuild(tt.builds, tt.build)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("selectBuild() error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("selectBuild() unexpected error: %v", err)
			}
			if got.ID != tt.wantID {
				t.Errorf("selectBuild() = build %d, want %d", got.ID, tt.wantID)
			}
		})
	}
}

// fakeFill serves a minimal Fill v3 API for one project.
type fakeFill struct {
	srv     *httptest.Server
	project string
	// builds per version; a missing version returns Fill's 404 body.
	builds map[string]string
	jar    []byte

	mu         sync.Mutex
	userAgents []string
}

func newFakeFill(t *testing.T, project string) *fakeFill {
	t.Helper()
	f := &fakeFill{project: project, builds: map[string]string{}, jar: []byte("fake server jar bytes")}
	mux := http.NewServeMux()
	mux.HandleFunc("/v3/projects/"+project, func(w http.ResponseWriter, r *http.Request) {
		f.record(r)
		_, _ = fmt.Fprintf(w, `{"project":{"id":%q,"name":"X"},"versions":{"26.3":["26.3","26.3-rc-3"],"1.21":["1.21.11","1.21.4"],"1.20":["1.20.6"]}}`, project)
	})
	mux.HandleFunc("/v3/projects/"+project+"/versions/", func(w http.ResponseWriter, r *http.Request) {
		f.record(r)
		rest := strings.TrimPrefix(r.URL.Path, "/v3/projects/"+project+"/versions/")
		version := strings.TrimSuffix(rest, "/builds")
		body, ok := f.builds[version]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"ok":false,"error":"version_not_found","message":"No version was found with the given identifier."}`))
			return
		}
		_, _ = w.Write([]byte(body))
	})
	mux.HandleFunc("/objects/", func(w http.ResponseWriter, r *http.Request) {
		f.record(r)
		_, _ = w.Write(f.jar)
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeFill) record(r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.userAgents = append(f.userAgents, r.UserAgent())
}

// buildJSON renders one Fill build whose jar is served by the fake server.
func (f *fakeFill) buildJSON(id int, channel, sha string) string {
	name := fmt.Sprintf("%s-%d.jar", f.project, id)
	return fmt.Sprintf(`{"id":%d,"time":"2026-10-29T23:21:46Z","channel":%q,"commits":[],"downloads":{"server:default":{"name":%q,"checksums":{"sha256":%q},"size":21,"url":"%s/objects/%s/%s"}}}`,
		id, channel, name, sha, f.srv.URL, sha, name)
}

func (f *fakeFill) provider(t *testing.T) *PaperProvider {
	t.Helper()
	c, err := cache.NewWithDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return &PaperProvider{
		cache:       c,
		client:      &http.Client{Timeout: 5 * time.Second},
		projectName: f.project,
		baseURL:     f.srv.URL + "/v3",
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func TestFillDownloadJar(t *testing.T) {
	f := newFakeFill(t, "folia")
	good := sha256Hex(f.jar)
	f.builds["26.3"] = "[" + f.buildJSON(6, "ALPHA", good) + "," + f.buildJSON(5, "BETA", good) + "," + f.buildJSON(4, "BETA", good) + "]"
	p := f.provider(t)

	path, url, sum, err := p.DownloadJar(context.Background(), "26.3", "latest")
	if err != nil {
		t.Fatalf("DownloadJar() error: %v", err)
	}
	if !strings.HasSuffix(url, "/folia-5.jar") {
		t.Errorf("download URL = %s, want the BETA build 5 jar", url)
	}
	if sum != good {
		t.Errorf("checksum = %s, want %s", sum, good)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading downloaded jar: %v", err)
	}
	if string(data) != string(f.jar) {
		t.Errorf("jar contents = %q, want %q", data, f.jar)
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.userAgents) < 2 {
		t.Fatalf("expected API and download requests, got %d", len(f.userAgents))
	}
	for _, ua := range f.userAgents {
		if !strings.HasPrefix(ua, "mcinit/") || !strings.Contains(ua, "github.com/jackh54/mcinit") {
			t.Errorf("User-Agent = %q, want mcinit/<version> (github.com/jackh54/mcinit)", ua)
		}
	}
}

func TestFillDownloadJarChecksumMismatch(t *testing.T) {
	f := newFakeFill(t, "paper")
	wrong := sha256Hex([]byte("some other jar"))
	f.builds["1.21.11"] = "[" + f.buildJSON(132, "STABLE", wrong) + "]"
	p := f.provider(t)

	_, _, _, err := p.DownloadJar(context.Background(), "1.21.11", "latest")
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("DownloadJar() error = %v, want checksum mismatch", err)
	}

	// Nothing may be left in the cache: neither the jar nor the temp file.
	jarPath := p.cache.GetJarPath("paper", "1.21.11", "132")
	if _, statErr := os.Stat(jarPath); !os.IsNotExist(statErr) {
		t.Errorf("jar left in cache after mismatch: %s", jarPath)
	}
	if _, statErr := os.Stat(jarPath + ".tmp"); !os.IsNotExist(statErr) {
		t.Errorf("temp file left in cache after mismatch: %s.tmp", jarPath)
	}
	entries, _ := filepath.Glob(filepath.Join(p.cache.GetBaseDir(), "**", "*.jar"))
	if len(entries) != 0 {
		t.Errorf("unexpected jars in cache: %v", entries)
	}
}

func TestFillMissingVersion(t *testing.T) {
	f := newFakeFill(t, "paper")
	p := f.provider(t)

	for name, call := range map[string]func() error{
		"DownloadJar":    func() error { _, _, _, err := p.DownloadJar(context.Background(), "9.9.9", "latest"); return err },
		"GetLatestBuild": func() error { _, err := p.GetLatestBuild(context.Background(), "9.9.9"); return err },
		"GetChecksum":    func() error { _, _, err := p.GetChecksum(context.Background(), "9.9.9", "latest"); return err },
	} {
		t.Run(name, func(t *testing.T) {
			err := call()
			if err == nil {
				t.Fatal("expected an error for a missing version")
			}
			msg := err.Error()
			if !strings.Contains(msg, `paper has no version "9.9.9"`) {
				t.Errorf("error = %q, want it to name the project and version", msg)
			}
		})
	}
}

func TestFillNoBuildsForVersion(t *testing.T) {
	f := newFakeFill(t, "paper")
	f.builds["1.21.4"] = "[]"
	p := f.provider(t)

	_, err := p.GetLatestBuild(context.Background(), "1.21.4")
	if err == nil || !strings.Contains(err.Error(), "paper 1.21.4: no builds available") {
		t.Fatalf("GetLatestBuild() error = %v, want no builds available", err)
	}
}

func TestFillServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	p := &PaperProvider{client: srv.Client(), projectName: "paper", baseURL: srv.URL + "/v3"}

	_, err := p.GetLatestBuild(context.Background(), "1.21.11")
	if err == nil || !strings.Contains(err.Error(), "returned 500") {
		t.Fatalf("GetLatestBuild() error = %v, want status 500", err)
	}
}

func TestFillGetAvailableVersionsKeepsOrder(t *testing.T) {
	f := newFakeFill(t, "paper")
	p := f.provider(t)

	got, err := p.GetAvailableVersions(context.Background())
	if err != nil {
		t.Fatalf("GetAvailableVersions() error: %v", err)
	}
	want := []string{"26.3", "26.3-rc-3", "1.21.11", "1.21.4", "1.20.6"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAvailableVersions() = %v, want %v", got, want)
	}
}

func TestFillProvidersShareImplementation(t *testing.T) {
	for _, name := range []string{"paper", "folia", "velocity", "waterfall"} {
		prov, err := Get(name)
		if err != nil {
			t.Fatal(err)
		}
		var pp *PaperProvider
		switch v := prov.(type) {
		case *PaperProvider:
			pp = v
		case *FoliaProvider:
			pp = v.PaperProvider
		case *VelocityProvider:
			pp = v.PaperProvider
		case *WaterfallProvider:
			pp = v.PaperProvider
		default:
			t.Fatalf("%s: unexpected provider type %T", name, prov)
		}
		if pp.baseURL != fillBaseURL || pp.projectName != name {
			t.Errorf("%s: baseURL=%q project=%q, want %q/%q", name, pp.baseURL, pp.projectName, fillBaseURL, name)
		}
	}
}
