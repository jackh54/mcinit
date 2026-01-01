package jvmflags

import (
	"fmt"
	"strings"
)

// Minimal returns minimal JVM flags for basic operation
func Minimal(xms, xmx string) []string {
	return []string{
		fmt.Sprintf("-Xms%s", xms),
		fmt.Sprintf("-Xmx%s", xmx),
		"-XX:+UseG1GC",
		"-XX:+ParallelRefProcEnabled",
		"-XX:MaxGCPauseMillis=200",
	}
}

// MinimalString returns minimal flags as a single string
func MinimalString(xms, xmx string) string {
	return strings.Join(Minimal(xms, xmx), " ")
}

// GetMinimalFlags returns minimal flags with memory settings
func GetMinimalFlags(xms, xmx string) []string {
	return Minimal(xms, xmx)
}

