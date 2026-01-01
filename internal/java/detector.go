package java

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jackh54/mcinit/internal/utils"
)

// Installation represents a Java installation
type Installation struct {
	Path    string
	Version string
	Major   int
	Minor   int
	Patch   int
}

// Detector handles Java detection
type Detector struct{}

// NewDetector creates a new Detector instance
func NewDetector() *Detector {
	return &Detector{}
}

// Detect finds Java installations on the system
// Priority: PATH first, then common locations
func (d *Detector) Detect() ([]*Installation, error) {
	var installations []*Installation
	seen := make(map[string]bool)

	// 1. Check JAVA_HOME environment variable
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		javaPath := filepath.Join(javaHome, "bin", utils.GetExecutableName("java"))
		if utils.PathExists(javaPath) {
			if inst, err := d.checkJava(javaPath); err == nil && !seen[inst.Path] {
				installations = append(installations, inst)
				seen[inst.Path] = true
			}
		}
	}

	// 2. Search PATH for java executable
	if javaPath, err := exec.LookPath("java"); err == nil {
		if absPath, err := filepath.Abs(javaPath); err == nil {
			if inst, err := d.checkJava(absPath); err == nil && !seen[inst.Path] {
				installations = append(installations, inst)
				seen[inst.Path] = true
			}
		}
	}

	// 3. Check common locations based on platform
	commonPaths := d.getCommonPaths()
	for _, basePath := range commonPaths {
		if !utils.PathExists(basePath) {
			continue
		}

		// Search for Java installations in this path
		found := d.searchPath(basePath)
		for _, inst := range found {
			if !seen[inst.Path] {
				installations = append(installations, inst)
				seen[inst.Path] = true
			}
		}
	}

	if len(installations) == 0 {
		return nil, fmt.Errorf("no Java installation found")
	}

	return installations, nil
}

// DetectVersion detects the version of a specific Java installation
func (d *Detector) DetectVersion(javaPath string) (*Installation, error) {
	return d.checkJava(javaPath)
}

// FindByMajorVersion finds a Java installation with a specific major version
func (d *Detector) FindByMajorVersion(majorVersion int) (*Installation, error) {
	installations, err := d.Detect()
	if err != nil {
		return nil, err
	}

	for _, inst := range installations {
		if inst.Major == majorVersion {
			return inst, nil
		}
	}

	return nil, fmt.Errorf("Java %d not found", majorVersion)
}

// FindBest finds the best Java installation (highest version)
func (d *Detector) FindBest() (*Installation, error) {
	installations, err := d.Detect()
	if err != nil {
		return nil, err
	}

	if len(installations) == 0 {
		return nil, fmt.Errorf("no Java installation found")
	}

	best := installations[0]
	for _, inst := range installations[1:] {
		if inst.Major > best.Major ||
			(inst.Major == best.Major && inst.Minor > best.Minor) ||
			(inst.Major == best.Major && inst.Minor == best.Minor && inst.Patch > best.Patch) {
			best = inst
		}
	}

	return best, nil
}

// checkJava verifies a Java executable and extracts its version
func (d *Detector) checkJava(javaPath string) (*Installation, error) {
	if !utils.PathExists(javaPath) {
		return nil, fmt.Errorf("java executable not found: %s", javaPath)
	}

	// Run java -version
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run java -version: %w", err)
	}

	// Parse version from output
	version, major, minor, patch, err := parseJavaVersion(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Java version: %w", err)
	}

	return &Installation{
		Path:    javaPath,
		Version: version,
		Major:   major,
		Minor:   minor,
		Patch:   patch,
	}, nil
}

// getCommonPaths returns common Java installation paths for the current platform
func (d *Detector) getCommonPaths() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"C:\\Program Files\\Java",
			"C:\\Program Files (x86)\\Java",
			"C:\\Program Files\\Eclipse Adoptium",
			"C:\\Program Files\\Microsoft\\jdk",
		}

	case "darwin":
		return []string{
			"/Library/Java/JavaVirtualMachines",
			"/System/Library/Java/JavaVirtualMachines",
		}

	case "linux":
		return []string{
			"/usr/lib/jvm",
			"/usr/java",
			"/opt/java",
			"/opt/jdk",
		}

	default:
		return []string{}
	}
}

// searchPath searches for Java installations in a directory
func (d *Detector) searchPath(basePath string) []*Installation {
	var installations []*Installation

	// Walk the directory tree
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Look for java executable
		if info.IsDir() {
			return nil
		}

		if info.Name() == "java" || info.Name() == "java.exe" {
			if inst, err := d.checkJava(path); err == nil {
				installations = append(installations, inst)
			}
		}

		return nil
	})

	return installations
}

// parseJavaVersion parses Java version from java -version output
func parseJavaVersion(output string) (version string, major, minor, patch int, err error) {
	// Java version output format varies:
	// Java 8: java version "1.8.0_292"
	// Java 11+: openjdk version "11.0.11" 2021-04-20
	// Java 17+: openjdk version "17.0.1" 2021-10-19 LTS

	// Try to match version string
	versionRegex := regexp.MustCompile(`version "([^"]+)"`)
	matches := versionRegex.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", 0, 0, 0, fmt.Errorf("could not parse Java version from output")
	}

	version = matches[1]

	// Parse version components
	// Handle both 1.8.0 (old format) and 11.0.11 (new format)
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return version, 0, 0, 0, fmt.Errorf("invalid version format: %s", version)
	}

	// Check if old format (1.x.y)
	if parts[0] == "1" && len(parts) >= 2 {
		// Java 8 and earlier: 1.8.0_292
		fmt.Sscanf(parts[1], "%d", &major)
		if len(parts) >= 3 {
			// Handle patch with optional build number (e.g., 0_292)
			patchStr := strings.Split(parts[2], "_")[0]
			fmt.Sscanf(patchStr, "%d", &minor)
		}
	} else {
		// Java 9+: 11.0.11, 17.0.1, etc.
		fmt.Sscanf(parts[0], "%d", &major)
		if len(parts) >= 2 {
			fmt.Sscanf(parts[1], "%d", &minor)
		}
		if len(parts) >= 3 {
			// Handle patch with optional build info
			patchStr := strings.Split(parts[2], " ")[0]
			patchStr = strings.Split(patchStr, "-")[0]
			patchStr = strings.Split(patchStr, "+")[0]
			fmt.Sscanf(patchStr, "%d", &patch)
		}
	}

	return version, major, minor, patch, nil
}
