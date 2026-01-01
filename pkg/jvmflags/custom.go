package jvmflags

import (
	"fmt"
	"strings"
)

// Custom parses and validates custom JVM flags
func Custom(xms, xmx string, customFlags []string) []string {
	flags := []string{
		fmt.Sprintf("-Xms%s", xms),
		fmt.Sprintf("-Xmx%s", xmx),
	}

	// Add custom flags
	flags = append(flags, customFlags...)

	return flags
}

// CustomString returns custom flags as a single string
func CustomString(xms, xmx string, customFlags []string) string {
	return strings.Join(Custom(xms, xmx, customFlags), " ")
}

// ParseCustomFlags parses a space-separated string of custom flags
func ParseCustomFlags(flagsString string) []string {
	if flagsString == "" {
		return []string{}
	}

	// Split by spaces, handling quoted strings
	var flags []string
	var current strings.Builder
	inQuote := false

	for _, char := range flagsString {
		switch char {
		case '"':
			inQuote = !inQuote
		case ' ':
			if inQuote {
				current.WriteRune(char)
			} else if current.Len() > 0 {
				flags = append(flags, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		flags = append(flags, current.String())
	}

	return flags
}

// ValidateFlags performs basic validation on JVM flags
func ValidateFlags(flags []string) error {
	// Basic validation: ensure flags start with - or are valid format
	for _, flag := range flags {
		if len(flag) == 0 {
			continue
		}
		if !strings.HasPrefix(flag, "-") && !strings.HasPrefix(flag, "--") {
			return fmt.Errorf("invalid flag format: %s (flags must start with - or --)", flag)
		}
	}
	return nil
}

// GetCustomFlags returns custom flags with memory settings
func GetCustomFlags(xms, xmx string, customFlags []string) []string {
	return Custom(xms, xmx, customFlags)
}

