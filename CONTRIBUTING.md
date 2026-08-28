# Contributing to jdiff

Thank you for your interest in contributing to `jdiff`!

## Code of Conduct

We are committed to providing a friendly, safe, and welcoming environment for all contributors. Please be respectful and constructive in discussions and reviews.

## Development Workflow

### Prerequisites

- Go 1.22 or higher installed
- Standard Go command-line tools

### Local Setup

1. Clone or initialize the repository:
   ```bash
   git clone https://github.com/example/jdiff.git
   cd jdiff
   ```

2. Verify code quality and run tests:
   ```bash
   go fmt ./...
   go vet ./...
   go test -v ./...
   go build ./...
   ```

### Coding Guidelines

- **Idiomatic Go**: Write clean, standard Go code conforming to Effective Go standards.
- **Minimal Dependencies**: Prefer the Go standard library wherever possible.
- **Testing**: Every new feature or fix must include unit tests and, where appropriate, integration tests.
- **Error Handling**: Provide clear, user-friendly error messages with proper exit codes.
- **Formatting**: Always format your code with `go fmt` before submitting pull requests.

## Submitting Pull Requests

1. Create a feature branch (`git checkout -b feature/your-feature-name`).
2. Make your changes with clear, descriptive commit messages.
3. Ensure all tests pass (`go test -v -race ./...`).
4. Submit a Pull Request describing your changes, motivation, and test coverage.
