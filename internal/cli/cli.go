package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"jdiff/internal/config"
	"jdiff/internal/diff"
	"jdiff/internal/matcher"
	"jdiff/internal/patch"
	"jdiff/internal/render"
	"jdiff/internal/stats"
	"jdiff/internal/version"
)

const (
	// ExitCodeOK indicates successful execution with no differences detected.
	ExitCodeOK = 0
	// ExitCodeDiff indicates differences were detected.
	ExitCodeDiff = 1
	// ExitCodeError indicates an operational error, invalid argument, or resource limit violation.
	ExitCodeError = 2
)

const usage = `jdiff - JSON Structural Diff

Usage:
  jdiff [options] <old.json> <new.json>
  jdiff apply [options] <patch.json> <input.json>

Options:
  --help                 Show help
  --version              Show version
  --output <format>      Output format: human, json, unified, patch (default: human)
  --output-file <file>   Write output to a file instead of stdout
  --verify-patch         Generate, apply, and verify the patch
  --stats                Display performance and memory statistics
  --max-file-size <size> Maximum allowed input file size (e.g. 100MB, 10KB, 500B)
  --max-changes <N>      Maximum number of differences to collect before truncating
  --exit-on-diff         Terminate comparison immediately upon discovering differences
  --quiet, -q            Suppress output and communicate exclusively via exit codes
  --no-color             Disable colored output
  --compact              Display compact diff output
  --verbose              Display additional comparison information
  --summary              Display only the change summary
  --ignore <path>        Ignore a JSON path (can be specified multiple times)
  --config <file>        Use a configuration file (defaults to .jdiff.json)
  --show-config          Show active comparison configuration

Arguments:
  <old.json>             Old JSON document (or - for stdin)
  <new.json>             New JSON document (or - for stdin)
`

// CLI manages command-line interface execution and I/O streams.
type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New creates a new CLI instance with given standard output and error writers and defaults stdin to os.Stdin.
func New(stdout, stderr io.Writer) *CLI {
	return NewWithStdin(os.Stdin, stdout, stderr)
}

