# API Testing Guide

Comprehensive testing scripts for Venio's RBAC API endpoints. Tests are available in multiple shells to support different development environments.

## Quick Start

### Prerequisites

- **PowerShell:** Windows (included), or PowerShell Core (all platforms)
- **Bash:** Linux, macOS, or Git Bash on Windows
- **Fish:** Any platform with Fish shell installed
- **All platforms need:** `curl` for HTTP requests

### Running Tests

#### PowerShell (Windows/Core)
```powershell
# Run with defaults (localhost:8080)
.\scripts\test-api.ps1

# Run with custom URL
.\scripts\test-api.ps1 -BaseURL http://localhost:3000 -Verbose 1

# Or via Makefile
make test-api
```

#### Bash (Linux/macOS/Git Bash)
```bash
# Run with defaults
./scripts/test-api.sh

# Run with custom URL and verbose output
./scripts/test-api.sh http://localhost:3000 1

# Via Makefile
make test-api
```

#### Fish Shell (Linux/macOS)
```fish
# Run with defaults
./scripts/test-api.fish

# Run with custom URL
./scripts/test-api.fish http://localhost:3000 1

# Via Makefile (if fish is your shell)
make test-api
```

## Test Suite Structure

The test suite consists of **7 phases**:

### Phase 1: Authentication & Token Generation
- Tests login for all 4 test user roles
- Generates JWT tokens for subsequent tests
- Validates credentials work correctly

**Test Users:**
- `admin@test.local` / `AdminPassword123!` → admin role
- `moderator@test.local` / `ModeratorPassword123!` → moderator role
- `user@test.local` / `UserPassword123!` → user role
- `guest@test.local` / `GuestPassword123!` → guest role

### Phase 2: Role Management Endpoints
- ✓ Create role (admin only)
- ✓ List roles (admin only)
- ✓ Get role by ID (admin only)
- ✓ Update role (admin only)
- ✓ Verify non-admin cannot create roles

### Phase 3: Permission Management Endpoints
- ✓ Create permission (admin only)
- ✓ List permissions (admin only)
- ✓ Get permission by ID (admin only)
- ✓ Update permission (admin only)

### Phase 4: Role-Permission Assignment
- ✓ Assign permission to role
- ✓ List role's permissions
- ✓ Remove permission from role

### Phase 5: User Role Management
- ✓ Get current user info
- ✓ List user's roles
- ✓ Assign role to user
- ✓ Remove role from user

### Phase 6: Permission-Based Access Control
- ✓ Test guest access (limited)
- ✓ Test moderator access (moderate)
- ✓ Verify authorization enforcement

### Phase 7: Cleanup
- ✓ Delete test role
- ✓ Delete test permission
- ✓ Ensure no test data remains

## Database Setup

Before running API tests, seed the database with default roles and permissions:

```bash
# Option 1: Via Makefile
make seed-data

# Option 2: Direct seeder
go run cmd/seeder/main.go

# Option 3: Via migration
docker compose up postgres -d
migrate -path migrations -database "postgres://venio:venio@localhost:5432/venio?sslmode=disable" up
go run cmd/seeder/main.go
```

This creates:
- **4 default roles:** admin, moderator, user, guest
- **16 permissions:** CRUD operations on users, roles, permissions, content, settings, and audit logs
- **4 test users:** One for each role with known credentials

## Customization

### Modify Base URL

All scripts accept the base URL as the first parameter:

```bash
# Bash
./scripts/test-api.sh http://localhost:3000

# Fish
./scripts/test-api.fish http://api.local:8000

# PowerShell
.\scripts\test-api.ps1 -BaseURL http://api.staging.com
```

### Verbose Output

All scripts support verbose output (shows error details):

```bash
# Bash
./scripts/test-api.sh http://localhost:8080 1

# Fish
./scripts/test-api.fish http://localhost:8080 1

# PowerShell
.\scripts\test-api.ps1 -Verbose 1
```

### Modify Test Data

Edit the test credentials in your preferred script:

**PowerShell (`scripts/test-api.ps1`)**
```powershell
$testUsers = @{
    "admin" = @{
        email    = "admin@test.local"
        password = "AdminPassword123!"
    }
    # ... modify as needed
}
```

