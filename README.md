# jdiff

A fast, lightweight command-line tool built in Go that compares two JSON documents, clearly displays differences across multiple presentation formats, generates RFC 6902-compliant JSON Patch documents, and provides configurable performance and resource controls.

**Current Version:** `v0.9.0`

## Features

- **Performance & Memory Metrics (`--stats`)**: Measures precise parse times, comparison times, memory allocation, and input sizes.
- **Configurable Resource Controls**: Protect against oversized inputs with `--max-file-size` and cap difference collection with `--max-changes`.
- **Automation & CI/CD Modes**: `--quiet` and `--exit-on-diff` with standardized, script-friendly exit codes (`0`, `1`, `2`).
- **JSON Patch Generation (RFC 6902 subset)**: Generates valid JSON Patch documents with `add`, `remove`, and `replace` operations via `--output patch`.
- **In-Memory Patch Application (`jdiff apply`)**: Apply a JSON Patch document to an input JSON file and stream the resulting document.
- **Patch Verification (`--verify-patch`)**: Atomically diffs, generates a patch, applies it in memory, and verifies that the patched document matches the target.
- **RFC 6901 JSON Pointer Paths**: Correctly encodes and decodes JSON Pointer path tokens, including `~0` and `~1` special character escaping.
- **Multiple Output Formats**:
  - **`human` (default)**: Colorized, human-readable terminal output.
  - **`json`**: Strict, machine-readable JSON output with summary metrics and delta objects.
  - **`unified`**: Unified diff representation (`@@ <path>`).
  - **`patch`**: RFC 6902 JSON Patch document.
- **Selective JSON Comparison & Ignore Rules**: Ignore dynamic fields via CLI `--ignore` flags or `.jdiff.json` configuration files with wildcard support (`*`, `[*]`).
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
| `--version`, `-v` | Display application version (`jdiff v0.9.0`) |
| `--output <format>` | Select output format: `human` (default), `json`, `unified`, `patch` |
| `--output-file <file>` | Write diff/patch output directly to a file |
| `--stats` | Display performance, timing, and memory allocation metrics |
| `--max-file-size <size>` | Maximum allowed input file size (e.g. `100MB`, `10KB`, `500B`) |
| `--max-changes <N>` | Maximum number of differences to collect before truncating |
| `--exit-on-diff` | Terminate comparison immediately upon discovering differences |
| `--quiet`, `-q` | Suppress output and communicate exclusively via exit codes |
| `--verify-patch` | Generate, apply, and verify the patch in memory |
| `--ignore <path>` | Ignore a JSON path or pattern (repeatable) |
| `--config <file>` | Use a configuration file (defaults to `.jdiff.json`) |
| `--show-config` | Display active ignore configuration and exit |
| `--no-color` | Disable ANSI terminal color escape sequences |
| `--compact` | Display compact single-line diffs |
| `--verbose` | Display active ignore rules and comparison file context |
| `--summary` | Suppress individual change diffs and display only the summary counts |

---

## Exit Codes & Automation

| Code | Meaning |
|---|---|
| `0` | Success / No differences found (or patch applied/verified) |
| `1` | Differences detected between documents |
| `2` | Operational, parsing, I/O, or resource limit error |

### CI/CD Example
```bash
if jdiff --quiet config.prod.json config.staging.json; then
    echo "Configurations match."
else
    echo "Configurations differ! Review changes before deploying."
fi
```

---

## Performance Benchmarking

Run the benchmark suite:

```bash
go test -bench=. ./...
```

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

## License

This project is licensed under the [MIT License](LICENSE).
