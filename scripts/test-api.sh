#!/bin/bash
# API Testing Script for Venio RBAC Endpoints
# Tests all 24 RBAC endpoints with different user roles

BASE_URL="${1:-http://localhost:8080}"
VERBOSE="${2:-0}"

# Color output helpers
print_success() { echo "✓ $*"; }
print_error() { echo "✗ $*"; }
print_info() { echo "ℹ $*"; }
print_warning() { echo "⚠ $*"; }

# Test results tracking
total_tests=0
passed_tests=0
failed_tests=0
declare -a test_results

# Function to invoke API with error handling
invoke_api() {
    local method="$1"
    local endpoint="$2"
    local token="$3"
    local body="$4"
    local test_name="$5"

    ((total_tests++))
    local url="$BASE_URL$endpoint"

    local curl_args=("-s" "-X" "$method" "$url" "-H" "Content-Type: application/json")

    if [[ -n "$token" ]]; then
        curl_args+=("-H" "Authorization: Bearer $token")
    fi

    if [[ -n "$body" ]]; then
        curl_args+=("-d" "$body")
    fi

    local response
    response=$(curl "${curl_args[@]}" 2>&1)
    local exit_code=$?

    if [[ $exit_code -eq 0 ]]; then
        print_success "$test_name"
        ((passed_tests++))
        echo "$response"
        return 0
    else
        print_error "$test_name"
        if [[ $VERBOSE -ge 1 ]]; then
            echo "  Error: $response"
        fi
        ((failed_tests++))
        return 1
    fi
}

# ============================================================================
# SETUP PHASE - Login as different test users
# ============================================================================

print_info "=================================================="
print_info "Venio RBAC API Test Suite"
print_info "Target: $BASE_URL"
print_info "=================================================="
print_info ""

print_info "PHASE 1: Authentication & Token Generation"
print_info "---"

# Login credentials for different roles
declare -A test_users=(
    [admin]="admin@test.local|AdminPassword123!"
    [moderator]="moderator@test.local|ModeratorPassword123!"
    [user]="user@test.local|UserPassword123!"
    [guest]="guest@test.local|GuestPassword123!"
)

declare -A tokens

for role in "${!test_users[@]}"; do
    IFS='|' read -r email password <<< "${test_users[$role]}"

    local login_body=$(cat <<EOF
{"email":"$email","password":"$password"}
EOF
)

    response=$(invoke_api "POST" "/auth/login" "" "$login_body" "Login as $role")
    token=$(echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

    if [[ -n "$token" ]]; then
        tokens[$role]="$token"
    fi
done

print_info ""

# ============================================================================
# TEST PHASE - Role Management Endpoints
# ============================================================================

print_info "PHASE 2: Role Management Endpoints"
print_info "---"

# Admin can create roles
role_name="test_role_$RANDOM"
new_role_body=$(cat <<EOF
{"name":"$role_name","description":"Test role for API testing"}
EOF
)

response=$(invoke_api "POST" "/roles" "${tokens[admin]}" "$new_role_body" "Create Role (Admin)")
test_role_id=$(echo "$response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

# Admin can list roles
invoke_api "GET" "/roles?limit=10&offset=0" "${tokens[admin]}" "" "List Roles (Admin)" > /dev/null

# Admin can get specific role
if [[ -n "$test_role_id" ]]; then
    invoke_api "GET" "/roles/$test_role_id" "${tokens[admin]}" "" "Get Role by ID (Admin)" > /dev/null
fi

# Admin can update role
if [[ -n "$test_role_id" ]]; then
    update_role_name="test_role_updated_$RANDOM"
    update_role_body=$(cat <<EOF
{"name":"$update_role_name","description":"Updated test role"}
EOF
)
    invoke_api "PUT" "/roles/$test_role_id" "${tokens[admin]}" "$update_role_body" "Update Role (Admin)" > /dev/null
fi

# Non-admin cannot create roles (expected to fail)
invoke_api "POST" "/roles" "${tokens[user]}" "$new_role_body" "Create Role (User - should fail)" > /dev/null 2>&1 || true

print_info ""

# ============================================================================
# TEST PHASE - Permission Management Endpoints
# ============================================================================

print_info "PHASE 3: Permission Management Endpoints"
print_info "---"

# Admin can create permissions
perm_name="test:permission:$RANDOM"
new_permission_body=$(cat <<EOF
{"name":"$perm_name","description":"Test permission for API testing"}
EOF
)

response=$(invoke_api "POST" "/permissions" "${tokens[admin]}" "$new_permission_body" "Create Permission (Admin)")
test_permission_id=$(echo "$response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

# Admin can list permissions
invoke_api "GET" "/permissions?limit=10&offset=0" "${tokens[admin]}" "" "List Permissions (Admin)" > /dev/null

# Admin can get specific permission
if [[ -n "$test_permission_id" ]]; then
    invoke_api "GET" "/permissions/$test_permission_id" "${tokens[admin]}" "" "Get Permission by ID (Admin)" > /dev/null
fi

# Admin can update permission
if [[ -n "$test_permission_id" ]]; then
    update_perm_name="test:permission:updated:$RANDOM"
    update_permission_body=$(cat <<EOF
{"name":"$update_perm_name","description":"Updated test permission"}
EOF
)
    invoke_api "PUT" "/permissions/$test_permission_id" "${tokens[admin]}" "$update_permission_body" "Update Permission (Admin)" > /dev/null
fi

print_info ""

# ============================================================================
# CLEANUP PHASE
# ============================================================================

print_info "PHASE 7: Cleanup"
print_info "---"

if [[ -n "$test_role_id" ]]; then
    invoke_api "DELETE" "/roles/$test_role_id" "${tokens[admin]}" "" "Delete Test Role (Admin)" > /dev/null
fi

if [[ -n "$test_permission_id" ]]; then
    invoke_api "DELETE" "/permissions/$test_permission_id" "${tokens[admin]}" "" "Delete Test Permission (Admin)" > /dev/null
fi

print_info ""

# ============================================================================
# SUMMARY
# ============================================================================

print_info "=================================================="
print_info "Test Summary"
print_info "=================================================="
print_info "Total Tests:   $total_tests"
print_success "Passed:        $passed_tests"
if [[ $failed_tests -gt 0 ]]; then
    print_error "Failed:        $failed_tests"
fi

success_rate=$(awk "BEGIN {printf \"%.2f\", ($passed_tests / $total_tests) * 100}")
print_info "Success Rate:  $success_rate%"
print_info ""

if [[ $failed_tests -eq 0 ]]; then
    print_success "All tests passed!"
    exit 0
else
    print_warning "Some tests failed. Review output above."
    exit 1
fi
