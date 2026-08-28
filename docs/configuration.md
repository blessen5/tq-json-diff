# Configuration & Ignore Rules

`jdiff` provides flexible mechanisms to ignore dynamic, environment-specific, or timestamp fields during JSON comparisons.

---

## 1. CLI Flags (`--ignore`)

Pass `--ignore` one or more times to ignore specific paths or wildcard patterns:

```bash
jdiff --ignore timestamp --ignore "users[*].session_id" old.json new.json
```

---

## 2. Configuration File (`.jdiff.json`)

By default, `jdiff` checks the working directory for a `.jdiff.json` file.

### Schema
```json
{
  "ignore": [
    "timestamp",
    "metadata.build_id",
    "servers[*].last_ping",
    "metrics.*"
  ]
}
```

### Specifying a Custom Config File
```bash
jdiff --config ./deploy/jdiff-rules.json old.json new.json
```

---

## 3. Supported Patterns

| Pattern | Description | Example Target |
|---|---|---|
| `key` | Ignores exact root key `key` and any nested `key` subtree | `{"key": ...}` |
| `a.b.c` | Exact nested object path | `{"a": {"b": {"c": 123}}}` |
| `items[0]` | Exact array element at index 0 | `{"items": ["ignored", ...}` |
| `items[*]` | All elements in an array | `{"items": [...]}` |
| `users[*].id` | Property `id` across every object in array `users` | `{"users": [{"id": 1}]}` |
| `config.*` | All immediate child keys under `config` | `{"config": {"a": 1, "b": 2}}` |

---

## 4. Viewing Active Configuration

Use `--show-config` to inspect active ignore rules:

```bash
jdiff --show-config
```
