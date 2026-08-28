# Advanced Features in jdiff

`jdiff` provides a suite of advanced capabilities designed for high-assurance API contract validation, data pipeline testing, and modern developer workflows.

---

## 1. Semantic Schema & Breaking Change Analyzer

Traditional diff tools only report value changes. `jdiff` classifies differences according to **API contract compatibility**:

- **`[BREAKING]`**: Missing fields, removed properties, or mutated data types (`string` $\to$ `number`, `object` $\to$ `array`, `any` $\to$ `null`) that cause deserialization failures in downstream clients.
- **`[ADDITIVE]`**: Non-breaking additions (e.g., newly introduced endpoints or optional keys).
- **`[VALUE_CHANGE]`**: Value-level modifications where data types remain identical.

### Usage
```bash
# Display breaking change analysis report
jdiff --breaking old_api.json new_api.json

# CI/CD Guard: exit code 1 if breaking changes exist, 0 if 100% backward-compatible
jdiff --check-breaking old_api.json new_api.json
```

---

## 2. Intelligent Array Alignment & Identity Matching

Standard JSON comparison tools match array elements strictly by ordinal index (`old[i]` vs `new[i]`). If an item is prepended or reordered, standard diff engines falsely report every element as modified.

`jdiff` supports automatic and explicit identity-based array matching:
- **`--array-match auto`**: Automatically inspects object items in arrays for primary key fields (`id`, `_id`, `uuid`, `key`, `name`, `slug`, `code`, `email`, `username`).
- **`--array-key <field>`**: Match on a specific unique property name.
- Accurately reports internal item mutations and records moved/reordered items without false-positive full object replacements.

### Usage
```bash
# Auto-detect identity keys in array objects
jdiff --array-match auto old_users.json new_users.json

# Align array elements by specific property
jdiff --array-key uuid old_manifest.json new_manifest.json
```

---

## 3. Numeric & Temporal Fuzzy Tolerance

Eliminates false alarms in automated test suites caused by floating-point precision jitter or clock drift:

- **`--numeric-tolerance <val|percent>`**:
  - Absolute delta: `--numeric-tolerance 0.001` (matches $|a - b| \le 0.001$)
  - Relative percentage: `--numeric-tolerance 1%` (matches $|a - b| / |a| \le 0.01$)
- **`--time-tolerance <duration>`**:
  - Compares ISO-8601 timestamps and ignores differences within duration: `--time-tolerance 5s`, `--time-tolerance 1m`.

### Usage
```bash
# Ignore floating point precision noise
jdiff --numeric-tolerance 0.0001 metrics_old.json metrics_new.json

# Ignore timestamp clock drift between servers
jdiff --time-tolerance 10s nodeA_event.json nodeB_event.json
```

---

## 4. Reverse Rollback Inverse Patch Generator

`jdiff` can compute the exact **mathematical inverse** RFC 6902 JSON Patch that mutates `new.json` back into `old.json`.

### Usage
```bash
# 1. Generate the inverse rollback patch
jdiff --output rollback --output-file rollback.patch.json old_state.json new_state.json

# 2. Revert the new document to the old document anytime!
jdiff apply rollback.patch.json new_state.json
```

---

## 5. Standalone Interactive HTML Visualizer

Generates a modern, single-file HTML report with side-by-side collapsible tree views, real-time search filtering, colored badges, and a one-click *"Copy JSON Patch"* button — with **zero CDN or external network dependencies**:

### Usage
```bash
jdiff --output html --output-file diff_report.html old_config.json new_config.json
```
