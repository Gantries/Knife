# Contributing to knife

Thank you for your interest in contributing to knife! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

- Go 1.21.5 or later
- Git
- Make (optional, for using Makefile commands)

### Setting Up Your Development Environment

1. **Fork the repository** on GitHub and clone your fork:

```bash
git clone https://github.com/your-username/knife.git
cd knife
```

1. **Add the upstream remote**:

```bash
git remote add upstream https://github.com/gantries/knife.git
```

1. **Install dependencies**:

```bash
go mod download
```

1. **Run tests** to verify your setup:

```bash
go test ./...
```

## Development Workflow

1. **Create a branch** for your work:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

1. **Make your changes** following our [Coding Standards](#coding-standards).

2. **Write tests** for your changes. See our [Testing guidelines](#testing).

3. **Run all tests**:

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

1. **Run linters** (if configured):

```bash
go vet ./...
golangci-lint run  # if you have golangci-lint installed
```

1. **Commit your changes** with a clear message:

```bash
git commit -m "feat: add new feature description"
```

1. **Push to your fork**:

```bash
git push origin feature/your-feature-name
```

1. **Create a Pull Request** on GitHub.

## Coding Standards

### Go Conventions

- Follow standard [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format your code

### Naming

- Use clear, descriptive names
- Package names should be lowercase, single words
- Exported functions should have documentation comments
- Acronyms should be capitalized (e.g., `HTTP`, `URL`)

### Documentation

- All exported functions, types, and constants must have documentation
- Use godoc comment format (starting with the name)
- Include usage examples for complex APIs

Example:

```go
// Serialize converts the given value to JSON bytes.
// It returns an error if the value cannot be serialized.
func Serialize[T any](v T) (buf []byte, err error) {
    // implementation
}
```

### Error Handling

- Always handle errors, never ignore them
- Use wrapped errors with context: `fmt.Errorf("operation failed: %w", err)`
- Prefer returning errors over panicking in library code

### Testing

- Aim for 80%+ test coverage
- Write table-driven tests for multiple scenarios
- Use `t.Run()` for sub-tests
- Mock external dependencies

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./pkg/orm/...

# Run verbose tests
go test -v ./...
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

```go
func TestYourFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid input",
            input: "test",
            want:  "expected",
        },
        {
            name:    "invalid input",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("YourFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Submitting Changes

### Pull Request Guidelines

1. **Title**: Use a clear, descriptive title (e.g., "feat: add Oracle dialect support")

2. **Description**: Include:
   - What changes were made and why
   - Related issues (using `#issue-number`)
   - Screenshots for UI changes (if applicable)
   - Breaking changes (if any)

3. **Commits**: Keep commits atomic and well-formatted

4. **CI Checks**: Ensure all CI checks pass before requesting review

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>: <description>

[optional body]

[optional footer]
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test changes
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `ci`: CI/CD changes
- `chore`: Other changes

Examples:

```text
feat: add support for Oracle Instant Client

Fixes #123

BREAKING CHANGE: The Oracle dialect API has changed
```

### Review Process

1. Automated checks must pass
2. At least one maintainer approval required
3. Address all review comments
4. Squash commits if requested by maintainer

## Reporting Issues

### Before Creating an Issue

1. Search existing issues to avoid duplicates
2. Check if the issue is already fixed in the latest version

### Creating a Bug Report

Include:

- **Title**: Clear, concise description
- **Description**: What happened and what you expected
- **Steps to reproduce**: Minimal reproduction case
- **Environment**: Go version, OS, database version
- **Logs**: Relevant error messages or stack traces

### Feature Requests

Include:

- **Use case**: What problem would this solve?
- **Proposed solution**: How should it work?
- **Alternatives**: What other approaches did you consider?
- **Additional context**: Examples, mockups, etc.

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **Discussions**: For questions and general discussion (if enabled)

Thank you for contributing to knife!
