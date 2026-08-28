# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents, clearly displays differences across multiple presentation formats, and generates RFC 6902-compliant JSON Patch documents with built-in patch application and verification.

**Current Version:** `v0.8.0`

## Features

- **JSON Patch Generation (RFC 6902 subset)**: Generates valid JSON Patch documents with `add`, `remove`, and `replace` operations via `--output patch`.
- **In-Memory Patch Application (`jdiff apply`)**: Apply a JSON Patch document to an input JSON file and output the resulting document.
- **Patch Verification (`--verify-patch`)**: Atomically diffs, generates a patch, applies it in memory, and verifies that the patched document matches the target.
- **RFC 6901 JSON Pointer Paths**: Correctly encodes and decodes JSON Pointer path tokens, including `~0` and `~1` special character escaping.
- **Safe Array Patch Ordering**: Orders multi-element array removals in descending index sequence to prevent index shifting errors during sequential patch application.
- **Multiple Output Formats**:
  - **`human` (default)**: Colorized, human-readable terminal output.
  - **`json`**: Strict, machine-readable JSON output with summary metrics and delta objects.
  - **`unified`**: Unified diff representation (`@@ <path>`).
  - **`patch`**: RFC 6902 JSON Patch document.
- **Standard Input (`-`) & Shell Piping**: Compare JSON streams directly from standard input.
- **Selective JSON Comparison & Ignore Rules**: Ignore dynamic fields via CLI `--ignore` flags or `.jdiff.json` configuration files with wildcard support (`*`, `[*]`).
- **Granular Array Comparison**: Deterministic index-based array element diffing (`users[0].name`, `languages[1]`).
- **Deep Recursive Comparison**: Recursively traverses nested objects and arrays without falsely reporting parent containers.
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
# Compare two JSON documents
jdiff [options] <old.json> <new.json>

# Apply a JSON Patch document
jdiff apply [options] <patch.json> <input.json>
```

### Options & Flags

| Flag | Description |
|---|---|
| `--help`, `-h` | Display usage and available options |
| `--version`, `-v` | Display application version (`jdiff v0.8.0`) |
| `--output <format>` | Select output format: `human` (default), `json`, `unified`, `patch` |
| `--output-file <file>` | Write diff/patch output directly to a file |
| `--verify-patch` | Generate, apply, and verify the patch in memory |
| `--ignore <path>` | Ignore a JSON path or pattern (repeatable) |
| `--config <file>` | Use a configuration file (defaults to `.jdiff.json`) |
| `--show-config` | Display active ignore configuration and exit |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Display active ignore rules and comparison file context |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## JSON Patch Workflow

### 1. Generating a JSON Patch (`--output patch`)

```bash
jdiff --output patch examples/output-old.json examples/output-new.json
```

**Output**:
```json
[
  {
    "op": "replace",
    "path": "/config/debug",
    "value": true
  },
  {
    "op": "replace",
    "path": "/config/timeout",
    "value": 60
  },
  {
    "op": "replace",
    "path": "/version",
    "value": 2
  },
  {
    "op": "add",
    "path": "/endpoints/2",
    "value": "/oauth"
  }
]
```

### 2. Verifying a Patch (`--verify-patch`)

```bash
jdiff --verify-patch examples/output-old.json examples/output-new.json
```

**Output**:
```text
Patch verification successful.
```

### 3. Applying a Patch (`jdiff apply`)

```bash
# Generate patch to file
jdiff --output patch --output-file update.patch.json old.json new.json

# Apply patch to old document
jdiff apply --output-file updated.json update.patch.json old.json
```

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff executed, patch verified, or help printed) |
| `1` | Operational error (invalid JSON, patch application failure, unsupported format) |

## Current Limitations

- **RFC 6902 Scope**: Supports `add`, `remove`, and `replace` operations. Operations `move`, `copy`, and `test` are not currently emitted.
- **Array Alignment**: Uses index-based array alignment (`old[i]` vs `new[i]`).

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [x] **v0.4.0**: Professional terminal presentation (ANSI colors, `--no-color`, `--compact`, `--verbose`, `--summary`).
- [x] **v0.5.0**: Granular index-based array comparison, array paths (`users[0].name`), nested arrays, and arrays of objects.
- [x] **v0.6.0**: Ignore rules, wildcard paths (`*.key`, `users[*].id`), configuration files (`.jdiff.json`), and `--show-config`.
- [x] **v0.7.0**: Multiple output formats (`human`, `json`, `unified`), `--output-file`, and standard input (`-`) piping.
- [x] **v0.8.0**: JSON Patch generation (RFC 6902), patch application (`jdiff apply`), and verification (`--verify-patch`).

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
