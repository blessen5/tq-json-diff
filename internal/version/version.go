package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Version holds the current application version.
	Version = "0.1.0-dev"

	// GitCommit holds the git commit hash at build time.
	GitCommit = "none"

	// BuildDate holds the ISO-8601 build timestamp.
	BuildDate = "unknown"
)

// Info returns formatted version and build metadata.
func Info() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("jdiff version %s\n", Version))
	sb.WriteString(fmt.Sprintf("  commit:     %s\n", GitCommit))
	sb.WriteString(fmt.Sprintf("  built at:   %s\n", BuildDate))
	sb.WriteString(fmt.Sprintf("  go version: %s\n", runtime.Version()))
	sb.WriteString(fmt.Sprintf("  os/arch:    %s/%s", runtime.GOOS, runtime.GOARCH))
	return sb.String()
}

// Short returns just the version number.
func Short() string {
	return Version
}
