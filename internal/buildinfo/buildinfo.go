// Package buildinfo holds the version stamped into the binary at release time.
package buildinfo

// Version is set from main (GoReleaser stamps main.version via -ldflags).
var Version = "dev"

// UserAgent identifies mcinit to upstream APIs such as PaperMC Fill.
func UserAgent() string {
	return "mcinit/" + Version + " (github.com/jackh54/mcinit)"
}
