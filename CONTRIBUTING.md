# Contributing to Venio

Thank you for your interest in contributing to Venio! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_lusoris/venio.git
   cd venio
   ```

3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/lusoris/venio.git
   ```

4. **Install dependencies:**
   ```bash
   # Install Go 1.23+
   # Install Docker & Docker Compose
   # Install development tools
   make setup  # TODO: Add setup target to Makefile
   ```

5. **Copy environment template:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

6. **Start development environment:**
   ```bash
   make dev
   ```

### Development Workflow

1. **Create a feature branch:**
   ```bash
   git checkout develop
   git pull upstream develop
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Run tests:**
   ```bash
   make test
   make lint
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature"
   ```
   
   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `style:` - Code style changes (formatting, etc.)
   - `refactor:` - Code refactoring
   - `test:` - Test changes
   - `chore:` - Build/tooling changes

5. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** to the `develop` branch

## Coding Standards

### Go Style Guide

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `goimports` for formatting (runs automatically with pre-commit hooks)
- Run `golangci-lint` before committing
- Write tests for new features
- Aim for 80%+ test coverage
- Document exported functions and types

### Code Organization

```
internal/
├── api/          # HTTP handlers & routes
├── services/     # Business logic
├── database/     # Database layer
├── models/       # Data models
├── config/       # Configuration
└── providers/    # External API clients
```

### Testing

- Write unit tests for all new code
- Place tests in `*_test.go` files
- Use table-driven tests where appropriate
- Mock external dependencies

Example:
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## Pull Request Process

1. **Ensure CI passes** - All tests and lints must pass
2. **Update documentation** - Document new features
3. **Add tests** - Maintain >80% coverage
4. **Follow template** - Fill out the PR template completely
5. **Wait for review** - Maintainers will review your PR
6. **Address feedback** - Make requested changes
7. **Squash commits** - Before merge if requested

### PR Guidelines

- Keep PRs focused and small
- One feature/fix per PR
- Reference related issues
- Include screenshots for UI changes
- Update CHANGELOG.md for user-facing changes

## Reporting Issues

### Bug Reports

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md) and include:
- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Logs/screenshots

### Feature Requests

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md) and include:
- Problem description
- Proposed solution
- Use cases
- Alternatives considered

## Questions?

- Check existing [GitHub Discussions](https://github.com/lusoris/venio/discussions)
- Ask in our [Discord](#) (coming soon)
- Open a [Question issue](.github/ISSUE_TEMPLATE/question.md)

## License

By contributing, you agree that your contributions will be licensed under the GPL v3.0 License.

---

**Thank you for contributing to Venio!** 🎉