// NewWithStdin creates a CLI instance with custom stdin, stdout, and stderr.
func NewWithStdin(stdin io.Reader, stdout, stderr io.Writer) *CLI {
	return &CLI{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// Run parses arguments and executes the requested operation.
func (c *CLI) Run(args []string) int {
	totalStart := time.Now()

	if len(args) > 0 && args[0] == "apply" {
		return c.runApply(args[1:])
	}

	var (
		outputFormatStr string
		outputFile      string
		verifyPatch     bool
		showStats       bool
		maxFileSizeStr  string
		maxChanges      int
		exitOnDiff      bool
		quiet           bool
		noColor         bool
		compact         bool
		verbose         bool
		summary         bool
		showConfig      bool
		configFile      string
		cliIgnores      []string
		positionals     []string
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
		case arg == "--output":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --output")
				return ExitCodeError
			}
			i++
			outputFormatStr = args[i]
		case strings.HasPrefix(arg, "--output="):
			outputFormatStr = strings.TrimPrefix(arg, "--output=")
		case arg == "--output-file":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --output-file")
				return ExitCodeError
			}
			i++
			outputFile = args[i]
		case strings.HasPrefix(arg, "--output-file="):
			outputFile = strings.TrimPrefix(arg, "--output-file=")
		case arg == "--verify-patch":
			verifyPatch = true
		case arg == "--stats":
			showStats = true
		case arg == "--max-file-size":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --max-file-size")
				return ExitCodeError
			}
			i++
			maxFileSizeStr = args[i]
		case strings.HasPrefix(arg, "--max-file-size="):
			maxFileSizeStr = strings.TrimPrefix(arg, "--max-file-size=")
		case arg == "--max-changes":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --max-changes")
				return ExitCodeError
			}
			i++
			val, err := strconv.Atoi(args[i])
			if err != nil || val <= 0 {
				fmt.Fprintf(c.stderr, "jdiff: invalid value for --max-changes: %q\n", args[i])
				return ExitCodeError
			}
			maxChanges = val
		case strings.HasPrefix(arg, "--max-changes="):
			val, err := strconv.Atoi(strings.TrimPrefix(arg, "--max-changes="))
			if err != nil || val <= 0 {
				fmt.Fprintf(c.stderr, "jdiff: invalid value for --max-changes: %q\n", arg)
				return ExitCodeError
			}
			maxChanges = val
		case arg == "--exit-on-diff":
			exitOnDiff = true
		case arg == "--quiet" || arg == "-q":
			quiet = true
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
		case strings.HasPrefix(arg, "--config="):
			configFile = strings.TrimPrefix(arg, "--config=")
		case arg == "--ignore":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff: missing argument for --ignore")
				return ExitCodeError
			}
			i++
			cliIgnores = append(cliIgnores, args[i])
		case strings.HasPrefix(arg, "--ignore="):
			cliIgnores = append(cliIgnores, strings.TrimPrefix(arg, "--ignore="))
		case strings.HasPrefix(arg, "-") && arg != "-":
			fmt.Fprintf(c.stderr, "jdiff: unknown flag %q\n", arg)
			return ExitCodeError
		default:
			positionals = append(positionals, arg)
		}
	}

	// Validate output format
	format, err := render.ParseFormat(outputFormatStr)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: %v\n", err)
		return ExitCodeError
	}

	// Parse max-file-size limit if specified
	var maxFileSizeBytes int64
	if maxFileSizeStr != "" {
		parsedLimit, err := stats.ParseSize(maxFileSizeStr)
		if err != nil {
			fmt.Fprintf(c.stderr, "jdiff: invalid --max-file-size: %v\n", err)
			return ExitCodeError
		}
		maxFileSizeBytes = parsedLimit
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

	if oldPath == "-" && newPath == "-" {
		fmt.Fprintln(c.stderr, "jdiff: cannot read both inputs from stdin")
		return ExitCodeError
	}

	// Read input documents respecting max-file-size
	oldData, err := c.readInputWithLimit(oldPath, maxFileSizeBytes)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: %s: %v\n", oldPath, err)
		return ExitCodeError
	}

	newData, err := c.readInputWithLimit(newPath, maxFileSizeBytes)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: %s: %v\n", newPath, err)
		return ExitCodeError
	}

	var mStart runtime.MemStats
	if showStats {
		runtime.ReadMemStats(&mStart)
	}

	parseStart := time.Now()
	var oldVal, newVal any
	decOld := json.NewDecoder(bytes.NewReader(oldData))
	decOld.UseNumber()
	if err := decOld.Decode(&oldVal); err != nil {
		fmt.Fprintf(c.stderr, "jdiff: failed to parse %s: %v\n", oldPath, err)
		return ExitCodeError
	}

	decNew := json.NewDecoder(bytes.NewReader(newData))
	decNew.UseNumber()
	if err := decNew.Decode(&newVal); err != nil {
		fmt.Fprintf(c.stderr, "jdiff: failed to parse %s: %v\n", newPath, err)
		return ExitCodeError
	}
	parseDuration := time.Since(parseStart)

	compareStart := time.Now()
	diffResult, err := diff.CompareBytesWithOptions(oldData, newData, diff.Options{
		Matcher:    pathMatcher,
		MaxChanges: maxChanges,
		EarlyExit:  exitOnDiff,
	})
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff: %v\n", err)
		return ExitCodeError
	}
	compareDuration := time.Since(compareStart)
	totalDuration := time.Since(totalStart)

	var statInfo *stats.Stats
	if showStats {
		var mEnd runtime.MemStats
		runtime.ReadMemStats(&mEnd)

		allocDiff := mEnd.TotalAlloc - mStart.TotalAlloc
		statInfo = &stats.Stats{
			OldSize:      int64(len(oldData)),
			OldIsStdin:   oldPath == "-",
			NewSize:      int64(len(newData)),
			NewIsStdin:   newPath == "-",
			ParseTime:    parseDuration,
			CompareTime:  compareDuration,
			TotalTime:    totalDuration,
			AllocBytes:   allocDiff,
			ChangesCount: len(diffResult.Changes),
		}
	}

	// Handle --verify-patch
	if verifyPatch {
		patchDoc := patch.Generate(diffResult)
		ok, err := patch.Verify(oldData, newData, patchDoc)
		if err != nil || !ok {
			if !quiet {
				fmt.Fprintln(c.stdout, "Patch verification failed.")
			}
			return ExitCodeDiff
		}
		if !quiet {
			fmt.Fprintln(c.stdout, "Patch verification successful.")
		}
		return ExitCodeOK
	}

	// Handle --quiet mode
	if quiet {
		if diffResult.HasChanges() {
			return ExitCodeDiff
		}
		return ExitCodeOK
	}

	// Color detection logic
	colorEnabled := true
	if noColor || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		colorEnabled = false
	}

	renderOpts := render.Options{
		Format:       format,
		Color:        colorEnabled,
		Compact:      compact,
		Verbose:      verbose,
		SummaryOnly:  summary,
		OldPath:      oldPath,
		NewPath:      newPath,
		IgnoredRules: activeRules,
		Stats:        statInfo,
	}

	// Handle statistics routing when patch format is selected
	if format == render.FormatPatch && showStats && statInfo != nil {
		render.Render(c.stderr, diffResult, render.Options{
			Format: render.FormatHuman,
			Stats:  statInfo,
		})
	}

	// Determine output destination writer
	outWriter := c.stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(c.stderr, "jdiff: failed to create output file %s: %v\n", outputFile, err)
			return ExitCodeError
		}
		defer f.Close()
		outWriter = f
	}

	if err := render.Render(outWriter, diffResult, renderOpts); err != nil {
		fmt.Fprintf(c.stderr, "jdiff: render error: %v\n", err)
		return ExitCodeError
	}

	if diffResult.HasChanges() {
		return ExitCodeDiff
	}

	return ExitCodeOK
}

