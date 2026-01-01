package jvmflags

import (
	"fmt"
	"strings"
)

// Aikar returns Aikar's famous flags for optimal GC performance
// Source: https://docs.papermc.io/paper/aikars-flags
func Aikar(xms, xmx string) []string {
	return []string{
		fmt.Sprintf("-Xms%s", xms),
		fmt.Sprintf("-Xmx%s", xmx),
		"-XX:+UseG1GC",
		"-XX:+ParallelRefProcEnabled",
		"-XX:MaxGCPauseMillis=200",
		"-XX:+UnlockExperimentalVMOptions",
		"-XX:+DisableExplicitGC",
		"-XX:+AlwaysPreTouch",
		"-XX:G1NewSizePercent=30",
		"-XX:G1MaxNewSizePercent=40",
		"-XX:G1HeapRegionSize=8M",
		"-XX:G1ReservePercent=20",
		"-XX:G1HeapWastePercent=5",
		"-XX:G1MixedGCCountTarget=4",
		"-XX:InitiatingHeapOccupancyPercent=15",
		"-XX:G1MixedGCLiveThresholdPercent=90",
		"-XX:G1RSetUpdatingPauseTimePercent=5",
		"-XX:SurvivorRatio=32",
		"-XX:+PerfDisableSharedMem",
		"-XX:MaxTenuringThreshold=1",
		"-Dusing.aikars.flags=https://mcflags.emc.gs",
		"-Daikars.new.flags=true",
	}
}

// AikarString returns Aikar's flags as a single string
func AikarString(xms, xmx string) string {
	return strings.Join(Aikar(xms, xmx), " ")
}

// GetAikarFlags returns Aikar's flags with memory settings
func GetAikarFlags(xms, xmx string) []string {
	return Aikar(xms, xmx)
}

