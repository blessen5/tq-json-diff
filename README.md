# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values with precise structural path tracking, intelligent array element comparison, selective ignore rules, and rich terminal presentation modes.

**Current Version:** `v0.6.0`

## Features

- **Selective JSON Comparison & Ignore Rules**: Ignore dynamic or generated fields (e.g. timestamps, session IDs, request tokens) via CLI `--ignore` flags or configuration files.
- **Flexible Path Pattern Matching**:
  - **Exact Path**: `timestamp`, `metadata.created_at`, `users[0].session_id`
  - **Object Subtree Pruning**: Ignoring `metadata` excludes the entire `metadata` sub-hierarchy.
  - **Wildcard Keys (`*`)**: `*.timestamp`, `users.*.email`
  - **Array Wildcards (`[*]`)**: `users[*].id`, `data.groups[*].values[*]`
- **Configuration File Support (`.jdiff.json`)**: Automatically loads `.jdiff.json` if present in the current working directory, or an explicit file via `--config <path>`.
- **Config Inspection (`--show-config`)**: Inspect merged active ignore rules without running a diff.
- **Granular Array Comparison**: Performs deterministic index-based comparison of arrays, detecting added, removed, and modified elements.
- **Deep Recursive Comparison**: Recursively traverses nested objects and arrays without falsely reporting parent containers.
- **Precise JSON Path Engine**: Dot-separated property keys and bracketed array indices (`data.groups[0].values[1]`).
- **Semantic ANSI Color Highlighting**: Added values in green (`+`), removed values in red (`-`), modified values in yellow, and path identifiers in cyan.
- **Multiple Presentation Modes**:
  - **Standard**: Hierarchical diff with change and ignored summaries.
  - **Compact (`--compact`)**: Streamlined single-line representations per modification.
  - **Summary-Only (`--summary`)**: Displays only change counters, ignored metrics, and totals.
  - **Verbose (`--verbose`)**: Lists ignored patterns and file path context before the diff.
- **CI/CD & Scripting Support (`--no-color`)**: Automatic color detection with explicit `--no-color` override and standard `NO_COLOR` environment variable support.
- **Zero External Dependencies**: Built strictly using the Go standard library.

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
| `--version`, `-v` | Display application version (`jdiff v0.6.0`) |
| `--ignore <path>` | Ignore a JSON path or pattern (repeatable) |
| `--config <file>` | Use a configuration file (defaults to `.jdiff.json`) |
| `--show-config` | Display active ignore configuration and exit |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Display active ignore rules and comparison file context |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## Ignore Rules & Configuration

### 1. Command-Line Ignore Rules

You can supply one or more `--ignore` flags:

```bash
jdiff --ignore timestamp --ignore "*.updated_at" --ignore "users[*].session_id" old.json new.json
```

### 2. Configuration File (`.jdiff.json`)

Create a `.jdiff.json` file in your repository or working directory:

```json
{
  "ignore": [
    "timestamp",
    "request_id",
    "*.updated_at",
    "users[*].session_id"
  ]
}
```

Running `jdiff old.json new.json` will automatically load `.jdiff.json` if present.

To specify a custom configuration file:

```bash
jdiff --config ./custom-rules.json old.json new.json
```

### 3. Rule Precedence & Merging

When both CLI `--ignore` flags and a configuration file are present:
1. Rules are merged into a unified set.
2. CLI rules take precedence in order of evaluation.
3. Duplicate rules are automatically deduplicated.

### 4. Viewing Active Rules

```bash
jdiff --config .jdiff.json --show-config
```

**Output**:
```text
Ignore rules:
  timestamp
  request_id
  *.updated_at
  users[*].session_id
```

---

## Selective Comparison Example

Comparing [`examples/ignore-old.json`](examples/ignore-old.json) and [`examples/ignore-new.json`](examples/ignore-new.json) with [`examples/.jdiff.json`](examples/.jdiff.json):

```bash
jdiff --config examples/.jdiff.json examples/ignore-old.json examples/ignore-new.json
```

**Output**:
```text
MODIFIED
  users[0].name
    - "Alice"
    + "Alice Cooper"

  version
    - "v0.5.0"
    + "v0.6.0"

Summary:
  Added:     0
  Removed:   0
  Modified:  2
  Ignored:   5
```

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff executed or help/version/config printed) |
| `1` | Operational error (invalid arguments, unreadable files, invalid JSON syntax, invalid config) |

## Current Limitations

- **Array Comparison**: Uses index-based matching (`old[i]` vs `new[i]`). Advanced ID/key-based entity matching across differing positions and heuristic reorder detection are planned for future releases.
- **Wildcard Syntax**: Supports exact paths, object segment wildcard (`*`), and array index wildcard (`[*]`). Arbitrary regex or full globs are not supported.

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [x] **v0.4.0**: Professional terminal presentation (ANSI colors, `--no-color`, `--compact`, `--verbose`, `--summary`).
- [x] **v0.5.0**: Granular index-based array comparison, array paths (`users[0].name`), nested arrays, and arrays of objects.
- [x] **v0.6.0**: Ignore rules, wildcard paths (`*.key`, `users[*].id`), configuration files (`.jdiff.json`), and `--show-config`.
- [ ] **v0.7.0**: Standardized JSON Patch (RFC 6902) export.

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
