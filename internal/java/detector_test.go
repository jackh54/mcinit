package java

import (
	"testing"
)

func TestParseJavaVersion(t *testing.T) {
	tests := []struct {
		name          string
		output        string
		wantVersion   string
		wantMajor     int
		wantErr       bool
	}{
		{
			name:        "Java 8",
			output:      `java version "1.8.0_292"`,
			wantVersion: "1.8.0_292",
			wantMajor:   8,
			wantErr:     false,
		},
		{
			name:        "Java 11",
			output:      `openjdk version "11.0.11" 2021-04-20`,
			wantVersion: "11.0.11",
			wantMajor:   11,
			wantErr:     false,
		},
		{
			name:        "Java 17",
			output:      `openjdk version "17.0.1" 2021-10-19 LTS`,
			wantVersion: "17.0.1",
			wantMajor:   17,
			wantErr:     false,
		},
		{
			name:        "Java 21",
			output:      `openjdk version "21.0.1" 2023-10-17`,
			wantVersion: "21.0.1",
			wantMajor:   21,
			wantErr:     false,
		},
		{
			name:    "invalid output",
			output:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, major, _, _, err := parseJavaVersion(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJavaVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if version != tt.wantVersion {
					t.Errorf("parseJavaVersion() version = %s, want %s", version, tt.wantVersion)
				}
				if major != tt.wantMajor {
					t.Errorf("parseJavaVersion() major = %d, want %d", major, tt.wantMajor)
				}
			}
		})
	}
}

func TestInstallationGetUptime(t *testing.T) {
	inst := &Installation{
		Path:    "/usr/bin/java",
		Version: "17.0.1",
		Major:   17,
		Minor:   0,
		Patch:   1,
	}

	if inst.Path != "/usr/bin/java" {
		t.Errorf("Installation.Path = %s, want /usr/bin/java", inst.Path)
	}

	if inst.Major != 17 {
		t.Errorf("Installation.Major = %d, want 17", inst.Major)
	}
}

