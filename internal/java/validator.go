package java

import (
	"fmt"
)

// Validator validates Java installations
type Validator struct{}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateVersion checks if a Java version meets requirements
func (v *Validator) ValidateVersion(inst *Installation, requiredMajor int) error {
	if inst == nil {
		return fmt.Errorf("no Java installation provided")
	}

	if inst.Major < requiredMajor {
		return fmt.Errorf("java %d or higher required, found Java %d (%s)", requiredMajor, inst.Major, inst.Version)
	}

	return nil
}

// ValidateForMinecraft checks if Java version is suitable for a Minecraft version
func (v *Validator) ValidateForMinecraft(inst *Installation, mcVersion string) error {
	if inst == nil {
		return fmt.Errorf("no Java installation provided")
	}

	// Get required Java version for Minecraft version
	requiredMajor := v.GetRequiredJavaVersion(mcVersion)

	return v.ValidateVersion(inst, requiredMajor)
}

// GetRequiredJavaVersion returns the minimum Java version required for a Minecraft version
func (v *Validator) GetRequiredJavaVersion(mcVersion string) int {
	// Minecraft version to Java version mapping:
	// 1.12-1.16.5: Java 8+
	// 1.17-1.17.1: Java 16+
	// 1.18-1.20.4: Java 17+
	// 1.20.5+: Java 21+ (for latest versions)

	// Simple version comparison based on major version
	// For a production system, this should be more robust
	if len(mcVersion) == 0 {
		return 17 // Default to Java 17
	}

	// Extract major version
	var major, minor int
	_, _ = fmt.Sscanf(mcVersion, "%d.%d", &major, &minor)

	// Minecraft 1.12-1.16
	if major == 1 && minor >= 12 && minor <= 16 {
		return 8
	}

	// Minecraft 1.17
	if major == 1 && minor == 17 {
		return 16
	}

	// Minecraft 1.18-1.20.4
	if major == 1 && minor >= 18 && minor <= 20 {
		return 17
	}

	// Minecraft 1.20.5+ (assume latest needs Java 21)
	if major == 1 && minor >= 21 {
		return 21
	}

	// Default to Java 17 for unknown versions
	return 17
}

// GetRecommendedJavaVersion returns the recommended Java version for a Minecraft version
func (v *Validator) GetRecommendedJavaVersion(mcVersion string) int {
	required := v.GetRequiredJavaVersion(mcVersion)

	// Recommend the same or next LTS version
	switch required {
	case 8:
		return 11 // Recommend Java 11 over 8
	case 16:
		return 17 // Recommend Java 17 over 16
	case 17:
		return 17
	case 21:
		return 21
	default:
		return 17
	}
}

// IsCompatible checks if a Java version is compatible with a Minecraft version
func (v *Validator) IsCompatible(javaMajor int, mcVersion string) bool {
	required := v.GetRequiredJavaVersion(mcVersion)
	return javaMajor >= required
}

