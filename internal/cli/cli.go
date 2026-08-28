package cli

import (
	"fmt"
	"io"
	"strings"

	"jdiff/internal/version"
)

const (
	// ExitCodeOK indicates successful execution with no differences or informational command completed.
	ExitCodeOK = 0
	// ExitCodeDiffFound indicates differences were found between JSON files (reserved for Phase 2+).
	ExitCodeDiffFound = 1
	// ExitCodeError indicates an error occurred (invalid arguments, I/O error, invalid JSON).
	ExitCodeError = 2
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
// Returns an integer exit code.
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
		fmt.Fprintf(c.Stdout, "jdiff version %s\n", version.Short())
		return ExitCodeOK
	}

	// Check if --version or --help appears anywhere in args
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			fmt.Fprint(c.Stdout, HelpText)
			return ExitCodeOK
		}
		if arg == "--version" || arg == "-v" {
			fmt.Fprintf(c.Stdout, "jdiff version %s\n", version.Short())
			return ExitCodeOK
		}
	}

	// Check positional arguments
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(c.Stderr, "Error: unknown flag %q\n\n%s", arg, HelpText)
			return ExitCodeError
		}
		positional = append(positional, arg)
	}

	if len(positional) < 2 {
		fmt.Fprintf(c.Stderr, "Error: two JSON file paths are required (<old.json> <new.json>)\n\n%s", HelpText)
		return ExitCodeError
	}

	if len(positional) > 2 {
		fmt.Fprintf(c.Stderr, "Error: too many arguments provided (expected 2, got %d)\n\n%s", len(positional), HelpText)
		return ExitCodeError
	}

	oldPath := positional[0]
	newPath := positional[1]

	// Phase 1 placeholder: JSON diffing will be implemented in Phase 2
	fmt.Fprintf(c.Stdout, "Comparing %s and %s...\n(JSON comparison engine will be initialized in Phase 2)\n", oldPath, newPath)
	return ExitCodeOK
}
