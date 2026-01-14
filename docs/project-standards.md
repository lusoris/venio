# Venio Project Standards

## Version Numbering

**Format:** `YYYY.MM.PATCH`

- **YYYY**: Year (e.g., 2026)
- **MM**: Month (e.g., 01 for January)
- **PATCH**: Incremental patch number starting from 0

**Examples:**
- `2026.01.0` - January 2026, initial release
- `2026.01.1` - January 2026, first patch
- `2026.02.0` - February 2026, new monthly release

**Rationale:**
- Calendar Versioning (CalVer) provides clear temporal context
- Patch number allows multiple releases within a month
- Better tracking for continuous delivery and monthly release cycles
- More informative than traditional SemVer for this project

## Technology Stack

### Database

**PostgreSQL Version:** 18.1

**Policy:** Always use the latest stable version of PostgreSQL

- Ensures access to latest features and performance improvements
- Security patches and bug fixes
- Better query optimization and indexing capabilities
- As long as the version is marked as stable, we upgrade

**Current Version:** `postgres:18.1-alpine` (Docker image)

### Dependencies

**Go Version:** 1.25.5

**Frontend:**
- Next.js 15.1.6
- React 19

**Principle:** Use bleeding-edge stable versions for all dependencies
- If a version is marked stable, we use it
- Allows us to leverage latest features and performance improvements
- Requires thorough testing but provides competitive advantage

## Development Workflow

### Git Hooks (Lefthook)

Pre-commit hooks run automatically on `git commit`:
- **format**: Auto-format Go files with `goimports`
- **lint**: Run `golangci-lint` with auto-fix
- **test**: Run all tests

Pre-push hooks run automatically on `git push`:
- **test**: Run full test suite
- **lint**: Final linting check

**Bypassing Hooks:**
Only use `--no-verify` when necessary (e.g., WIP commits, CI issues)

### Code Quality

- All Go code must pass `golangci-lint`
- 100% formatted with `goimports` and `gofmt`
- Minimum test coverage: TBD
- All new features require tests

## Deployment

### Environment Management

- **development**: Local development with hot reload
- **staging**: Pre-production testing environment
- **production**: Live production environment

### Release Process

1. Version bump following YYYY.MM.PATCH format
2. Update CHANGELOG.md
3. Tag release: `git tag v2026.01.0`
4. Push with tags: `git push --tags`
5. CI/CD handles building and deployment

## Docker Standards

### Image Tags

- Latest stable version for all base images
- Alpine variants preferred for smaller image size
- Multi-platform builds (linux/amd64, linux/arm64)

### Compose Files

- `docker-compose.yml` - Production configuration
- `docker-compose.dev.yml` - Development overrides
- `docker-compose.test.yml` - Testing environment

## Documentation

- All public APIs must be documented with OpenAPI/Swagger
- README.md must be kept up-to-date
- Architecture decisions documented in docs/architecture.md
- Development setup in docs/development.md

## Security

- Never commit secrets or credentials
- Use environment variables for all sensitive data
- Regular dependency updates for security patches
- Follow OWASP best practices