func (c *CLI) runApply(args []string) int {
	var (
		outputFile  string
		positionals []string
	)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--output-file":
			if i+1 >= len(args) {
				fmt.Fprintln(c.stderr, "jdiff apply: missing argument for --output-file")
				return ExitCodeError
			}
			i++
			outputFile = args[i]
		case strings.HasPrefix(arg, "--output-file="):
			outputFile = strings.TrimPrefix(arg, "--output-file=")
		case arg == "--help" || arg == "-h":
			fmt.Fprintln(c.stdout, "Usage: jdiff apply [options] <patch.json> <input.json>")
			return ExitCodeOK
		case strings.HasPrefix(arg, "-") && arg != "-":
			fmt.Fprintf(c.stderr, "jdiff apply: unknown flag %q\n", arg)
			return ExitCodeError
		default:
			positionals = append(positionals, arg)
		}
	}

	if len(positionals) < 2 {
		fmt.Fprintln(c.stderr, "jdiff apply: missing required arguments: <patch.json> <input.json>")
		return ExitCodeError
	}
	if len(positionals) > 2 {
		fmt.Fprintln(c.stderr, "jdiff apply: too many arguments provided")
		return ExitCodeError
	}

	patchPath := positionals[0]
	inputPath := positionals[1]

	if patchPath == "-" && inputPath == "-" {
		fmt.Fprintln(c.stderr, "jdiff apply: cannot read both patch and input from stdin")
		return ExitCodeError
	}

	patchData, err := c.readInputWithLimit(patchPath, 0)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff apply: failed to read patch %s: %v\n", patchPath, err)
		return ExitCodeError
	}

	inputData, err := c.readInputWithLimit(inputPath, 0)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff apply: failed to read input %s: %v\n", inputPath, err)
		return ExitCodeError
	}

	var patchDoc patch.Patch
	if err := json.Unmarshal(patchData, &patchDoc); err != nil {
		fmt.Fprintf(c.stderr, "jdiff apply: invalid JSON Patch document: %v\n", err)
		return ExitCodeError
	}

	var inputDoc any
	dec := json.NewDecoder(bytes.NewReader(inputData))
	dec.UseNumber()
	if err := dec.Decode(&inputDoc); err != nil {
		fmt.Fprintf(c.stderr, "jdiff apply: invalid input JSON document: %v\n", err)
		return ExitCodeError
	}

	resultDoc, err := patch.Apply(inputDoc, patchDoc)
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff apply: %v\n", err)
		return ExitCodeError
	}

	outBytes, err := json.MarshalIndent(resultDoc, "", "  ")
	if err != nil {
		fmt.Fprintf(c.stderr, "jdiff apply: failed to format output JSON: %v\n", err)
		return ExitCodeError
	}

	outWriter := c.stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(c.stderr, "jdiff apply: failed to create output file %s: %v\n", outputFile, err)
			return ExitCodeError
		}
		defer f.Close()
		outWriter = f
	}

	fmt.Fprintln(outWriter, string(outBytes))
	return ExitCodeOK
}

func (c *CLI) readInputWithLimit(path string, maxBytes int64) ([]byte, error) {
	var r io.Reader
	if path == "-" {
		if c.stdin == nil {
			return nil, fmt.Errorf("stdin not available")
		}
		r = c.stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if maxBytes > 0 {
			info, err := f.Stat()
			if err == nil && info.Size() > maxBytes {
				return nil, fmt.Errorf("input exceeds maximum allowed file size (%s)", stats.FormatBytes(maxBytes))
			}
		}
		r = f
	}

	if maxBytes > 0 {
		lr := io.LimitReader(r, maxBytes+1)
		data, err := io.ReadAll(lr)
		if err != nil {
			return nil, err
		}
		if int64(len(data)) > maxBytes {
			return nil, fmt.Errorf("input exceeds maximum allowed file size (%s)", stats.FormatBytes(maxBytes))
		}
		return data, nil
	}

	return io.ReadAll(r)
}
