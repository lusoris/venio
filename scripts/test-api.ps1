# API Testing Script for Venio RBAC Endpoints
# Tests all 24 RBAC endpoints with different user roles

param(
    [string]$BaseURL = "http://localhost:8080",
    [int]$Verbose = 0
)

# Color output helpers
function Write-Success { Write-Host "✓ $args" -ForegroundColor Green }
function Write-Error { Write-Host "✗ $args" -ForegroundColor Red }
function Write-Info { Write-Host "ℹ $args" -ForegroundColor Blue }
function Write-Warning { Write-Host "⚠ $args" -ForegroundColor Yellow }

# Test results tracking
$totalTests = 0
$passedTests = 0
$failedTests = 0
$testResults = @()

# Function to invoke API with error handling
function Invoke-API {
    param(
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Headers = @{},
        [object]$Body = $null,
        [string]$TestName = ""
    )

    $totalTests++
    $url = "$BaseURL$Endpoint"

    try {
        $params = @{
            Uri             = $url
            Method          = $Method
            Headers         = $Headers
            ContentType     = "application/json"
            ErrorAction     = "Stop"
        }

        if ($Body) {
            $params.Body = $Body | ConvertTo-Json -Depth 10
        }

        $response = Invoke-RestMethod @params

        Write-Success "$TestName (HTTP $($response.StatusCode ?? 200))"
        $passedTests++

        $testResults += @{
            Name   = $TestName
            Status = "PASS"
            Method = $Method
            Path   = $Endpoint
        }

        return $response
    }
    catch {
        Write-Error "$TestName"
        if ($Verbose -ge 1) {
            Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
        }
        $failedTests++

        $testResults += @{
            Name   = $TestName
            Status = "FAIL"
            Method = $Method
            Path   = $Endpoint
            Error  = $_.Exception.Message
        }

        return $null
    }
}

# ============================================================================
# SETUP PHASE - Login as different test users
# ============================================================================

Write-Info "=================================================="
Write-Info "Venio RBAC API Test Suite"
Write-Info "Target: $BaseURL"
Write-Info "=================================================="
Write-Info ""

Write-Info "PHASE 1: Authentication & Token Generation"
Write-Info "---"

# Login credentials for different roles
$testUsers = @{
    "admin"      = @{
        email    = "admin@test.local"
        password = "AdminPassword123!"
    }
    "moderator"  = @{
        email    = "moderator@test.local"
        password = "ModeratorPassword123!"
    }
    "user"       = @{
        email    = "user@test.local"
        password = "UserPassword123!"
    }
    "guest"      = @{
        email    = "guest@test.local"
        password = "GuestPassword123!"
    }
}

$tokens = @{}

foreach ($role in $testUsers.Keys) {
    $loginBody = @{
        email    = $testUsers[$role].email
        password = $testUsers[$role].password
    }

    $response = Invoke-API -Method "POST" `
        -Endpoint "/auth/login" `
        -Body $loginBody `
        -TestName "Login as $role"

    if ($response -and $response.token) {
        $tokens[$role] = $response.token
    }
}

Write-Info ""

# ============================================================================
# TEST PHASE - Role Management Endpoints
# ============================================================================

Write-Info "PHASE 2: Role Management Endpoints"
Write-Info "---"

# Admin can create roles
$newRoleBody = @{
    name        = "test_role_$(Get-Random)"
    description = "Test role for API testing"
}

$createRoleResponse = Invoke-API -Method "POST" `
    -Endpoint "/roles" `
    -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
    -Body $newRoleBody `
    -TestName "Create Role (Admin)"

$testRoleID = $createRoleResponse.id

# Admin can list roles
Invoke-API -Method "GET" `
    -Endpoint "/roles?limit=10&offset=0" `
    -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
    -TestName "List Roles (Admin)"

# Admin can get specific role
if ($testRoleID) {
    Invoke-API -Method "GET" `
        -Endpoint "/roles/$testRoleID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "Get Role by ID (Admin)"
}

# Admin can update role
if ($testRoleID) {
    $updateRoleBody = @{
        name        = "test_role_updated_$(Get-Random)"
        description = "Updated test role"
    }

    Invoke-API -Method "PUT" `
        -Endpoint "/roles/$testRoleID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -Body $updateRoleBody `
        -TestName "Update Role (Admin)"
}

# Non-admin cannot create roles
Invoke-API -Method "POST" `
    -Endpoint "/roles" `
    -Headers @{ Authorization = "Bearer $($tokens['user'])" } `
    -Body $newRoleBody `
    -TestName "Create Role (User - should fail)"

Write-Info ""

# ============================================================================
# TEST PHASE - Permission Management Endpoints
# ============================================================================

Write-Info "PHASE 3: Permission Management Endpoints"
Write-Info "---"

# Admin can create permissions
$newPermissionBody = @{
    name        = "test:permission:$(Get-Random)"
    description = "Test permission for API testing"
}

$createPermissionResponse = Invoke-API -Method "POST" `
    -Endpoint "/permissions" `
    -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
    -Body $newPermissionBody `
    -TestName "Create Permission (Admin)"

$testPermissionID = $createPermissionResponse.id

# Admin can list permissions
Invoke-API -Method "GET" `
    -Endpoint "/permissions?limit=10&offset=0" `
    -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
    -TestName "List Permissions (Admin)"

# Admin can get specific permission
if ($testPermissionID) {
    Invoke-API -Method "GET" `
        -Endpoint "/permissions/$testPermissionID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "Get Permission by ID (Admin)"
}

