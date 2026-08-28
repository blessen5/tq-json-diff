# jdiff Examples

This directory contains sample JSON pairs demonstrating structural differences across various JSON data types and depths.

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

### 2. Deeply Nested Object Changes
- [`nested-old.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-old.json)
- [`nested-new.json`](file:///c:/Users/bless/OneDrive/Desktop/tq-json-diff/examples/nested-new.json)

Run:
```bash
jdiff examples/nested-old.json examples/nested-new.json
```

Output:
```text
MODIFIED
  user.preferences.theme
    - "light"
    + "dark"

  user.profile.age
    - 28
    + 29

  user.profile.contact.email
    - "john.old@example.com"
    + "john.new@example.com"

ADDED
  user.profile.address
    + {"city":"Bengaluru","country":"India"}

REMOVED
  user.profile.contact.phone
    - "+1-555-0100"

Summary:
  Added:     1
  Removed:   1
  Modified:  3
```
