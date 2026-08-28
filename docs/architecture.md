# jdiff Architecture

`jdiff` is designed with a strictly decoupled, modular pipeline implemented entirely with Go's standard library.

---

## 1. System Pipeline

```
          ┌─────────────┐
          │  CLI Input  │ (Files, Stdin, Flags, Config)
          └──────┬──────┘
                 │
                 ▼
          ┌─────────────┐
          │ JSON Parser │ (json.Decoder with UseNumber)
          └──────┬──────┘
                 │
                 ▼
          ┌─────────────┐
          │ Diff Engine │ <─── [PathMatcher / Ignore Rules]
          └──────┬──────┘
                 │
                 ▼
          ┌─────────────┐
          │ Diff Result │ (Structured delta AST & metrics)
          └──────┬──────┘
                 │
  ┌──────────────┼──────────────┬──────────────┐
  ▼              ▼              ▼              ▼
┌───────┐  ┌───────────┐  ┌───────────┐  ┌───────────────┐
│ Human │  │ JSON Diff │  │  Unified  │  │  JSON Patch   │
└───────┘  └───────────┘  └───────────┘  └───────┬───────┘
                                                 │
                                                 ▼
                                         ┌───────────────┐
                                         │ Patch Applier │ (jdiff apply)
                                         └───────────────┘
```

---

## 2. Core Packages

### `internal/diff`
- Recursive structural comparator.
- Uses `json.Number` for lossless numeric comparisons.
- Employs `Path` AST representing exact object keys and array indices.
- Supports safety bounds (`MaxDepth`, `MaxChanges`, `EarlyExit`).

### `internal/matcher`
- Evaluates wildcard expressions (`*`, `[*]`, subtree pruning) against `Path` segments.
- Early-prunes subtree traversals for ignored paths.

### `internal/patch`
- Generates RFC 6902-compliant JSON Patch operations (`add`, `remove`, `replace`).
- Emits array removals in descending index order to ensure sequential in-memory applicability.
- Provides atomic, in-memory patch execution via `Apply` and `Verify`.

### `internal/render`
- Formats structured `DiffResult` into `human`, `json`, `unified`, or `patch` representations.
- Employs streaming output writers and ANSI color isolation.

### `internal/stats`
- Tracks memory allocations and parse/compare timings with millisecond resolution.

### `internal/config`
- Loads `.jdiff.json` and explicit `--config` configurations.
