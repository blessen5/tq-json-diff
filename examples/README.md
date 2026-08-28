# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types, depths, and arrays.

## Available Examples

- [`basic-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-old.json) & [`basic-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-new.json)
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json) & [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)
- [`arrays-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-old.json) & [`arrays-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-new.json)

---

## Array Comparison Demonstrations (v0.5.0)

### 1. Default Mode
```bash
jdiff examples/arrays-old.json examples/arrays-new.json
```

Output:
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

### 2. Compact Mode (`--compact`)
```bash
jdiff --compact examples/arrays-old.json examples/arrays-new.json
```

Output:
```text
MODIFIED languages[1]: "Python" → "Rust"
ADDED languages[3]: "TypeScript"
MODIFIED maintainers[0].role: "Lead" → "Creator"
ADDED maintainers[2]: {"name":"Bob","role":"Reviewer"}
REMOVED features[2]: "ansi-colors"

Summary:
  Added:     2
  Removed:   1
  Modified:  2
```

### 3. Summary-Only Mode (`--summary`)
```bash
jdiff --summary examples/arrays-old.json examples/arrays-new.json
```

Output:
```text
JSON Diff Summary

Added:     2
Removed:   1
Modified:  2
Total:     5
```
