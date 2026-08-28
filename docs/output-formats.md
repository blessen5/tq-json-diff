# Output Formats & Schemas

`jdiff` provides three distinct presentation formats for comparing JSON documents:

1. `human` (default)
2. `json` (machine-readable)
3. `unified` (unified diff format)

---

## 1. Human-Readable (`--output human`)

Hierarchical, developer-friendly terminal presentation with semantic ANSI color highlighting.

```bash
jdiff --output human old.json new.json
```

---

## 2. Machine-Readable JSON (`--output json`)

Outputs strict, deterministic JSON compliant with standard JSON parsers and automation pipelines.

```bash
jdiff --output json old.json new.json
```

### JSON Schema

```json
{
  "summary": {
    "added": 1,
    "removed": 1,
    "modified": 2,
    "ignored": 0,
    "total": 4
  },
  "changes": [
    {
      "path": "user.name",
      "type": "modified",
      "old": "John",
      "new": "James"
    },
    {
      "path": "user.email",
      "type": "added",
      "new": "james@example.com"
    },
    {
      "path": "user.legacy_id",
      "type": "removed",
      "old": 123
    }
  ]
}
```

### Properties

| Field | Type | Description |
|---|---|---|
| `summary.added` | `integer` | Count of added leaf elements / properties |
| `summary.removed` | `integer` | Count of removed leaf elements / properties |
| `summary.modified` | `integer` | Count of modified leaf elements / properties |
| `summary.ignored` | `integer` | Count of detected differences excluded by ignore rules |
| `summary.total` | `integer` | Total actual differences (`added + removed + modified`) |
| `changes[].path` | `string` | Hierarchical JSON path (e.g. `user.profile.age`, `tags[1]`) |
| `changes[].type` | `string` | Change category: `"added"`, `"removed"`, or `"modified"` |
| `changes[].old` | `any` | Original value (omitted for `"added"`) |
| `changes[].new` | `any` | Updated value (omitted for `"removed"`) |

### Summary-Only JSON (`--output json --summary`)

```bash
jdiff --output json --summary old.json new.json
```

```json
{
  "summary": {
    "added": 1,
    "removed": 0,
    "modified": 2,
    "ignored": 0,
    "total": 3
  }
}
```

---

## 3. Unified Diff Format (`--output unified`)

Provides standard unified diff style hunks (`@@ <path>`).

```bash
jdiff --output unified old.json new.json
```

```text
--- old.json
+++ new.json
@@ user.name
- "John"
+ "James"

@@ user.email
+ "james@example.com"
```

---

## Stdin & Piping

Use `-` to read one document from standard input:

```bash
cat old.json | jdiff - new.json
```

```bash
generate-payload | jdiff --output json template.json - > diff.json
```
