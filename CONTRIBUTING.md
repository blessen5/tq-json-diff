# Contributing to jdiff

Thank you for your interest in contributing to `jdiff`!

## Code Quality & Standards

1. **Zero External Dependencies**: Keep the project lightweight and implement all functionality using the Go standard library.
2. **Formatting**: All Go files must be formatted with `go fmt ./...`.
3. **Static Analysis**: Code must pass `go vet ./...` with zero diagnostics.
4. **Testing**:
   - Write comprehensive unit tests for all new functionality.
   - Run tests with `go test -v ./...`.
   - Run benchmarks with `go test -bench=. ./...`.
   - Run fuzz tests with `go test -fuzz=Fuzz ./...`.
5. **Deterministic Output**: Ensure all JSON object keys and diff outputs remain sorted and deterministic.

## Pull Request Workflow

1. Fork the repository and create a feature branch (`git checkout -b feat/my-feature`).
2. Implement your changes following idiomatic Go patterns.
3. Verify test coverage and static analysis pass locally.
4. Submit a descriptive Pull Request.
