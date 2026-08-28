# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents and clearly shows added, removed, and modified values with precise structural path tracking.

**Current Version:** `v0.3.0`

## Features

- **Deep Recursive Comparison**: Recursively traverses nested JSON structures to any depth and pinpoints the exact changed leaf values.
- **Precise JSON Path Engine**: Tracks dot-separated paths (e.g. `user.profile.contact.email`, `(root)` for document root).
- **Structural Delta Classification**: Distinguishes added properties, removed properties, modified primitive values, and new/removed nested sub-objects.
- **Explicit Type Change Detection**: Identifies type shifts between strings, numbers, booleans, null, objects, and arrays (e.g. `20` -> `"20"`).
- **Root Value Support**: Safely compares root primitives (`"hello"`, `10`, `true`, `null`), root arrays, and root objects.
- **Deterministic Traversal & Ordering**: Guarantees identical, predictable diff output across runs by sorting paths alphabetically.
- **Change Summary Statistics**: Displays a concise summary of added, removed, and modified counts, or `No differences found.` when identical.
- **Unchanged Value Suppression**: Omits matching keys by default to keep terminal diffs clean and actionable.
- **Zero External Dependencies**: Implemented entirely with Go standard library components.

## Supported JSON Structures (v0.3.0)

- **Objects**: Flat, nested, and deeply nested objects of arbitrary depth.
- **Strings**: Quoted string comparisons.
- **Numbers**: Floating-point and integer numbers with full precision preservation (`json.Number`).
- **Booleans**: `true` and `false` literals.
- **Null Values**: `null` literal support.
- **Arrays**: Deterministic whole-value equality comparison (see [Limitations](#current-limitations)).
- **Root Primitives**: JSON documents containing a single primitive, array, or null as the root.

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
jdiff --version    # Show application version (v0.3.0)
```

### Example 1: Deeply Nested Changes

Comparing [`examples/nested-old.json`](examples/nested-old.json) and [`examples/nested-new.json`](examples/nested-new.json):

```bash
jdiff examples/nested-old.json examples/nested-new.json
```

**Output**:
```text
MODIFIED
  user.preferences.theme
    - "light"
    + "dark"

  user.profile.age
    - 28
    + 29

  user.profile.contact.email
    - "john.old@example.com"
    + "john.new@example.com"

ADDED
  user.profile.address
    + {"city":"Bengaluru","country":"India"}

REMOVED
  user.profile.contact.phone
    - "+1-555-0100"

Summary:
  Added:     1
  Removed:   1
  Modified:  3
```

### Example 2: Identical Documents

```bash
jdiff f1.json f1.json
```

**Output**:
```text
No differences found.
```

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Command completed successfully (diff output or help/version printed) |
| `1` | Operational error (missing arguments, unreadable files, invalid JSON syntax) |

## Current Limitations

- **Array Comparison**: In version `v0.3.0`, arrays are compared as complete values. If elements inside an array differ, the array path is reported as modified. Fine-grained element tracking, longest-common-subsequence (LCS) alignment, item insertions/deletions within arrays, and item reordering detection are planned for upcoming releases.
- **Color Themes & Output Formats**: Output is currently formatted in human-readable plain terminal text. ANSI color syntax and JSON Patch (RFC 6902) format are planned for future versions.

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [ ] **v0.4.0**: Advanced array diffing (element alignment, insertions, deletions, reordering).
- [ ] **v0.5.0**: Terminal ANSI color themes and customizable formatting flags.
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

## License

This project is licensed under the [MIT License](LICENSE).
