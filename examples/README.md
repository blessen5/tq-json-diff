# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types, depths, arrays, ignore rules, and output formats.

## Available Examples

- [`basic-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-old.json) & [`basic-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-new.json)
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json) & [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)
- [`arrays-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-old.json) & [`arrays-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/arrays-new.json)
- [`ignore-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/ignore-old.json) & [`ignore-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/ignore-new.json)
- [`output-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/output-old.json) & [`output-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/output-new.json)
- [`.jdiff.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/.jdiff.json)

---

## Output Formats (v0.7.0)

### 1. Human (Default)
```bash
jdiff --output human examples/output-old.json examples/output-new.json
```

### 2. Machine-Readable JSON
```bash
jdiff --output json examples/output-old.json examples/output-new.json
```

**Output**:
```json
{
  "summary": {
    "added": 1,
    "ignored": 0,
    "modified": 3,
    "removed": 0,
    "total": 4
  },
  "changes": [
    {
      "new": true,
      "old": false,
      "path": "config.debug",
      "type": "modified"
    },
    {
      "new": 60,
      "old": 30,
      "path": "config.timeout",
      "type": "modified"
    },
    {
      "new": "/oauth",
      "path": "endpoints[2]",
      "type": "added"
    },
    {
      "new": 2,
      "old": 1,
      "path": "version",
      "type": "modified"
    }
  ]
}
```

### 3. Unified Diff Format
```bash
jdiff --output unified examples/output-old.json examples/output-new.json
```

**Output**:
```text
--- examples/output-old.json
+++ examples/output-new.json
@@ config.debug
- false
+ true

@@ config.timeout
- 30
+ 60

@@ endpoints[2]
+ "/oauth"

@@ version
- 1
+ 2
```

### 4. Stdin Piping & File Redirection
```bash
cat examples/output-old.json | jdiff --output json --output-file /tmp/diff.json - examples/output-new.json
```
