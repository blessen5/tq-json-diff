package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Version holds the current application version.
	Version = "v0.5.0"

	// GitCommit holds the git commit hash at build time.
	GitCommit = "none"

	// BuildDate holds the ISO-8601 build timestamp.
	BuildDate = "unknown"
)

// Info returns formatted version and build metadata.
func Info() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("jdiff %s\n", Version))
	sb.WriteString(fmt.Sprintf("  commit:     %s\n", GitCommit))
	sb.WriteString(fmt.Sprintf("  built at:   %s\n", BuildDate))
	sb.WriteString(fmt.Sprintf("  go version: %s\n", runtime.Version()))
	sb.WriteString(fmt.Sprintf("  os/arch:    %s/%s", runtime.GOOS, runtime.GOARCH))
	return sb.String()
}

// Short returns just the version string.
func Short() string {
	return Version
}
