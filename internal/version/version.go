package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current SemVer version. Set via ldflags.
	Version = "0.1.0-alpha"

	// Commit is the git commit SHA. Set via ldflags.
	Commit = "unknown"

	// BuildDate is the RFC3339 date when the binary was built. Set via ldflags.
	BuildDate = "unknown"
)

// Info represents the version info payload.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// GetInfo returns the unified version information.
func GetInfo() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a formatted version string.
func String() string {
	return fmt.Sprintf("PromptEngine v%s (commit: %s, built: %s, go: %s)", Version, Commit, BuildDate, runtime.Version())
}
