# jdiff Architecture & Design Specification

## Overview

`jdiff` is a lightweight, high-performance command-line utility built in Go designed to compute and display structural differences between two JSON documents.

## Design Principles

1. **Zero External Dependencies**: Use Go standard library (`encoding/json`, `fmt`, `os`, `io`, `strings`, `reflect`) to ensure maximum portability, minimal binary footprint, and fast compilation.
2. **Deterministic Traversal**: Traverse JSON trees (objects, arrays, primitives) in deterministic, sorted order to produce predictable, reproducible diffs.
3. **Clear Structural Delta**: Classify structural changes into explicit delta types:
   - Added fields/elements (`+`)
   - Removed fields/elements (`-`)
   - Modified fields/values (`~` or `-` / `+`)
   - Type mismatches (e.g. string converted to array)
4. **Standard CLI Contract**:
   - Exit code `0`: Comparison succeeded, no differences found (or help/version printed).
   - Exit code `1`: Differences detected between input documents (Phase 2+).
   - Exit code `2`: Operational error (invalid arguments, missing files, invalid JSON syntax).
5. **Modular Architecture**:
   - `internal/version`: Version and build metadata.
   - `internal/cli`: CLI argument handling, flag parsing, and help/version commands.
   - `internal/diff` *(Phase 2)*: Core AST traversal and difference computation engine.
   - `internal/formatter` *(Phase 3)*: Colorized and plain-text output formatters (terminal diff, JSON patch, summary).
   - `internal/loader` *(Phase 2)*: Safe JSON file reading and syntax validation with line/column reporting.

## Package Hierarchy

```
jdiff/
├── cmd/
│   └── jdiff/              # Canonical binary entry point
├── internal/
│   ├── cli/                # Command-line interface & argument parser
│   ├── version/            # Build & version information
│   ├── loader/             # (Phase 2) JSON parsing & validation
│   ├── diff/               # (Phase 2) Recursive diff calculation engine
│   └── formatter/          # (Phase 3) Terminal & color output formatting
├── docs/                   # Architecture and reference documentation
├── examples/               # Sample JSON datasets & test fixtures
├── tests/                  # End-to-end and integration tests
├── go.mod                  # Module specification
├── main.go                 # Root entry point
├── README.md               # User guide and quick start
├── CONTRIBUTING.md         # Contribution guidelines
├── LICENSE                 # License terms
└── .gitignore              # Ignored build artifacts
```
