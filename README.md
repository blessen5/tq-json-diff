# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values with precise structural path tracking, intelligent array element comparison, selective ignore rules, and multiple output formats.

**Current Version:** `v0.7.0`

## Features

- **Multiple Output Formats**:
  - **`human` (default)**: Colorized, human-readable terminal output.
  - **`json`**: Strict, machine-readable JSON output compliant with automated toolchains.
  - **`unified`**: Unified diff representation (`@@ <path>`).
- **Standard Input (`-`) & Shell Piping**: Compare JSON streams directly from standard input (e.g. `cat old.json | jdiff - new.json`).
- **File Output (`--output-file`)**: Redirect clean diff results directly to a target file.
- **Clean Stream Separation**: Diff output is written to `stdout` (or file); operational errors are routed exclusively to `stderr`.
- **Selective JSON Comparison & Ignore Rules**: Ignore dynamic fields (timestamps, IDs) via CLI `--ignore` flags or `.jdiff.json` configuration files with wildcard support (`*`, `[*]`).
- **Granular Array Comparison**: Deterministic index-based array element diffing (`users[0].name`, `languages[1]`).
- **Deep Recursive Comparison**: Traverses arbitrary JSON nesting depths without flagging parent structures.
- **Precise JSON Path Engine**: Dot-separated property keys and bracketed array indices.
- **Zero External Dependencies**: Implemented entirely with the Go standard library.

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
| `--version`, `-v` | Display application version (`jdiff v0.7.0`) |
| `--output <format>` | Select output format: `human` (default), `json`, `unified` |
| `--output-file <file>` | Write diff output directly to a file |
| `--ignore <path>` | Ignore a JSON path or pattern (repeatable) |
| `--config <file>` | Use a configuration file (defaults to `.jdiff.json`) |
| `--show-config` | Display active ignore configuration and exit |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Display active ignore rules and comparison file context |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## Output Formats & Automation

### 1. Machine-Readable JSON (`--output json`)

```bash
jdiff --output json examples/output-old.json examples/output-new.json
```

**Output**:
```json
{
  "summary": {
    "added": 1,
    "ignored": 0,
    "modified": 3,
    "removed": 0,
    "total": 4
  },
  "changes": [
    {
      "new": true,
      "old": false,
      "path": "config.debug",
      "type": "modified"
    },
    {
      "new": 60,
      "old": 30,
      "path": "config.timeout",
      "type": "modified"
    },
    {
      "new": "/oauth",
      "path": "endpoints[2]",
      "type": "added"
    },
    {
      "new": 2,
      "old": 1,
      "path": "version",
      "type": "modified"
    }
  ]
}
```

### 2. Unified Diff Format (`--output unified`)

```bash
jdiff --output unified examples/output-old.json examples/output-new.json
```

**Output**:
```text
--- examples/output-old.json
+++ examples/output-new.json
@@ config.debug
- false
+ true

@@ config.timeout
- 30
+ 60

@@ endpoints[2]
+ "/oauth"

@@ version
- 1
+ 2
```

### 3. Piping and Standard Input (`-`)

```bash
# Pipe old document via stdin
cat old.json | jdiff - new.json

# Pipe dynamic API output and write JSON diff to file
curl -s https://api.example.com/v2/config | jdiff --output json --output-file diff.json production.json -
```

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff executed or help/version/config printed) |
| `1` | Operational error (invalid arguments, unreadable files, invalid JSON, unsupported format) |

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [x] **v0.4.0**: Professional terminal presentation (ANSI colors, `--no-color`, `--compact`, `--verbose`, `--summary`).
- [x] **v0.5.0**: Granular index-based array comparison, array paths (`users[0].name`), nested arrays, and arrays of objects.
- [x] **v0.6.0**: Ignore rules, wildcard paths (`*.key`, `users[*].id`), configuration files (`.jdiff.json`), and `--show-config`.
- [x] **v0.7.0**: Multiple output formats (`human`, `json`, `unified`), `--output-file`, and standard input (`-`) piping.
- [ ] **v0.8.0**: Standardized JSON Patch (RFC 6902) export.

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

## License

This project is licensed under the [MIT License](LICENSE).
