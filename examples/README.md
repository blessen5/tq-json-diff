# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types, depths, arrays, and selective ignore rules.

## Available Examples

- [`basic-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-old.json) & [`basic-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-new.json)
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json) & [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)
- [`arrays-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-old.json) & [`arrays-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-new.json)
- [`ignore-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/ignore-old.json) & [`ignore-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/ignore-new.json)
- [`.jdiff.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/.jdiff.json)

---

## Selective Comparison Demonstrations (v0.6.0)

### 1. Without Ignore Rules (Full Diff)
```bash
jdiff examples/ignore-old.json examples/ignore-new.json
```

Outputs modifications across `timestamp`, `request_id`, `metadata.updated_at`, `users[0].session_id`, `users[1].session_id`, `users[0].name`, and `version`.

### 2. Using CLI `--ignore` Flag
```bash
jdiff --ignore timestamp --ignore request_id examples/ignore-old.json examples/ignore-new.json
```

### 3. Using Configuration File (`--config`)
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

### 4. Viewing Active Rules (`--show-config`)
```bash
jdiff --config examples/.jdiff.json --show-config
```

**Output**:
```text
Ignore rules:
  timestamp
  request_id
  *.updated_at
  users[*].session_id
```
