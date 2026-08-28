# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values.

## Features

- **Structural Diffing**: Compares JSON documents by key and value structure rather than raw line-by-line text differences.
- **Zero External Dependencies**: Built entirely with Go standard library components.
- **Clear Terminal Output**: High-visibility reporting of additions, deletions, modifications, and type shifts.
- **Cross-Platform**: Compiles into a single standalone binary for Linux, macOS, and Windows.

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
jdiff --version    # Show application version information
```

### Example Help Output

```text
jdiff - JSON Structural Diff

Usage:
  jdiff <old.json> <new.json>

Commands:
  --help       Show help
  --version    Show version
```

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success / No differences found / Informational output (`--help`, `--version`) |
| `1` | Differences detected between the two JSON files |
| `2` | Error occurred (invalid arguments, file unreadable, JSON syntax error) |

## Development

Verify formatting, linting, and tests:

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
│   └── jdiff/              # Application CLI entry point
├── internal/
│   ├── cli/                # Command-line interface and argument parsing
│   └── version/            # Version and build metadata
├── docs/                   # Architecture and technical documentation
├── examples/               # Sample JSON datasets
├── tests/                  # Integration test suite
├── go.mod                  # Go module definition
├── main.go                 # Root application entry point
├── README.md               # Project documentation
├── CONTRIBUTING.md         # Contribution guidelines
├── LICENSE                 # License terms
└── .gitignore              # Ignored files
```

## License

This project is licensed under the [MIT License](LICENSE).
