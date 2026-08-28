package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"jdiff/internal/diff"
	"jdiff/internal/version"
)

const (
	// ExitCodeOK indicates the command completed successfully.
	ExitCodeOK = 0
	// ExitCodeError indicates invalid input or operational error.
	ExitCodeError = 1
)

// HelpText is the standard help message displayed by jdiff.
const HelpText = `jdiff - JSON Structural Diff

Usage:
  jdiff <old.json> <new.json>

Commands:
  --help       Show help
  --version    Show version
`

// CLI encapsulates CLI options and standard streams.
type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
}

// New creates a new CLI instance with the specified output streams.
func New(stdout, stderr io.Writer) *CLI {
	return &CLI{
		Stdout: stdout,
		Stderr: stderr,
	}
}

// Run parses arguments and executes the requested command.
// Returns an integer exit code (0 for success, 1 for error).
func (c *CLI) Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprint(c.Stderr, HelpText)
		return ExitCodeError
	}

	// Handle flags and commands
	switch strings.TrimSpace(args[0]) {
	case "--help", "-h", "help":
		fmt.Fprint(c.Stdout, HelpText)
		return ExitCodeOK
	case "--version", "-v", "version":
		fmt.Fprintf(c.Stdout, "jdiff %s\n", version.Short())
		return ExitCodeOK
	}

	// Check if --version or --help appears anywhere in args
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			fmt.Fprint(c.Stdout, HelpText)
			return ExitCodeOK
		}
		if arg == "--version" || arg == "-v" {
			fmt.Fprintf(c.Stdout, "jdiff %s\n", version.Short())
			return ExitCodeOK
		}
	}

	// Check positional arguments and flags
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(c.Stderr, "jdiff: unknown flag %q\n\n%s", arg, HelpText)
			return ExitCodeError
		}
		positional = append(positional, arg)
	}

	if len(positional) < 2 {
		fmt.Fprintln(c.Stderr, "jdiff: missing input files (expected <old.json> <new.json>)")
		return ExitCodeError
	}

	if len(positional) > 2 {
		fmt.Fprintf(c.Stderr, "jdiff: too many arguments provided (expected 2, got %d)\n", len(positional))
		return ExitCodeError
	}

	oldPath := positional[0]
	newPath := positional[1]

	oldBytes, err := os.ReadFile(oldPath)
	if err != nil {
		fmt.Fprintf(c.Stderr, "jdiff: failed to read %s: %v\n", oldPath, err)
		return ExitCodeError
	}

	// Pre-validate JSON for specific error messaging per file
	var testVal any
	if err := json.Unmarshal(oldBytes, &testVal); err != nil {
		fmt.Fprintf(c.Stderr, "jdiff: invalid JSON in %s: %v\n", oldPath, err)
		return ExitCodeError
	}

	newBytes, err := os.ReadFile(newPath)
	if err != nil {
		fmt.Fprintf(c.Stderr, "jdiff: failed to read %s: %v\n", newPath, err)
		return ExitCodeError
	}

	if err := json.Unmarshal(newBytes, &testVal); err != nil {
		fmt.Fprintf(c.Stderr, "jdiff: invalid JSON in %s: %v\n", newPath, err)
		return ExitCodeError
	}

	res, err := diff.CompareBytes(oldBytes, newBytes)
	if err != nil {
		fmt.Fprintf(c.Stderr, "jdiff: error comparing files: %v\n", err)
		return ExitCodeError
	}

	if err := res.Format(c.Stdout); err != nil {
		fmt.Fprintf(c.Stderr, "jdiff: error writing output: %v\n", err)
		return ExitCodeError
	}

	return ExitCodeOK
}
