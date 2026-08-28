# Performance, Resource Controls & Benchmarks

`jdiff` provides built-in metrics and configurable resource controls to handle larger JSON payloads safely and predictably in developer workflows and CI/CD pipelines.

---

## 1. Performance Statistics (`--stats`)

Use `--stats` to display input sizes, memory allocation, and parse/compare timings:

```bash
jdiff --stats old.json new.json
```

**Terminal Output**:
```text
Comparison Statistics:
  Old size:      12.4 MB
  New size:      13.1 MB
  Changes:       148
  Parse time:    42 ms
  Compare time:  31 ms
  Total time:    76 ms
  Allocated:     24.5 MB
```

When paired with `--output json`, statistics are embedded in the JSON object:

```json
{
  "summary": { ... },
  "changes": [ ... ],
  "statistics": {
    "old_size": 13002342,
    "new_size": 13736345,
    "old_is_stdin": false,
    "new_is_stdin": false,
    "parse_time_ms": 42,
    "compare_time_ms": 31,
    "total_time_ms": 76,
    "alloc_bytes": 25690112,
    "changes_count": 148
  }
}
```

---

## 2. Resource Controls

### Maximum File Size (`--max-file-size <size>`)
Enforce input file size limits to protect against unbounded resource consumption:

```bash
jdiff --max-file-size 100MB old.json new.json
```

Supported unit suffixes (case-insensitive):
- `B` (bytes)
- `KB`, `K` (kilobytes: $1024$ bytes)
- `MB`, `M` (megabytes: $1024 \times 1024$ bytes)
- `GB`, `G` (gigabytes: $1024 \times 1024 \times 1024$ bytes)

Exceeding the threshold halts reading immediately and exits with status code `2`.

### Maximum Changes Limit (`--max-changes <N>`)
Limit change collection to at most $N$ differences before terminating traversal early:

```bash
jdiff --max-changes 50 old.json new.json
```

When reached, the output explicitly indicates truncation:
```text
[Diff output truncated: maximum changes limit 50 reached]
```

---

## 3. Automation Modes

### Early Exit on First Difference (`--exit-on-diff`)
Terminates diff exploration upon finding the first difference:

```bash
jdiff --exit-on-diff old.json new.json
```

### Quiet Mode (`--quiet`, `-q`)
Suppresses standard diff output, communicating exclusively through exit codes:

```bash
jdiff --quiet old.json new.json
```

---

## 4. Standardized Exit Codes

| Exit Code | Meaning |
|---|---|
| `0` | Success / No differences found (or successful patch apply/verify) |
| `1` | Differences detected between documents (or patch verification failed) |
| `2` | Operational, parsing, I/O, or resource limit error |

---

## 5. Benchmarks

Run the benchmark suite using standard Go tooling:

```bash
go test -bench=. ./...
```

The benchmark suite tests:
1. `BenchmarkSmallJSON`: Small objects (5-10 properties)
2. `BenchmarkMediumJSON`: Typical service configuration files
3. `BenchmarkDeeplyNestedJSON`: 20 levels of nesting recursion
4. `BenchmarkLargeArray`: Array with 1,000 elements
5. `BenchmarkLargeObject`: Object with 1,000 keys
6. `BenchmarkFewChangesInLargeInput`: 1,000-key object with a single change

---

## 6. Architecture & Memory Notes

- `jdiff` loads and parses input JSON documents into an in-memory JSON abstract syntax tree using standard library parsers.
- Memory consumption scales proportionally with total input size.
- For streaming comparison of unbounded multi-gigabyte continuous streams, full streaming diff engines are planned for future major releases.
