# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values.

**Current Version:** `v0.2.0`

## Features

- **Structural Diffing**: Compares JSON documents by key and value structure rather than raw line-by-line text differences.
- **Dot-Separated Path Tracking**: Pinpoints exact modification locations (e.g. `user.profile.email`).
- **Unchanged Value Suppression**: Unchanged values are hidden by default to keep diff outputs concise.
- **Zero External Dependencies**: Built entirely with Go standard library components (`encoding/json`, `reflect`, etc.).
- **Deterministic Output**: Keys are sorted alphabetically to ensure reproducible diff representations across runs.
- **Cross-Platform**: Compiles into a single standalone binary for Linux, macOS, and Windows.

## Supported JSON Structures (v0.2.0)

- **Objects**: Flat and deeply nested objects.
- **Strings**: Quoted string comparisons.
- **Numbers**: Floating-point and integer numbers with full precision preservation (`json.Number`).
- **Booleans**: `true` and `false` literals.
- **Null Values**: `null` literal support.
- **Arrays**: Deterministic whole-value equality comparison (see [Limitations](#current-limitations)).

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
jdiff <old.json> <new.json>
```

### Options & Flags

```bash
jdiff --help       # Show help and usage instructions
jdiff --version    # Show application version (v0.2.0)
```

### Example

Comparing `old.json` and `new.json`:

**`old.json`**:
```json
{
  "name": "Blessen",
  "age": 19,
  "city": "Kochi"
}
```

**`new.json`**:
```json
{
  "name": "Blessen",
  "age": 20,
  "city": "Bengaluru",
  "country": "India"
}
```

**Command**:
```bash
jdiff old.json new.json
```

**Output**:
```text
MODIFIED:
    age
        - 19
        + 20

    city
        - "Kochi"
        + "Bengaluru"

ADDED:
    country
        + "India"
```

*(Note: `name` is omitted from the output because it did not change.)*

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff executed or help/version printed) |
| `1` | Operational error (missing arguments, unreadable files, invalid JSON syntax) |

## Current Limitations

- **Array Comparison**: In version `v0.2.0`, arrays are treated as whole values. If an array's contents differ, the entire array is reported as modified. Advanced array element matching, LCS-based diffing, insertions, deletions, and item reordering detection are planned for future versions.
- **Format Options**: Output is currently formatted in human-readable terminal text. JSON Patch (RFC 6902) and colorized themes are planned for upcoming releases.

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, recursive object traversal, path tracking, primitive comparisons, deterministic output.
- [ ] **v0.3.0**: Advanced array diffing (element alignment, insertions, deletions, reordering).
- [ ] **v0.4.0**: Terminal color highlights and customized formatting options.
- [ ] **v0.5.0**: Standardized JSON Patch (RFC 6902) generation and summary statistics.

## Development & Testing

```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Run unit and integration test suite
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
│   │   ├── cli.go               # Command-line interface and argument parsing
│   │   └── cli_test.go          # CLI unit tests
│   ├── diff/
│   │   ├── diff.go              # Recursive JSON structural diff engine
│   │   └── diff_test.go         # Diff engine unit tests
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
