# Development & Testing Guide

## 1. Prerequisites
- **Go**: 1.20+ (standard Go toolchain)
- No external libraries required (100% standard library).

---

## 2. Building

Build the `jdiff` binary locally:

```bash
go build -o jdiff .
```

---

## 3. Running Tests

### Unit & Integration Tests
```bash
go test -v ./...
```

### Static Analysis
```bash
go vet ./...
```

### Benchmarks
```bash
go test -bench=. -benchmem ./...
```

### Native Fuzz Testing
```bash
# Fuzz JSON comparison engine
go test -fuzz=FuzzCompareBytes -fuzztime=10s ./internal/diff

# Fuzz JSON Pointer engine
go test -fuzz=FuzzPointer -fuzztime=10s ./internal/patch

# Fuzz Patch Application
go test -fuzz=FuzzPatchApply -fuzztime=10s ./internal/patch
```
