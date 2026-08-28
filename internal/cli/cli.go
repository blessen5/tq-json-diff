package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"jdiff/internal/config"
	"jdiff/internal/diff"
	"jdiff/internal/matcher"
	"jdiff/internal/render"
	"jdiff/internal/version"
)

const (
	// ExitCodeOK indicates successful execution.
	ExitCodeOK = 0
	// ExitCodeError indicates an error occurred during execution.
	ExitCodeError = 1
)

const usage = `jdiff - JSON Structural Diff

Usage:
  jdiff [options] <old.json> <new.json>

Options:
  --help            Show help
  --version         Show version
  --no-color        Disable colored output
  --compact         Display compact diff output
  --verbose         Display additional comparison information
  --summary         Display only the change summary
  --ignore <path>   Ignore a JSON path (can be specified multiple times)
  --config <file>   Use a configuration file (defaults to .jdiff.json)
  --show-config     Show active comparison configuration
`

// CLI manages command-line interface execution and I/O streams.
type CLI struct {
	stdout io.Writer
	stderr io.Writer
}

// New creates a new CLI instance with given standard output and error writers.
func New(stdout, stderr io.Writer) *CLI {
	return &CLI{
		stdout: stdout,
		stderr: stderr,
	}
}

// Run parses arguments and executes the diff operation.
func (c *CLI) Run(args []string) int {
	var (
		noColor     bool
		compact     bool
		verbose     bool
		summary     bool
		showConfig  bool
		configFile  string
		cliIgnores  []string
		positionals []string
	)

	// Check if help or version was requested
	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "help" {
			fmt.Fprint(c.stdout, usage)
			return ExitCodeOK
		}
		if arg == "--version" || arg == "-v" || arg == "version" {
			fmt.Fprintln(c.stdout, fmt.Sprintf("jdiff %s", version.Short()))
			return ExitCodeOK
		}
	}

	// Parse flags and positionals
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--no-color":
			noColor = true
		case arg == "--compact":
			compact = true
		case arg == "--verbose":
			verbose = true
		case arg == "--summary":
			summary = true
		case arg == "--show-config":
			showConfig = true
		case arg == "--config":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --config")
				return ExitCodeError
			}
			i++
			configFile = args[i]
		case arg == "--ignore":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --ignore")
				return ExitCodeError
			}
			i++
			cliIgnores = append(cliIgnores, args[i])
		case strings.HasPrefix(arg, "--ignore="):
			cliIgnores = append(cliIgnores, strings.TrimPrefix(arg, "--ignore="))
		case strings.HasPrefix(arg, "--config="):
			configFile = strings.TrimPrefix(arg, "--config=")
		case strings.HasPrefix(arg, "-"):
			fmt.Fprintf(c.stderr, "jdiff: unknown flag %q\n", arg)
			return ExitCodeError
		default:
			positionals = append(positionals, arg)
		}
	}

	// Load configuration rules
	var cfgIgnores []string
	if configFile != "" {
		cfg, err := config.Load(configFile)
		if err != nil {
			fmt.Fprintf(c.stderr, "jdiff: %v\n", err)
			return ExitCodeError
		}
		cfgIgnores = cfg.Ignore
	} else {
		cfg, found, err := config.LoadDefault()
		if err != nil {
			fmt.Fprintf(c.stderr, "jdiff: %v\n", err)
			return ExitCodeError
		}
		if found && cfg != nil {
			cfgIgnores = cfg.Ignore
		}
	}

	activeRules := config.Merge(cliIgnores, cfgIgnores)

	// Validate rule syntax upfront
	pathMatcher, err := matcher.New(activeRules)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: invalid ignore rule: %v\n", err)
		return ExitCodeError
	}

	// Handle --show-config
	if showConfig {
		if len(activeRules) == 0 {
			fmt.Fprintln(c.stdout, "No ignore rules configured.")
		} else {
			fmt.Fprintln(c.stdout, "Ignore rules:")
			for _, r := range activeRules {
				fmt.Fprintf(c.stdout, "  %s\n", r)
			}
		}
		return ExitCodeOK
	}

	// Positional arguments validation
	if len(positionals) == 0 {
		fmt.Fprint(c.stderr, usage)
		return ExitCodeError
	}

	if len(positionals) < 2 {
		fmt.Fprintln(c.stderr, "jdiff: missing input files. Please provide both <old.json> and <new.json>")
		return ExitCodeError
	}

	if len(positionals) > 2 {
		fmt.Fprintln(c.stderr, "jdiff: too many arguments provided. Expected 2 JSON files.")
		return ExitCodeError
	}

	oldPath := positionals[0]
	newPath := positionals[1]

	oldData, err := os.ReadFile(oldPath)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: failed to read %s: %v\n", oldPath, err)
		return ExitCodeError
	}

	newData, err := os.ReadFile(newPath)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: failed to read %s: %v\n", newPath, err)
		return ExitCodeError
	}

	diffResult, err := diff.CompareBytesWithOptions(oldData, newData, diff.Options{
		Matcher: pathMatcher,
	})
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: %v\n", err)
		return ExitCodeError
	}

	// Color detection logic
	colorEnabled := true
	if noColor || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		colorEnabled = false
	}

	renderOpts := render.Options{
		Color:        colorEnabled,
		Compact:      compact,
		Verbose:      verbose,
		SummaryOnly:  summary,
		OldPath:      oldPath,
		NewPath:      newPath,
		IgnoredRules: activeRules,
	}

	if err := render.Render(c.stdout, diffResult, renderOpts); err != nil {
		fmt.Fprintf(c.stderr, "jdiff: render error: %v\n", err)
		return ExitCodeError
	}

	return ExitCodeOK
}
