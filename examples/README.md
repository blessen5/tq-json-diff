# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types and depths.

## Available Examples

- [`basic-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-old.json) & [`basic-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-new.json)
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json) & [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)

---

## Output Modes Demonstrations

### 1. Default Mode
```bash
jdiff examples/basic-old.json examples/basic-new.json
```
Output:
```text
MODIFIED
  age
    - 19
    + 20

  city
    - "Kochi"
    + "Bengaluru"

ADDED
  country
    + "India"

Summary:
  Added:     1
  Removed:   0
  Modified:  2
```

### 2. Compact Mode (`--compact`)
```bash
jdiff --compact examples/basic-old.json examples/basic-new.json
```
Output:
```text
MODIFIED age: 19 → 20
MODIFIED city: "Kochi" → "Bengaluru"
ADDED country: "India"

Summary:
  Added:     1
  Removed:   0
  Modified:  2
```

### 3. Summary-Only Mode (`--summary`)
```bash
jdiff --summary examples/basic-old.json examples/basic-new.json
```
Output:
```text
JSON Diff Summary

Added:     1
Removed:   0
Modified:  2
Total:     3
```

### 4. No-Color Mode (`--no-color`)
```bash
jdiff --no-color examples/basic-old.json examples/basic-new.json
```
Disables all ANSI terminal color sequences (ideal for CI/CD pipelines, log files, and redirected pipes).

### 5. Verbose Mode (`--verbose`)
```bash
jdiff --verbose examples/basic-old.json examples/basic-new.json
```
Output:
```text
Comparing:
  Old: examples/basic-old.json
  New: examples/basic-new.json

Changes:
MODIFIED
  age
    - 19
    + 20

  city
    - "Kochi"
    + "Bengaluru"

ADDED
  country
    + "India"

Summary:
  Added:     1
  Removed:   0
  Modified:  2
```
