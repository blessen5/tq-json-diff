# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values with precise structural path tracking, intelligent array element comparison, and rich terminal presentation modes.

**Current Version:** `v0.5.0`

## Features

- **Granular Array Comparison**: Performs deterministic index-based comparison of arrays, detecting added elements, removed elements, and modified values at specific indexes (e.g. `languages[1]`, `users[0].name`).
- **Deep Recursive Comparison**: Recursively traverses nested JSON structures and nested arrays to arbitrary depths without falsely reporting parent containers.
- **Precise JSON Path Engine**: Tracks dot-separated keys and bracketed array indexes (`data.groups[0].values[1]`, `[0]` for root arrays).
- **Arrays of Objects**: Diffing objects nested within array items down to the property level.
- **Explicit Type Change Detection**: Identifies type shifts between strings, numbers, booleans, null, objects, and arrays.
- **Root JSON Value Support**: Compares root primitives (`"hello"`, `10`, `true`, `null`), root arrays, and root objects safely.
- **Semantic ANSI Color Highlighting**: Added values in green (`+`), removed values in red (`-`), modified values in yellow, and path identifiers in cyan.
- **Multiple Presentation Modes**:
  - **Standard**: Hierarchical diff with change summaries.
  - **Compact (`--compact`)**: Streamlined single-line representations per modification (`MODIFIED items[1]: "B" → "C"`).
  - **Summary-Only (`--summary`)**: Displays only high-level change counters and totals.
  - **Verbose (`--verbose`)**: Includes file paths and comparison context.
- **CI/CD & Scripting Support (`--no-color`)**: Automatic color detection with explicit `--no-color` override and standard `NO_COLOR` environment variable support.
- **Zero External Dependencies**: Implemented entirely using the Go standard library.

## Supported JSON Structures (v0.5.0)

- **Objects**: Flat, nested, and deeply nested objects.
- **Arrays**: Index-based element comparison, arrays of primitives, arrays of objects, and nested arrays.
- **Strings**: Quoted string comparisons.
- **Numbers**: Floating-point and integer numbers with full precision preservation (`json.Number`).
- **Booleans**: `true` and `false` literals.
- **Null Values**: `null` literal support.
- **Root Primitives**: Single primitive, array, or null root documents.

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
| `--version`, `-v` | Display application version (`jdiff v0.5.0`) |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Include comparison file context before diff output |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## Array Comparison Example

Comparing [`examples/arrays-old.json`](examples/arrays-old.json) and [`examples/arrays-new.json`](examples/arrays-new.json):

```bash
jdiff examples/arrays-old.json examples/arrays-new.json
```

**Output**:
```text
MODIFIED
  languages[1]
    - "Python"
    + "Rust"

  maintainers[0].role
    - "Lead"
    + "Creator"

ADDED
  languages[3]
    + "TypeScript"

  maintainers[2]
    + {"name":"Bob","role":"Reviewer"}

REMOVED
  features[2]
    - "ansi-colors"

Summary:
  Added:     2
  Removed:   1
  Modified:  2
```

---

## Current Array Limitations

> [!NOTE]
> `jdiff v0.5.0` uses **deterministic index-based array comparison** (`old[i]` vs `new[i]`).
>
> It does not currently perform:
> - ID/key-based entity matching across different array indexes
> - Longest-common-subsequence (LCS) alignment or fuzzy similarity matching
> - Element move and reorder detection
>
> If an element is inserted at the beginning of an array, subsequent items at shifted indexes are evaluated by their corresponding index position. Advanced key-based alignment and heuristic reorder detection are planned for upcoming releases.

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff executed or help/version printed) |
| `1` | Operational error (missing arguments, unreadable files, invalid JSON syntax) |

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [x] **v0.4.0**: Professional terminal presentation (ANSI colors, `--no-color`, `--compact`, `--verbose`, `--summary`).
- [x] **v0.5.0**: Granular index-based array comparison, array paths (`users[0].name`), nested arrays, and arrays of objects.
- [ ] **v0.6.0**: Standardized JSON Patch (RFC 6902) export.
- [ ] **v0.7.0**: Identity/key-based array alignment and reorder detection.

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