**Bash (`scripts/test-api.sh`)**
```bash
declare -A test_users=(
    [admin]="admin@test.local|AdminPassword123!"
    # ... modify as needed
)
```

**Fish (`scripts/test-api.fish`)**
```fish
set -l test_users admin "admin@test.local" "AdminPassword123!" \
                  # ... modify as needed
```

## Expected Results

Successful test run output:

```
==================================================
Venio RBAC API Test Suite
Target: http://localhost:8080
==================================================

PHASE 1: Authentication & Token Generation
---
✓ Login as admin
✓ Login as moderator
✓ Login as user
✓ Login as guest

PHASE 2: Role Management Endpoints
---
✓ Create Role (Admin)
✓ List Roles (Admin)
✓ Get Role by ID (Admin)
✓ Update Role (Admin)
✓ Create Role (User - should fail)

... more tests ...

==================================================
Test Summary
==================================================
Total Tests:   45
✓ Passed:        45
Success Rate:  100.00%

✓ All tests passed!
```

## Troubleshooting

### Connection Refused
```
✗ Login as admin
Error: Connection refused
```
**Solution:** Ensure the server is running and accessible at the specified URL.

```bash
curl http://localhost:8080/health  # Check if server is alive
```

### Invalid Credentials
```
✗ Login as admin
Error: Invalid credentials
```
**Solution:** Ensure test users exist in the database. Run seed-data:

```bash
make seed-data
```

### Permission Denied
```bash
Permission denied: ./scripts/test-api.sh
```
**Solution:** Make scripts executable:

```bash
chmod +x scripts/test-api.sh
chmod +x scripts/test-api.fish
```

### Snyk Not Available
In PowerShell pre-push hooks, Snyk is optional:
- If available: security scan runs
- If not installed: message shown, hook continues

Install Snyk:
```bash
npm install -g snyk
snyk auth
```

## Cross-Platform Development

These scripts are designed for developers using Visual Studio Code across Windows and Linux:

1. **Windows:** Use PowerShell (`test-api.ps1`)
2. **Linux:** Use Bash (`test-api.sh`) or Fish (`test-api.fish`)
3. **macOS:** Use Bash (`test-api.sh`) or Fish (`test-api.fish`)

All scripts produce identical output and test coverage. Choose based on your shell preference.

### Git Synchronization Tips

When syncing between Windows and Linux:

```bash
# Set Git to auto-convert line endings
git config core.autocrlf true  # Windows
git config core.autocrlf input # Linux/macOS

# Or globally
git config --global core.autocrlf true  # Windows
git config --global core.autocrlf input # Linux/macOS
```

The scripts use platform-native syntax:
- **PowerShell:** `$variables`, `@{hashtables}`, `Invoke-RestMethod`
- **Bash:** `$variables`, associative arrays, `curl`
- **Fish:** `set -l variables`, `function declarations`, `$expansion`

## Development Workflow

Typical workflow when implementing RBAC features:

```bash
# 1. Start development environment
make dev

# 2. Seed database with test data
make seed-data

# 3. Run API tests
make test-api

# 4. If tests pass, make changes to backend
# 5. Rerun tests
make test-api

# 6. Commit changes (lefthook will run automatically)
git commit -m "feat: implement new RBAC endpoint"
```

## Integration with CI/CD

For GitHub Actions or other CI/CD:

```yaml
# .github/workflows/test.yml
- name: Run API Tests
  run: |
    # Install dependencies
    go mod download

    # Start services
    docker-compose up -d postgres redis

    # Run migrations
    go run cmd/migrate/main.go

    # Seed data
    go run cmd/seeder/main.go

    # Start server (background)
    go run cmd/venio/main.go &
    sleep 2

    # Run tests (bash in CI)
    bash scripts/test-api.sh
```

## Contributing

When adding new RBAC endpoints:

1. Add corresponding tests to all three script versions
2. Keep identical test phases and naming
3. Test on target platform before committing
4. Update this README if behavior changes

---

**Last Updated:** January 2026
**Test Scripts:** PowerShell, Bash, Fish
**Coverage:** 45+ RBAC API tests
**Success Rate:** 100% on passing systems
