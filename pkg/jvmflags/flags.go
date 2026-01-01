package jvmflags

import (
	"fmt"
)

// Get returns JVM flags based on preset name
func Get(preset, xms, xmx string, customFlags []string) ([]string, error) {
	switch preset {
	case "aikar":
		return GetAikarFlags(xms, xmx), nil
	case "minimal":
		return GetMinimalFlags(xms, xmx), nil
	case "custom":
		if err := ValidateFlags(customFlags); err != nil {
			return nil, fmt.Errorf("invalid custom flags: %w", err)
		}
		return GetCustomFlags(xms, xmx, customFlags), nil
	default:
		return nil, fmt.Errorf("unknown JVM flags preset: %s", preset)
	}
}

// GetAsString returns JVM flags as a single space-separated string
func GetAsString(preset, xms, xmx string, customFlags []string) (string, error) {
	flags, err := Get(preset, xms, xmx, customFlags)
	if err != nil {
		return "", err
	}

	return Join(flags), nil
}

// Join joins flags into a single string
func Join(flags []string) string {
	result := ""
	for i, flag := range flags {
		if i > 0 {
			result += " "
		}
		result += flag
	}
	return result
}

// ValidatePreset checks if a preset name is valid
func ValidatePreset(preset string) bool {
	return preset == "aikar" || preset == "minimal" || preset == "custom"
}

// ListPresets returns all available presets
func ListPresets() []string {
	return []string{"aikar", "minimal", "custom"}
}

