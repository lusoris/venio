# Branch-Specific Documentation Strategy

This document explains the branch-specific documentation strategy for Venio.

## Overview

Venio uses a **tiered documentation structure** to separate concerns and ensure that sensitive development documentation is not exposed to end users.

## Branch Structure

```
main (production)
├── docs/
│   ├── INDEX.md              ✅ Available
│   ├── user/                 ✅ Available (end-user docs)
│   ├── admin/                ✅ Available (admin/deployment docs)
│   └── dev/                  ❌ EXCLUDED (developer docs)
│
develop (development)
├── docs/
│   ├── INDEX.md              ✅ Available
│   ├── user/                 ✅ Available
│   ├── admin/                ✅ Available
│   └── dev/                  ✅ Available (full developer docs)
```

## Documentation Categories

### 1. User Documentation (`docs/user/`)

**Audience:** End users
**Availability:** Both `main` and `develop` branches
**Content:**
- Getting started guide
- FAQ
- Feature usage tutorials
- Troubleshooting for users

**Files:**
- `getting-started.md` - Account setup, basic usage
- `faq.md` - Common questions

### 2. Administrator Documentation (`docs/admin/`)

**Audience:** System administrators, DevOps
**Availability:** Both `main` and `develop` branches
**Content:**
- Configuration reference
- Deployment guides
- Monitoring setup
- Backup/restore procedures

**Files:**
- `configuration.md` - Environment variables, settings
- `deployment.md` - Docker, Kubernetes, production setup

### 3. Developer Documentation (`docs/dev/`)

**Audience:** Software developers, contributors
**Availability:** **ONLY in `develop` branch**
**Content:**
- Architecture details
- API documentation
- Development setup guides
- Best practices
- Security hardening
- Internal implementation details

**Files:**
- `api.md` - API endpoints and authentication
- `architecture.md` - System design, components
- `best-practices.md` - Framework-specific patterns (Gin, pgx, Next.js, etc.)
- `security-hardening.md` - OWASP Top 10, security guidelines
- `development.md` - Development environment setup
- `windows-setup.md` - Windows-specific setup
- `project-guidelines.md` - Coding standards, AI instructions
- `project-standards.md` - Tech stack versions, CalVer
- `PROJECT_STATUS.md` - Implementation status, roadmap

## Rationale

### Why Separate Developer Docs?

1. **Security:** Internal architecture details shouldn't be public
2. **Clarity:** End users don't need developer documentation
3. **Maintenance:** Allows rapid iteration on dev docs without affecting releases
4. **Branch Protection:** Developer docs evolve with `develop` branch

### Why Keep User/Admin Docs in Both?

1. **Accessibility:** Users of production releases need documentation
2. **Versioning:** Documentation matches the release version
3. **Searchability:** End users can find docs in GitHub's `main` branch

## Implementation

### `.gitignore` in `main` Branch

To exclude `docs/dev/` from the `main` branch:

```gitignore
# Exclude developer documentation (available only in develop)
docs/dev/
```

### Syncing Process

When merging `develop` → `main`:

1. **User docs** (`docs/user/`) → Synced
2. **Admin docs** (`docs/admin/`) → Synced
3. **Dev docs** (`docs/dev/`) → **NOT synced** (gitignored in main)
4. **INDEX.md** → Synced (with branch-aware content)

### Maintaining INDEX.md

`docs/INDEX.md` is branch-aware:

```markdown
## 💻 Developer Documentation ([dev/](dev/))
**Only available in `develop` branch** (not synced to `main`).
```

This message alerts readers that dev docs are in the development branch.

## Workflow

### Adding New Documentation

**User or Admin Docs:**
1. Create/edit in `docs/user/` or `docs/admin/`
2. Commit to `develop`
3. Will be synced to `main` on merge

**Developer Docs:**
1. Create/edit in `docs/dev/`
2. Commit to `develop`
3. **Will NOT be synced to `main`** (by design)

### Updating Documentation

**Routine Updates:**
- Edit in `develop` branch
- Changes auto-sync to `main` (user/admin only)

**Major Refactoring:**
- Update all three categories
- Test links across categories
- Verify INDEX.md references

## Benefits

1. **Security:** Internal details stay internal
2. **Simplicity:** Users see only relevant docs
3. **Flexibility:** Dev docs can change rapidly
4. **Clarity:** Clear separation of concerns

## Revision

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial branch strategy documentation |
