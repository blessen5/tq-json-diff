# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values with precise structural path tracking and rich presentation modes.

**Current Version:** `v0.4.0`

## Features

- **Semantic ANSI Color Highlighting**: Added values in green (`+`), removed values in red (`-`), modified values in yellow, and path identifiers in cyan.
- **Multiple Presentation Modes**:
  - **Standard**: Human-readable hierarchical diff with change summaries.
  - **Compact (`--compact`)**: Streamlined single-line representations per modification (`MODIFIED user.age: 19 → 20`).
  - **Summary-Only (`--summary`)**: Displays only high-level change counters and totals.
  - **Verbose (`--verbose`)**: Includes file paths and comparison context.
- **Color Safety & CI/CD Support (`--no-color`)**: Automatic color detection with explicit `--no-color` override and standard `NO_COLOR` environment variable support.
- **Deep Recursive Comparison**: Traverses nested JSON structures to any depth without falsely reporting parent nodes.
- **Precise JSON Path Engine**: Tracks dot-separated paths (e.g. `user.profile.contact.email`).
- **Explicit Type Change Detection**: Identifies type shifts between strings, numbers, booleans, null, objects, and arrays.
- **Root JSON Value Support**: Compares root primitives (`"hello"`, `10`, `true`, `null`), arrays, and objects safely.
- **Deterministic Traversal & Ordering**: Predictable alphabetical key and path ordering across all executions.
- **Zero External Dependencies**: Built entirely with the Go standard library.

## Installation

### From Source

```bash
go install ./cmd/jdiff
```

Or build locally:

```bash
go build -o jdiff .
```

## Usage

```bash
jdiff [options] <old.json> <new.json>
```

### Options & Flags

| Flag | Description |
|---|---|
| `--help`, `-h` | Display usage and available options |
| `--version`, `-v` | Display application version (`jdiff v0.4.0`) |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Include comparison file context before diff output |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## Output Examples

### 1. Standard Mode (Default)

```bash
jdiff examples/basic-old.json examples/basic-new.json
```

```text
MODIFIED
  age
    - 19
    + 20

  city
    - "Kochi"
    + "Bengaluru"

ADDED
  country
    + "India"

Summary:
  Added:     1
  Removed:   0
  Modified:  2
```

### 2. Compact Mode (`--compact`)

```bash
jdiff --compact examples/basic-old.json examples/basic-new.json
```

```text
MODIFIED age: 19 → 20
MODIFIED city: "Kochi" → "Bengaluru"
ADDED country: "India"

Summary:
  Added:     1
  Removed:   0
  Modified:  2
```

### 3. Summary-Only Mode (`--summary`)

```bash
jdiff --summary examples/basic-old.json examples/basic-new.json
```

```text
JSON Diff Summary

Added:     1
Removed:   0
Modified:  2
Total:     3
```

### 4. CI/CD & Scripting (`--no-color`)

```bash
jdiff --no-color --compact old.json new.json
```

Ensures clean, plain-text output without ANSI escape characters for log files and automated diff pipelines.

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff executed or help/version printed) |
| `1` | Operational error (missing arguments, unreadable files, invalid JSON syntax) |

## Current Limitations

- **Array Comparison**: In version `v0.4.0`, arrays are evaluated using deterministic whole-value equality. If elements within an array differ, the entire array is flagged as modified. Fine-grained element tracking, longest-common-subsequence (LCS) alignment, item insertions/deletions, and item reordering detection are planned for v0.5.0.
- **Machine-Readable Formats**: Output currently focuses on human-readable terminal presentations. JSON Patch (RFC 6902) export is scheduled for upcoming releases.

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [x] **v0.4.0**: Professional terminal presentation (ANSI colors, `--no-color`, `--compact`, `--verbose`, `--summary`).
- [ ] **v0.5.0**: Advanced array diffing (element alignment, insertions, deletions, reordering).
- [ ] **v0.6.0**: Standardized JSON Patch (RFC 6902) export.

## Development & Testing

```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Run test suite
go test -v ./...

# Build binary
go build ./...
```

## Project Structure

```text
jdiff/
├── cmd/
│   └── jdiff/
│       └── main.go              # Canonical command entry point
├── internal/
│   ├── cli/
│   │   ├── cli.go               # Command-line interface and flag handling
│   │   └── cli_test.go          # CLI unit tests
│   ├── diff/
│   │   ├── diff.go              # Deep structural diff engine
│   │   ├── diff_test.go         # Diff engine unit tests
│   │   └── path.go              # Structured JSON Path representation
│   ├── render/
│   │   ├── render.go            # Presentation layer (ANSI colors, compact, verbose, summary)
│   │   └── render_test.go       # Render unit tests
│   └── version/
│       ├── version.go           # Version and build metadata
│       └── version_test.go      # Version tests
├── docs/                        # Architecture and reference documentation
├── examples/                    # Sample JSON datasets
├── tests/                       # Integration test suite
├── go.mod                       # Go module definition
├── main.go                      # Root application entry point
├── README.md                    # Project documentation
├── CONTRIBUTING.md              # Contribution guidelines
├── LICENSE                      # License terms
└── .gitignore                   # Ignored files
```

## License

This project is licensed under the [MIT License](LICENSE).
