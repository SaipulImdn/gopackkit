# Contributing to GoPackKit

Thank you for your interest in contributing to GoPackKit! This document explains the process for contributing to this project.

## Code of Conduct

By participating in this project, you are expected to adhere to a code of conduct that respects all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:

1. **A clear description** of the problem
2. **Steps to reproduce** the issue
3. **Expected behavior** vs **actual behavior**
4. **Environment details** (OS, Go version, etc.)
5. **Code samples** if relevant

### Suggesting Features

To propose a new feature:

1. Check if the feature already exists in issues
2. Create a new issue with the "enhancement" label
3. Explain the use case and benefits of the feature
4. Provide implementation examples if possible

### Pull Requests

#### Before Creating a PR

1. Fork this repository
2. Create a feature branch: `git checkout -b feature/feature-name`
3. Ensure your code follows the style guidelines
4. Add tests for new code
5. Make sure all tests pass
6. Update documentation if needed

#### PR Guidelines

1. **Clear title**: Use a descriptive title
2. **Detailed description**: Explain the changes you made
3. **Link to issues**: Reference related issues
4. **Small changes**: Split large changes into several small PRs
5. **Tests included**: Every new feature must include tests

#### Code Style

- Use `gofmt` for formatting
- Follow Go naming conventions
- Write clear, descriptive variable names
- Add comments for exported functions/types
- Maximum line length: 120 characters

#### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run all checks
make check
```

#### Commit Messages

Use a clear commit message format:

```
type: short description

Longer description if needed

Fixes #123
```

Types:
- `feat`: new feature
- `fix`: bug fix
- `docs`: documentation changes
- `style`: formatting changes
- `refactor`: code refactoring
- `test`: adding tests
- `chore`: maintenance tasks

#### Example Workflow

```bash
# 1. Fork and clone the repository
git clone https://github.com/yourusername/gopackkit.git
cd gopackkit

# 2. Create a feature branch
git checkout -b feature/new-validator-rule

# 3. Make changes and commit
git add .
git commit -m "feat: add new email validation rule"

# 4. Push to your fork
git push origin feature/new-validator-rule

# 5. Create a Pull Request on GitHub
```

## Development Setup

### Prerequisites

- Go 1.20 or newer
- Git
- Make (optional, but recommended)

### Local Development

```bash
# Clone repository
git clone https://github.com/saipulimdn/gopackkit.git
cd gopackkit

# Install dependencies
make deps

# Install development tools
make install-tools

# Run tests
make test

# Run linting
make lint

# Run all checks
make check
```

### Using Docker

```bash
# Development environment
docker-compose up gopackkit-dev

# Run tests in container
docker-compose up gopackkit-test

# Run with services (MinIO, Redis)
docker-compose up -d minio redis
```

## Project Structure

```
gopackkit/
├── .github/workflows/    # CI/CD pipelines
├── config/              # Configuration module
├── httpclient/          # HTTP client module
├── jwt/                 # JWT token module
├── logger/              # Logging module
├── minio/               # MinIO client module
├── password/            # Password hashing module
├── validator/           # Validation module
├── .golangci.yml        # Linter configuration
├── Dockerfile           # Container configuration
├── Makefile             # Development commands
├── README.md            # Project documentation
├── docker-compose.yml   # Development services
└── go.mod               # Go module file
```

## Module Guidelines

When adding a new module:

1. **Create a directory** with the module name
2. **Implement an interface** consistent with other modules
3. **Add comprehensive tests** with at least 80% coverage
4. **Write documentation** with examples
5. **Update the main README.md** with module documentation

### Module Structure

```
module-name/
├── module.go           # Main implementation
├── module_test.go      # Unit tests
├── examples/           # Usage examples
└── README.md           # Module documentation
```

### Interface Design

- Keep interfaces small and focused
- Use dependency injection
- Support configuration through structs
- Provide sensible defaults
- Handle errors properly

## Security Guidelines

- **No hardcoded secrets** in code
- **Validate all inputs** to prevent injection attacks
- **Use secure defaults** in configuration
- **Avoid regex** that may be vulnerable to ReDoS
- **Review dependencies** for vulnerabilities

## Documentation

- **Code comments**: All exported functions must have comments
- **README updates**: Update README.md when adding features
- **Examples**: Provide clear usage examples
- **API docs**: Use godoc conventions

## Testing Guidelines

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result, err := FunctionName(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test Categories

1. **Unit tests**: Test individual functions
2. **Integration tests**: Test module interactions
3. **Benchmark tests**: Performance testing
4. **Example tests**: Executable documentation

### Coverage Requirements

- Minimum 80% test coverage for new modules
- 100% coverage for critical security functions
- Tests must cover edge cases and error conditions

## Release Process

1. **Update version** in documentation
2. **Update CHANGELOG.md** with changes
3. **Create a release PR** to the main branch
4. **Tag the release** after merging
5. **GitHub Actions** will automatically create the release

## Questions?

If you have questions:

1. Check existing issues on GitHub
2. Create a new issue with the "question" label
3. Join the discussion in repository discussions

## Recognition

All contributors will be acknowledged in:
- The README.md contributors section
- Release notes
- Project documentation

Thank you for your contributions! 🎉
