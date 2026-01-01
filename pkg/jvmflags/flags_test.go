package jvmflags

import (
	"strings"
	"testing"
)

func TestAikar(t *testing.T) {
	flags := Aikar("2G", "4G")

	if len(flags) == 0 {
		t.Error("Aikar() returned no flags")
	}

	// Check for Xms and Xmx
	hasXms := false
	hasXmx := false
	for _, flag := range flags {
		if strings.Contains(flag, "-Xms2G") {
			hasXms = true
		}
		if strings.Contains(flag, "-Xmx4G") {
			hasXmx = true
		}
	}

	if !hasXms {
		t.Error("Aikar flags missing -Xms2G")
	}
	if !hasXmx {
		t.Error("Aikar flags missing -Xmx4G")
	}

	// Check for G1GC
	hasG1GC := false
	for _, flag := range flags {
		if flag == "-XX:+UseG1GC" {
			hasG1GC = true
			break
		}
	}
	if !hasG1GC {
		t.Error("Aikar flags missing G1GC")
	}
}

func TestMinimal(t *testing.T) {
	flags := Minimal("1G", "2G")

	if len(flags) == 0 {
		t.Error("Minimal() returned no flags")
	}

	// Should have at least Xms, Xmx, and UseG1GC
	if len(flags) < 3 {
		t.Errorf("Minimal() returned %d flags, expected at least 3", len(flags))
	}

	// Check for basic flags
	hasXms := false
	hasXmx := false
	for _, flag := range flags {
		if strings.Contains(flag, "-Xms1G") {
			hasXms = true
		}
		if strings.Contains(flag, "-Xmx2G") {
			hasXmx = true
		}
	}

	if !hasXms || !hasXmx {
		t.Error("Minimal flags missing Xms or Xmx")
	}
}

func TestCustom(t *testing.T) {
	customFlags := []string{"-XX:+PrintGCDetails", "-XX:MaxMetaspaceSize=256M"}
	flags := Custom("2G", "4G", customFlags)

	if len(flags) < 2+len(customFlags) {
		t.Error("Custom() didn't include all flags")
	}

	// Should include custom flags
	foundCustom := 0
	for _, flag := range flags {
		for _, custom := range customFlags {
			if flag == custom {
				foundCustom++
			}
		}
	}

	if foundCustom != len(customFlags) {
		t.Errorf("Custom() included %d custom flags, expected %d", foundCustom, len(customFlags))
	}
}

func TestParseCustomFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty string", "", 0},
		{"single flag", "-XX:+PrintGC", 1},
		{"multiple flags", "-XX:+PrintGC -Xlog:gc", 2},
		{"flags with spaces", "-XX:+PrintGC  -Xlog:gc", 2},
		{"quoted flag", "\"-Dfoo=bar baz\"", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := ParseCustomFlags(tt.input)
			if len(flags) != tt.expected {
				t.Errorf("ParseCustomFlags() returned %d flags, expected %d", len(flags), tt.expected)
			}
		})
	}
}

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name    string
		flags   []string
		wantErr bool
	}{
		{"valid flags", []string{"-Xmx4G", "-XX:+UseG1GC"}, false},
		{"invalid flag", []string{"invalid"}, true},
		{"empty flag", []string{""}, false}, // Empty flags are ignored
		{"mixed valid/invalid", []string{"-Xmx4G", "invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFlags(tt.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		preset      string
		xms         string
		xmx         string
		customFlags []string
		wantErr     bool
	}{
		{"aikar preset", "aikar", "2G", "4G", nil, false},
		{"minimal preset", "minimal", "2G", "4G", nil, false},
		{"custom preset", "custom", "2G", "4G", []string{"-XX:+PrintGC"}, false},
		{"invalid preset", "invalid", "2G", "4G", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, err := Get(tt.preset, tt.xms, tt.xmx, tt.customFlags)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(flags) == 0 {
				t.Error("Get() returned no flags")
			}
		})
	}
}

func TestValidatePreset(t *testing.T) {
	tests := []struct {
		preset string
		want   bool
	}{
		{"aikar", true},
		{"minimal", true},
		{"custom", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			if got := ValidatePreset(tt.preset); got != tt.want {
				t.Errorf("ValidatePreset(%s) = %v, want %v", tt.preset, got, tt.want)
			}
		})
	}
}

func TestListPresets(t *testing.T) {
	presets := ListPresets()
	
	if len(presets) != 3 {
		t.Errorf("ListPresets() returned %d presets, expected 3", len(presets))
	}

	expected := map[string]bool{"aikar": true, "minimal": true, "custom": true}
	for _, preset := range presets {
		if !expected[preset] {
			t.Errorf("Unexpected preset in list: %s", preset)
		}
	}
}