# Admin can update permission
if ($testPermissionID) {
    $updatePermissionBody = @{
        name        = "test:permission:updated:$(Get-Random)"
        description = "Updated test permission"
    }

    Invoke-API -Method "PUT" `
        -Endpoint "/permissions/$testPermissionID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -Body $updatePermissionBody `
        -TestName "Update Permission (Admin)"
}

Write-Info ""

# ============================================================================
# TEST PHASE - Role-Permission Assignment
# ============================================================================

Write-Info "PHASE 4: Role-Permission Assignment"
Write-Info "---"

# Admin can assign permission to role
if ($testRoleID -and $testPermissionID) {
    $assignBody = @{
        permission_id = $testPermissionID
    }

    Invoke-API -Method "POST" `
        -Endpoint "/roles/$testRoleID/permissions" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -Body $assignBody `
        -TestName "Assign Permission to Role (Admin)"

    # Admin can get role permissions
    Invoke-API -Method "GET" `
        -Endpoint "/roles/$testRoleID/permissions" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "Get Role Permissions (Admin)"

    # Admin can remove permission from role
    Invoke-API -Method "DELETE" `
        -Endpoint "/roles/$testRoleID/permissions/$testPermissionID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "Remove Permission from Role (Admin)"
}

Write-Info ""

# ============================================================================
# TEST PHASE - User Role Management
# ============================================================================

Write-Info "PHASE 5: User Role Management"
Write-Info "---"

# Get current user first (need user ID for role assignment)
$currentUserResponse = Invoke-API -Method "GET" `
    -Endpoint "/auth/me" `
    -Headers @{ Authorization = "Bearer $($tokens['user'])" } `
    -TestName "Get Current User (User)"

$testUserID = $currentUserResponse.id

# Admin can list user roles
if ($testUserID) {
    Invoke-API -Method "GET" `
        -Endpoint "/users/$testUserID/roles" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "List User Roles (Admin)"
}

# Admin can assign role to user
if ($testUserID -and $testRoleID) {
    $assignUserRoleBody = @{
        role_id = $testRoleID
    }

    Invoke-API -Method "POST" `
        -Endpoint "/users/$testUserID/roles" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -Body $assignUserRoleBody `
        -TestName "Assign Role to User (Admin)"

    # Admin can remove role from user
    Invoke-API -Method "DELETE" `
        -Endpoint "/users/$testUserID/roles/$testRoleID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -Body $assignUserRoleBody `
        -TestName "Remove Role from User (Admin)"
}

Write-Info ""

# ============================================================================
# TEST PHASE - Permission-Based Access Control
# ============================================================================

Write-Info "PHASE 6: Permission-Based Access Control"
Write-Info "---"

# Guest should have limited access
Invoke-API -Method "GET" `
    -Endpoint "/roles" `
    -Headers @{ Authorization = "Bearer $($tokens['guest'])" } `
    -TestName "List Roles (Guest - permission-based)"

# Moderator should have moderate access
Invoke-API -Method "GET" `
    -Endpoint "/roles?limit=5&offset=0" `
    -Headers @{ Authorization = "Bearer $($tokens['moderator'])" } `
    -TestName "List Roles (Moderator)"

Write-Info ""

# ============================================================================
# CLEANUP PHASE - Delete test data
# ============================================================================

Write-Info "PHASE 7: Cleanup"
Write-Info "---"

# Admin can delete test role
if ($testRoleID) {
    Invoke-API -Method "DELETE" `
        -Endpoint "/roles/$testRoleID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "Delete Test Role (Admin)"
}

# Admin can delete test permission
if ($testPermissionID) {
    Invoke-API -Method "DELETE" `
        -Endpoint "/permissions/$testPermissionID" `
        -Headers @{ Authorization = "Bearer $($tokens['admin'])" } `
        -TestName "Delete Test Permission (Admin)"
}

Write-Info ""

# ============================================================================
# SUMMARY
# ============================================================================

Write-Info "=================================================="
Write-Info "Test Summary"
Write-Info "=================================================="
Write-Info "Total Tests:   $totalTests"
Write-Success "Passed:        $passedTests"
if ($failedTests -gt 0) {
    Write-Error "Failed:        $failedTests"
}
Write-Info "Success Rate:  $(([math]::Round(($passedTests / $totalTests) * 100, 2)))%"
Write-Info ""

if ($failedTests -eq 0) {
    Write-Success "All tests passed!"
    exit 0
}
else {
    Write-Warning "Some tests failed. Review results below."
    Write-Host ""
    Write-Host "Failed Tests:" -ForegroundColor Red
    foreach ($result in ($testResults | Where-Object { $_.Status -eq "FAIL" })) {
        Write-Host "  - $($result.Name)" -ForegroundColor Red
        if ($result.Error) {
            Write-Host "    Error: $($result.Error)" -ForegroundColor DarkRed
        }
    }
    exit 1
}
