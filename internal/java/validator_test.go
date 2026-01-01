package java

import (
	"testing"
)

func TestGetRequiredJavaVersion(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		mcVersion string
		want      int
	}{
		{"Minecraft 1.12", "1.12.2", 8},
		{"Minecraft 1.16", "1.16.5", 8},
		{"Minecraft 1.17", "1.17.1", 16},
		{"Minecraft 1.18", "1.18.2", 17},
		{"Minecraft 1.19", "1.19.4", 17},
		{"Minecraft 1.20", "1.20.4", 17},
		{"Minecraft 1.21", "1.21.4", 21},
		{"Unknown version", "2.0.0", 17}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.GetRequiredJavaVersion(tt.mcVersion)
			if got != tt.want {
				t.Errorf("GetRequiredJavaVersion(%s) = %d, want %d", tt.mcVersion, got, tt.want)
			}
		})
	}
}

func TestGetRecommendedJavaVersion(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		mcVersion string
		want      int
	}{
		{"Minecraft 1.12", "1.12.2", 11}, // Recommend 11 over 8
		{"Minecraft 1.17", "1.17.1", 17},
		{"Minecraft 1.18", "1.18.2", 17},
		{"Minecraft 1.21", "1.21.4", 21},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.GetRecommendedJavaVersion(tt.mcVersion)
			if got != tt.want {
				t.Errorf("GetRecommendedJavaVersion(%s) = %d, want %d", tt.mcVersion, got, tt.want)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name         string
		inst         *Installation
		requiredMajor int
		wantErr      bool
	}{
		{
			name:         "valid version",
			inst:         &Installation{Major: 17, Version: "17.0.1"},
			requiredMajor: 17,
			wantErr:      false,
		},
		{
			name:         "higher version",
			inst:         &Installation{Major: 21, Version: "21.0.1"},
			requiredMajor: 17,
			wantErr:      false,
		},
		{
			name:         "lower version",
			inst:         &Installation{Major: 11, Version: "11.0.11"},
			requiredMajor: 17,
			wantErr:      true,
		},
		{
			name:         "nil installation",
			inst:         nil,
			requiredMajor: 17,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateVersion(tt.inst, tt.requiredMajor)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsCompatible(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		javaMajor int
		mcVersion string
		want      bool
	}{
		{"Java 17 with MC 1.18", 17, "1.18.2", true},
		{"Java 21 with MC 1.21", 21, "1.21.4", true},
		{"Java 21 with MC 1.18", 21, "1.18.2", true}, // Higher is OK
		{"Java 8 with MC 1.18", 8, "1.18.2", false},  // Too low
		{"Java 11 with MC 1.12", 11, "1.12.2", true}, // Higher than required
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsCompatible(tt.javaMajor, tt.mcVersion)
			if got != tt.want {
				t.Errorf("IsCompatible(%d, %s) = %v, want %v", tt.javaMajor, tt.mcVersion, got, tt.want)
			}
		})
	}
}

func TestValidateForMinecraft(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		inst      *Installation
		mcVersion string
		wantErr   bool
	}{
		{
			name:      "Java 17 with MC 1.18",
			inst:      &Installation{Major: 17, Version: "17.0.1"},
			mcVersion: "1.18.2",
			wantErr:   false,
		},
		{
			name:      "Java 21 with MC 1.21",
			inst:      &Installation{Major: 21, Version: "21.0.1"},
			mcVersion: "1.21.4",
			wantErr:   false,
		},
		{
			name:      "Java 8 with MC 1.21",
			inst:      &Installation{Major: 8, Version: "1.8.0"},
			mcVersion: "1.21.4",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateForMinecraft(tt.inst, tt.mcVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateForMinecraft() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

