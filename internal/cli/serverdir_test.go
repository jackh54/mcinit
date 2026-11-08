package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveConfig(t *testing.T) {
	serverDir := t.TempDir()
	cfg := filepath.Join(serverDir, "custom.json")
	if err := os.WriteFile(cfg, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Resolve symlinks (macOS /var -> /private/var) so comparisons are stable.
	serverDir, _ = filepath.EvalSymlinks(serverDir)
	cfg = filepath.Join(serverDir, "custom.json")

	otherDir := t.TempDir()
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(otherDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prevWD) })
	otherDir, _ = os.Getwd()

	tests := []struct {
		name    string
		flag    string
		wantCfg string
		wantDir string
		wantErr string
	}{
		{name: "no flag uses cwd", flag: "", wantCfg: filepath.Join(otherDir, "mcinit.json"), wantDir: otherDir},
		{name: "absolute config elsewhere", flag: cfg, wantCfg: cfg, wantDir: serverDir},
		{name: "relative config", flag: mustRel(t, otherDir, cfg), wantCfg: cfg, wantDir: serverDir},
		{name: "missing config", flag: filepath.Join(serverDir, "nope.json"), wantErr: "config file not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := cfgFile
			cfgFile = tt.flag
			t.Cleanup(func() { cfgFile = old })

			gotCfg, gotDir, err := resolveConfig()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("resolveConfig() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveConfig() error: %v", err)
			}
			if gotCfg != tt.wantCfg || gotDir != tt.wantDir {
				t.Errorf("resolveConfig() = (%s, %s), want (%s, %s)", gotCfg, gotDir, tt.wantCfg, tt.wantDir)
			}
		})
	}
}

func mustRel(t *testing.T, base, target string) string {
	t.Helper()
	rel, err := filepath.Rel(base, target)
	if err != nil {
		t.Fatal(err)
	}
	return rel
}
