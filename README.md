# jdiff

[![Go Version](https://img.shields.io/badge/go-1.20%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/badge/release-v1.0.0-green.svg)](https://github.com/blessen5/tq-son-diff)

A fast, lightweight, zero-dependency command-line utility and Go library that performs deep structural JSON comparison, semantic API breaking change analysis, intelligent identity-based array alignment, RFC 6902 JSON patch and rollback generation, and standalone interactive HTML diff reporting.

---

## Key Features

### 🛡️ Semantic Schema & Breaking Change Analyzer
- Categorizes changes into **`[BREAKING]`** (removed fields, mutated data types), **`[ADDITIVE]`** (backward-compatible additions), and **`[VALUE_CHANGE]`** (safe value updates).
- **`--check-breaking`**: CI/CD automation guard that returns exit code `1` only if breaking API contract changes exist.

### 🧠 Intelligent Identity-Based Array Alignment
- Eliminates false-positive full replacements when arrays are reordered or elements are inserted.
- **`--array-match auto`**: Automatically discovers object primary keys (`id`, `uuid`, `key`, `name`, `slug`, `code`).
- **`--array-key <field>`**: Aligns array elements on any custom property.

### 🎯 Numeric & Temporal Fuzzy Tolerance
- **`--numeric-tolerance <val|percent>`**: Eliminates floating-point precision jitter (e.g. `0.001` or `1%`).
- **`--time-tolerance <duration>`**: Compares ISO-8601 timestamps with duration tolerance (e.g. `5s`, `1m`) to ignore clock drift.

### ⏪ Reverse Rollback Inverse Patch Generator
- **`--output rollback`**: Computes the exact mathematical inverse RFC 6902 JSON Patch (`undo.patch.json`) to safely revert modified documents back to their original state.

### 🌐 Zero-Dependency Interactive HTML Visualizer
- **`--output html`**: Generates a self-contained, single-file HTML report with side-by-side tree views, real-time search filtering, and one-click JSON patch copying.

### ⚡ Deep Structural Diff Engine
- Recursively compares arbitrary nested objects and arrays with granular path reporting (`user.profile.name`, `tags[2]`).
- **Multiple Output Formats**: `human` (colorized terminal), `json` (machine-readable), `unified` (`@@ <path>`), `patch` (RFC 6902), `rollback` (inverse patch), and `html`.

### 🔄 In-Memory Patch Application & Verification
- `jdiff apply <patch.json> <input.json>`: Applies JSON Patch documents atomically in-memory.
- `--verify-patch`: Validates round-trip equivalence in memory.

### 🔒 Resource Controls & Security Boundaries
- `--max-file-size`: Bounded file size protection (`B`, `KB`, `MB`, `GB`).
- `--max-changes`: Caps difference collection to prevent runaway output.
- `--max-depth`: Protects against stack exhaustion on hostile/adversarial recursion.

### 🎯 Selective Comparison & Ignore Rules
- Filter out timestamps, IDs, and dynamic tokens via CLI `--ignore` or `.jdiff.json` configuration files with wildcard support (`*`, `[*]`).

---

## Quick Start

### Installation

```bash
# Build standalone binary locally
go build -o jdiff .

# Or install to $GOPATH/bin
go install ./cmd/jdiff
```

### Usage Examples

#### 1. Basic Structural Diff
```bash
jdiff old.json new.json
```

#### 2. Schema Breaking Change Analysis
```bash
# Display breaking change analysis report
jdiff --breaking old_api.json new_api.json

# Exit with code 1 if breaking changes exist, 0 if safe/additive
jdiff --check-breaking old_api.json new_api.json
```

#### 3. Smart Array Alignment by Identity Key
```bash
# Auto-detect primary keys (id, uuid, name, etc.)
jdiff --array-match auto old_users.json new_users.json

# Align on custom property
jdiff --array-key uuid old_manifest.json new_manifest.json
```

#### 4. Fuzzy Numeric & Timestamp Tolerance
```bash
# Ignore float rounding jitter within 0.001
jdiff --numeric-tolerance 0.001 metrics_old.json metrics_new.json

# Ignore timestamp clock drift within 5 seconds
jdiff --time-tolerance 5s event_old.json event_new.json
```

#### 5. Generate and Apply RFC 6902 JSON Patches
```bash
# Generate patch
jdiff --output patch --output-file update.patch.json old.json new.json

# Apply patch
jdiff apply update.patch.json old.json
```

#### 6. Generate Inverse Rollback Patch
```bash
# Generate inverse patch to undo changes
jdiff --output rollback --output-file undo.patch.json old.json new.json

# Revert document
jdiff apply undo.patch.json new.json
```

#### 7. Standalone HTML Visualizer
```bash
jdiff --output html --output-file report.html old.json new.json
```

#### 8. Performance & Memory Statistics
```bash
jdiff --stats old.json new.json
```

---

## CLI Options & Flags

```text
Usage:
  jdiff [options] <old.json> <new.json>
  jdiff apply [options] <patch.json> <input.json>

Options:
  --help                 Show help
  --version              Show version
  --output <format>      Output format: human, json, unified, patch, rollback, html (default: human)
  --output-file <file>   Write output to a file instead of stdout
  --verify-patch         Generate, apply, and verify the patch
  --breaking             Analyze and display API/schema breaking changes
  --check-breaking       Exit with status 1 if breaking changes exist, 0 if backward-compatible
  --array-match <mode>   Array alignment mode: index, auto, key (default: index)
  --array-key <field>    Object field to align arrays on (e.g. id, uuid, key)
  --numeric-tolerance <v> Ignore numeric float drift within delta or percent (e.g. 0.01, 1%)
  --time-tolerance <dur> Ignore timestamp drift within duration (e.g. 5s, 1m)
  --stats                Display performance and memory statistics
  --max-file-size <size> Maximum allowed input file size (e.g. 100MB, 10KB, 500B)
  --max-changes <N>      Maximum number of differences to collect before truncating
  --max-depth <N>        Maximum allowed JSON recursion depth (default: 1000)
  --exit-on-diff         Terminate comparison immediately upon discovering differences
  --quiet, -q            Suppress output and communicate exclusively via exit codes
  --no-color             Disable colored output
  --compact              Display compact diff output
  --verbose              Display additional comparison information
  --summary              Display only the change summary
  --ignore <path>        Ignore a JSON path (can be specified multiple times)
  --config <file>        Use a configuration file (defaults to .jdiff.json)
  --show-config          Show active comparison configuration
```

---

## Standardized Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success / No differences found (or patch applied / 100% backward-compatible in `--check-breaking`) |
| `1` | Differences detected between documents (or breaking changes found in `--check-breaking`) |
| `2` | Operational, parsing, I/O, or resource limit error |

---

## Comprehensive Documentation

- [Advanced Features Guide](docs/advanced-features.md)
- [Architecture Overview](docs/architecture.md)
- [Configuration & Ignore Rules](docs/configuration.md)
- [Output Formats & JSON Schema](docs/output-formats.md)
- [JSON Patch (RFC 6902) Guide](docs/json-patch.md)
- [Performance & Resource Controls](docs/performance.md)
- [Development, Testing & Fuzzing Guide](docs/development.md)

---

## License

This project is licensed under the [MIT License](LICENSE).
