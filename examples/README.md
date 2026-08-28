# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types.

## Available Examples

### 1. Basic Object Changes
- [`basic-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-old.json)
- [`basic-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/basic-new.json)

Run:
```bash
jdiff examples/basic-old.json examples/basic-new.json
```

Output:
```text
MODIFIED:
    age
        - 19
        + 20

    city
        - "Kochi"
        + "Bengaluru"

ADDED:
    country
        + "India"
```

### 2. Nested Hierarchy Changes
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json)
- [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)

Run:
```bash
jdiff examples/nested-old.json examples/nested-new.json
```

Output:
```text
MODIFIED:
    config.database.name
        - "auth_dev"
        + "auth_prod"

    config.database.pool
        - 5
        + 20

    config.debug
        - true
        + false

    config.host
        - "localhost"
        + "0.0.0.0"

    version
        - 1
        + 2

ADDED:
    config.database.ssl
        + true
```
