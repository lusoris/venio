---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "**"
description: Bleeding Edge Stable Dependency Policy
---

# Venio Dependency & Version Management Policy

## Core Principle

**Bleeding Edge Stable:** We always use the latest **stable** versions of dependencies, not legacy versions.

- ❌ Do NOT stay multiple versions behind
- ✅ Do USE latest stable release immediately
- ✅ Do LEVERAGE new features from newer versions
- ✅ Do FIX deprecations promptly

## Version Update Strategy

### When a new stable version is released:
1. **Update immediately** (within 1 sprint if possible)
2. **Test thoroughly** with existing codebase
3. **Leverage new features** where they improve performance/security/DX
4. **Fix deprecation warnings** in same update if needed

### Deprecation Handling
- When a feature is marked deprecated in a dependency, plan its replacement in **next 2 versions max**
- Do not wait for features to be removed to start fixing
- Create issues to track deprecation fixes
- Include deprecation fixes in routine maintenance

### Security Updates
- Security patches: Apply immediately (same day if possible)
- Use Snyk to monitor continuously
- Run `snyk test` before each commit
- Zero tolerance for known vulnerabilities

## Technology Stack Versions

### Hard Constraints (Always Latest Stable)
- **PostgreSQL:** Latest stable (currently 18.1)
- **Redis:** Latest stable (currently 8.4)
- **Go:** Latest stable (currently 1.25)
- **Node.js/npm:** Latest LTS + latest stable

### Rationale
- **Performance:** Newest versions have performance improvements
- **Features:** New capabilities for competitive advantage
- **Security:** Latest patches applied automatically
- **Compatibility:** Ecosystem moves forward together
- **DX:** New tooling and debugging capabilities

## Examples

### ✅ DO
```go
// Use new Go 1.25 feature: Container-aware GOMAXPROCS
// No need to manually set GOMAXPROCS, it auto-respects cgroups

// Use Go 1.25 sync.WaitGroup.Go() convenience method
wg := sync.WaitGroup{}
wg.Go(func() {
    // goroutine work
})
```

### ✅ DO
```sql
-- Use PostgreSQL 18.1 UUID v7 (timestamp-ordered)
SELECT uuidv7() as new_id;

-- Use PostgreSQL 18.1 virtual generated columns
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT,
    display_name TEXT GENERATED ALWAYS AS (full_name) STORED
);
```

### ❌ DON'T
```go
// Bad: Using old pattern when newer version offers improvement
// Instead of manually managing with sync.Mutex, use sync.WaitGroup.Go()

var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
```

### ❌ DON'T
```
"postgres:16-alpine"  // ❌ Outdated, we use 18.1
"redis:7-alpine"      // ❌ Outdated, we use 8.4
```

## Maintenance Cadence

- **Monthly:** Check for new stable releases of major dependencies
- **Weekly:** Run Snyk scans for vulnerabilities
- **Per-commit:** Validate dependency security before push
- **Per-PR:** Update dependencies if newer versions available

## Communication

When updating versions:
1. Document changes in commit message
2. Reference new features being utilized (if any)
3. Note deprecated patterns being fixed (if any)
4. Tag as `chore: upgrade` or `feat: leverage new features`

Example commit:
```
chore: upgrade golang.org/x/net to v0.38.0

- Fixes CVE-2024-45338 (DoS vulnerability)
- Fixes CVE-2025-22872 (Input validation)
- Leverages improved TLS support in new version
```

## For AI Assistants

When generating or modifying code:
- Check current versions in `go.mod`, `package.json`, etc.
- Use latest stable APIs and patterns (not deprecated ones)
- Suggest version upgrades if new features would improve code quality
- Flag any deprecated patterns for immediate fixing
- Run security scans (`snyk test`) on generated code

This is a strategic decision for competitive advantage through staying current with ecosystem evolution.
