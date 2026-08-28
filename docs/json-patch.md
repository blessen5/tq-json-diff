# JSON Patch (RFC 6902) Guide

`jdiff` supports generating and applying RFC 6902-compliant JSON Patch documents.

---

## 1. Generating a JSON Patch (`--output patch`)

To generate an RFC 6902 JSON Patch document between two JSON files:

```bash
jdiff --output patch old.json new.json
```

Output:
```json
[
  {
    "op": "replace",
    "path": "/version",
    "value": 2
  },
  {
    "op": "add",
    "path": "/database/pool_size",
    "value": 20
  },
  {
    "op": "remove",
    "path": "/legacy_token"
  }
]
```

---

## 2. Supported Operations

| Operation | Description | Schema |
|---|---|---|
| `add` | Adds a property to an object or inserts an element into an array | `{"op": "add", "path": "/...", "value": ...}` |
| `remove` | Removes a property from an object or deletes an array element | `{"op": "remove", "path": "/..."}` |
| `replace` | Replaces an existing value in an object or array | `{"op": "replace", "path": "/...", "value": ...}` |

> [!NOTE]
> Operations `move`, `copy`, and `test` are not currently emitted by `jdiff`.

---

## 3. RFC 6901 JSON Pointer Paths

All patch operations reference document paths using RFC 6901 JSON Pointer notation:

- Object properties: `/user/profile/name`
- Array indices: `/items/0`, `/users/2/name`
- Special character escaping:
  - `~` is escaped as `~0`
  - `/` is escaped as `~1`
  - Example: Property `"a/b"` becomes `/a~1b`

---

## 4. Array Patch Ordering Semantics

When multiple elements in an array are removed or added, `jdiff` orders patch operations safely:

1. **Replacements** (`replace`) are executed on existing index targets.
2. **Removals** (`remove`) on array elements are executed in **descending index order** (e.g. `/items/3` before `/items/2`) so index shifting does not alter remaining targets.
3. **Additions** (`add`) are executed in ascending order.

---

## 5. Applying a Patch (`jdiff apply`)

Apply a generated JSON Patch document to an input JSON file:

```bash
jdiff apply patch.json old.json
```

Or write the patched result directly to a file:

```bash
jdiff apply --output-file new-state.json patch.json old.json
```

---

## 6. Patch Verification (`--verify-patch`)

Validate that generating a patch from `old.json` and applying it produces a result structurally equivalent to `new.json`:

```bash
jdiff --verify-patch old.json new.json
```

Output:
```text
Patch verification successful.
```
