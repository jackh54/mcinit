package logs

import (
	"strings"
	"time"
)

// Filter provides log filtering functionality
type Filter struct{}

// NewFilter creates a new Filter instance
func NewFilter() *Filter {
	return &Filter{}
}

// MatchesGrep checks if a line matches a grep pattern (simple contains, case-insensitive)
func (f *Filter) MatchesGrep(line, pattern string) bool {
	if pattern == "" {
		return true
	}
	return strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
}

// IsSince checks if a line's timestamp is after the given time
// This is a simplified version - real implementation would parse log timestamps
func (f *Filter) IsSince(line string, since time.Time) bool {
	if since.IsZero() {
		return true
	}

	// In a real implementation, you would:
	// 1. Parse the timestamp from the log line
	// 2. Compare it with the 'since' time
	// For now, we'll return true (accept all lines)
	
	return true
}

// ParseDuration parses a duration string like "10m", "1h", "30s"
func ParseDuration(durationStr string) (time.Duration, error) {
	return time.ParseDuration(durationStr)
}

// GetSinceTime calculates the time point for "since" filtering
func GetSinceTime(durationStr string) (time.Time, error) {
	if durationStr == "" {
		return time.Time{}, nil
	}

	duration, err := ParseDuration(durationStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().Add(-duration), nil
}

