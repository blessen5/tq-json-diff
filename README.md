# jdiff

[![Go Version](https://img.shields.io/badge/go-1.20%2B-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/badge/release-v1.0.0-green.svg)](https://github.com/blessen5/tq-son-diff)

A fast, lightweight, zero-dependency command-line utility built in Go that compares two JSON documents, clearly displays structural differences across multiple presentation formats, generates RFC 6902-compliant JSON Patch documents, and provides configurable performance and resource controls.

**Current Version:** `v1.0.0`

---

## Features

- **Deep Structural JSON Diff**: Recursively compares arbitrary nested objects and arrays with granular path reporting (`user.profile.name`, `tags[2]`).
- **Multiple Output Formats**:
  - **`human` (default)**: Colorized, developer-friendly terminal diff.
  - **`json`**: Machine-readable JSON summary and delta stream.
  - **`unified`**: Unified patch representation (`@@ <path>`).
  - **`patch`**: RFC 6902 JSON Patch document (`add`, `remove`, `replace`).
- **In-Memory JSON Patch Application & Verification**:
  - `jdiff apply <patch.json> <input.json>` applies changes atomically.
  - `--verify-patch` validates round-trip equivalence in memory.
- **Resource Safeguards & Security**:
  - `--max-file-size`: Bounded file size protection (`B`, `KB`, `MB`, `GB`).
  - `--max-changes`: Caps difference collection to prevent runaway output.
  - `--max-depth`: Protects against stack exhaustion on hostile/deep recursion.
- **Selective Comparison & Ignore Rules**: Filter out timestamps and ephemeral fields via CLI `--ignore` or `.jdiff.json` configuration files with wildcard support (`*`, `[*]`).
- **Automation & Scripting**: `--quiet` and `--exit-on-diff` with standardized exit codes (`0`, `1`, `2`).
- **Zero External Dependencies**: Built entirely with Go standard library packages.

---

## Quick Start

### Installation

```bash
# Install directly via Go
go install ./cmd/jdiff

# Or build locally
go build -o jdiff .
```

### Basic Comparison

```bash
# Compare two JSON files
jdiff old.json new.json

# Machine-readable JSON output
jdiff --output json old.json new.json

# Generate RFC 6902 JSON Patch
jdiff --output patch old.json new.json

# Apply a generated patch
jdiff apply diff.patch.json old.json
```

---

## CLI Options & Flags

```text
jdiff [options] <old.json> <new.json>
jdiff apply [options] <patch.json> <input.json>
```

| Flag | Description |
|---|---|
| `--help`, `-h` | Display usage and available options |
| `--version`, `-v` | Display application version (`jdiff v1.0.0`) |
| `--output <format>` | Select output format: `human` (default), `json`, `unified`, `patch` |
| `--output-file <file>` | Write diff/patch output directly to a file |
| `--stats` | Display performance, timing, and memory allocation metrics |
| `--max-file-size <size>` | Maximum allowed input file size (e.g. `100MB`, `10KB`, `500B`) |
| `--max-changes <N>` | Maximum number of differences to collect before truncating |
| `--max-depth <N>` | Maximum allowed JSON recursion depth (default: `1000`) |
| `--exit-on-diff` | Terminate comparison immediately upon discovering differences |
| `--quiet`, `-q` | Suppress output and communicate exclusively via exit codes |
| `--verify-patch` | Generate, apply, and verify the patch in memory |
| `--ignore <path>` | Ignore a JSON path or wildcard pattern (repeatable) |
| `--config <file>` | Use a configuration file (defaults to `.jdiff.json`) |
| `--show-config` | Display active ignore configuration and exit |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Display active ignore rules and comparison file context |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## Standardized Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success / No differences found (or patch applied/verified) |
| `1` | Differences detected between documents |
| `2` | Operational, parsing, I/O, or resource limit error |

### Automation / CI/CD Example
```bash
if jdiff --quiet config.prod.json config.staging.json; then
    echo "Configurations match."
else
    echo "Configurations differ! Review changes before deployment."
fi
```

---

## Documentation

- [Architecture Overview](docs/architecture.md)
- [Configuration & Ignore Rules](docs/configuration.md)
- [Output Formats & JSON Schema](docs/output-formats.md)
- [JSON Patch (RFC 6902) Guide](docs/json-patch.md)
- [Performance & Resource Controls](docs/performance.md)
- [Development, Testing & Fuzzing Guide](docs/development.md)

---

## Roadmap

- [x] **v0.1.0**: CLI scaffolding, versioning, documentation, test harness.
- [x] **v0.2.0**: Core structural diff engine, object traversal, primitive comparisons, deterministic output.
- [x] **v0.3.0**: Deep recursive comparison, JSON Path engine, root primitive support, explicit type change detection, change summaries.
- [x] **v0.4.0**: Professional terminal presentation (ANSI colors, `--no-color`, `--compact`, `--verbose`, `--summary`).
- [x] **v0.5.0**: Granular index-based array comparison, array paths (`users[0].name`), nested arrays, and arrays of objects.
- [x] **v0.6.0**: Ignore rules, wildcard paths (`*.key`, `users[*].id`), configuration files (`.jdiff.json`), and `--show-config`.
- [x] **v0.7.0**: Multiple output formats (`human`, `json`, `unified`), `--output-file`, and standard input (`-`) piping.
- [x] **v0.8.0**: JSON Patch generation (RFC 6902), patch application (`jdiff apply`), and verification (`--verify-patch`).
- [x] **v0.9.0**: Performance metrics (`--stats`), resource controls (`--max-file-size`, `--max-changes`), `--exit-on-diff`, `--quiet`, and benchmarks.
- [x] **v1.0.0**: Production-ready stable release, deep nesting protection (`--max-depth`), native Go fuzz testing, property invariants, and architecture documentation.

---

## License

This project is licensed under the [MIT License](LICENSE).
