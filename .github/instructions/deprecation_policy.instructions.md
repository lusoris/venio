---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "**"
description: Deprecation Management Guidelines
---

# Deprecation Management & Feature Lifecycle

## Philosophy

**Pro-Active Modernization:** We fix deprecations BEFORE they become breaking changes, not after.

## Deprecation Timeline

### When Feature is Marked Deprecated

1. **Version N:** Feature marked as deprecated
   - Status: ⚠️ Can still use, but notice is given
   - Action: Create GitHub issue with label `deprecation`
   - Priority: Add to next sprint backlog

2. **Version N+1:** Deprecation continues
   - Status: ⚠️ Still works, warnings may increase
   - Action: Begin migration planning
   - Priority: Schedule migration in roadmap

3. **Version N+2:** Final deprecation warning
   - Status: ⚠️ Usually the last version before removal
   - Action: MUST migrate before N+3
   - Priority: Critical - schedule immediately

4. **Version N+3+:** Feature removed
   - Status: ❌ Code breaks
   - Action: Would cause build failures
   - **Our goal:** Never reach this state

## Required Actions

### At Version N (Deprecation Announced)
- [ ] Document deprecated pattern in code comment
- [ ] Create migration plan/issue
- [ ] Assess impact (how many usages)
- [ ] Plan replacement approach

Example:
```go
// DEPRECATED: Use context.WithoutCancel() instead (Go 1.21+)
// Migration: Replace ctx.Background() patterns where cancellation isn't needed
// Timeline: Must migrate before Go 1.26
func legacyPattern() {
    ctx := context.Background()
    // ...
}
```

### At Version N+1 (Deprecation Continues)
- [ ] Schedule migration in sprint
- [ ] Begin code refactoring
- [ ] Update tests to use new pattern

### At Version N+2 (Final Warning)
- [ ] Migration MUST be complete
- [ ] All tests passing
- [ ] Code review done
- [ ] Deployed to production

## Examples from Current Stack

### Go Deprecations to Watch

| Pattern | Deprecated | Removal | Status | Action |
|---------|-----------|---------|--------|--------|
| `sync.Mutex` patterns | — | — | ✅ Use `sync.WaitGroup.Go()` where applicable (Go 1.25) | Refactor goroutine patterns |
| Old `crypto/elliptic` methods | Go 1.21 | Go 1.26 | ✅ Complete | Use new `crypto/ecdsa` functions |
| SHA-1 in TLS 1.2 | Go 1.25 | Go 1.26+ | ⚠️ Actively using? | Add flag if needed: `GODEBUG=tlssha1=1` |

### PostgreSQL Deprecations

| Feature | Deprecated | Removed | Status | Action |
|---------|-----------|---------|--------|--------|
| MD5 password auth | PG 18 | PG 19+ | ⚠️ Plan migration | Use SCRAM/OAuth instead |
| Old `pg_upgrade` behavior | PG 18 | PG 19 | ✅ Already migrated | Using new statistics preservation |

### Node.js/Next.js Deprecations
- Track in `package.json` with `npm audit`
- Review Next.js breaking changes each major version
- Plan Next.js migrations in quarterly reviews

## Tracking System

### Issue Labels
```
deprecation          # Tracks deprecated features needing migration
breaking-change      # Breaking change in dependencies
upgrade              # Version upgrade needed
security             # Security deprecations (highest priority)
```

### Issue Title Format
```
[DEPRECATION] Feature Name - Removal in Version X.Y

Example:
[DEPRECATION] sync.Mutex manual goroutine management - Better pattern in Go 1.25+
[DEPRECATION] MD5 password authentication - Removal in PostgreSQL 19
```

### Issue Body Template
```markdown
## Feature Being Deprecated
[What is being deprecated]

## Announced In Version
[Version X.Y.Z]

## Removal Timeline
- Deprecated: Version A.B
- Final warning: Version C.D
- Removed: Version E.F (estimate)

## Current Usage
[Where is this used in codebase]

## Replacement Pattern
[What should be used instead]

## Migration Effort
[Estimated lines of code to change, number of files, etc.]

## Priority
[Critical/High/Medium/Low]

## Acceptance Criteria
- [ ] All usages identified
- [ ] Replacement pattern documented
- [ ] Code refactored
- [ ] Tests updated
- [ ] No warnings in builds
```

## Monitoring Tools

### Automated Checks
```bash
# Go
go vet ./...                    # Checks for deprecated patterns
golangci-lint run ./...         # Lints include deprecation checks

# Node/npm
npm audit                       # Security & deprecation warnings
npm outdated                    # Shows outdated packages

# Python (if used)
pip-audit                       # Audits for known vulnerabilities
```

### Manual Monitoring
1. **Release notes:** Read each major version's release notes
2. **Snyk:** Run weekly `snyk test` for vulnerability/deprecation alerts
3. **Changelogs:** Subscribe to important project changelogs
4. **Migration guides:** Document custom deprecations in code comments

## Ownership

- **Backend Deprecations:** Go team lead
- **Frontend Deprecations:** Frontend team lead
- **Database Deprecations:** Database/DevOps engineer
- **Security Deprecations:** ALL (immediate escalation)

## Never Ignore

These are NEVER acceptable reasons to defer deprecation fixes:

❌ "It still works"
❌ "No time this sprint"
❌ "It's not breaking yet"
❌ "Other teams do it"
❌ "We'll fix it next year"

✅ These ARE acceptable reasons to prioritize:

✅ "Security deprecation"
✅ "Breaking change coming in 1-2 versions"
✅ "Performance improvement in new pattern"
✅ "New feature requires migration"

## Example: Real Deprecation Fix

**Go 1.21's `context.WithoutCancel()` Introduction**

```go
// BEFORE: Go 1.20 and earlier
func processRequest(parentCtx context.Context) {
    ctx := context.Background() // ❌ Deprecated pattern (Go 1.21+)
    // Process without parent cancellation
}

// AFTER: Go 1.21+ (what we migrated to)
func processRequest(parentCtx context.Context) {
    ctx := context.WithoutCancel(parentCtx) // ✅ Modern pattern
    // Process without parent cancellation, but inherits values
}
```

**Our Process:**
1. ✅ Read Go 1.21 release notes → saw new feature
2. ✅ Created issue: `[DEPRECATION] context.Background() when WithoutCancel available`
3. ✅ Identified 7 usages across codebase
4. ✅ Migrated in next sprint (2 hours)
5. ✅ Tests confirmed behavior unchanged
6. ✅ Merged in PR

Result: **Cleaner code, better semantics, more efficient, future-proofed**

---

This proactive approach keeps us competitive and maintainable.
