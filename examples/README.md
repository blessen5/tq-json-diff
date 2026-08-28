# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types, depths, arrays, ignore rules, output formats, and JSON Patch generation.

## Available Examples

- [`basic-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-old.json) & [`basic-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-new.json)
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json) & [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)
- [`arrays-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-old.json) & [`arrays-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-new.json)
- [`ignore-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/ignore-old.json) & [`ignore-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/ignore-new.json)
- [`output-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/output-old.json) & [`output-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/output-new.json)
- [`.jdiff.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/.jdiff.json)

---

## JSON Patch Generation & Application (v0.8.0)

### 1. Generate Patch
```bash
jdiff --output patch examples/output-old.json examples/output-new.json
```

### 2. Verify Patch
```bash
jdiff --verify-patch examples/output-old.json examples/output-new.json
```

### 3. Apply Patch
```bash
jdiff --output patch --output-file diff.patch.json examples/output-old.json examples/output-new.json
jdiff apply diff.patch.json examples/output-old.json
```
