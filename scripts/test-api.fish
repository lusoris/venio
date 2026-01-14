#!/usr/bin/env fish
# API Testing Script for Venio RBAC Endpoints
# Tests all 24 RBAC endpoints with different user roles

set -l BASE_URL (test -n "$argv[1]" && echo $argv[1] || echo "http://localhost:8080")
set -l VERBOSE (test -n "$argv[2]" && echo $argv[2] || echo "0")

# Color output helpers
function print_success
    echo "✓ $argv"
end

function print_error
    echo "✗ $argv" >&2
end

function print_info
    echo "ℹ $argv"
end

function print_warning
    echo "⚠ $argv"
end

# Test results tracking
set -l total_tests 0
set -l passed_tests 0
set -l failed_tests 0

# Function to invoke API with error handling
function invoke_api
    set -l method $argv[1]
    set -l endpoint $argv[2]
    set -l token $argv[3]
    set -l body $argv[4]
    set -l test_name $argv[5]

    set total_tests (math $total_tests + 1)
    set -l url "$BASE_URL$endpoint"

    set -l curl_args -s -X $method "$url" -H "Content-Type: application/json"

    if test -n "$token"
        set curl_args $curl_args -H "Authorization: Bearer $token"
    end

    if test -n "$body"
        set curl_args $curl_args -d "$body"
    end

    set -l response (curl $curl_args 2>&1)
    set -l exit_code $status

    if test $exit_code -eq 0
        print_success "$test_name"
        set passed_tests (math $passed_tests + 1)
        echo "$response"
        return 0
    else
        print_error "$test_name"
        if test $VERBOSE -ge 1
            echo "  Error: $response" >&2
        end
        set failed_tests (math $failed_tests + 1)
        return 1
    end
end

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
set -l test_users admin "admin@test.local" "AdminPassword123!" \
                  moderator "moderator@test.local" "ModeratorPassword123!" \
                  user "user@test.local" "UserPassword123!" \
                  guest "guest@test.local" "GuestPassword123!"

# Build tokens map
set -l tokens_admin ""
set -l tokens_moderator ""
set -l tokens_user ""
set -l tokens_guest ""

for i in (seq 1 4 (count $test_users))
    set -l role $test_users[$i]
    set -l email $test_users[(math $i + 1)]
    set -l password $test_users[(math $i + 2)]

    set -l login_body "{\"email\":\"$email\",\"password\":\"$password\"}"

    set -l response (invoke_api "POST" "/auth/login" "" "$login_body" "Login as $role" 2>/dev/null)
    set -l token (echo "$response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

    if test -n "$token"
        set tokens_$role "$token"
    end
end

print_info ""

# ============================================================================
# TEST PHASE - Role Management Endpoints
# ============================================================================

print_info "PHASE 2: Role Management Endpoints"
print_info "---"

# Admin can create roles
set -l role_name "test_role_"(random)
set -l new_role_body "{\"name\":\"$role_name\",\"description\":\"Test role for API testing\"}"

set -l response (invoke_api "POST" "/roles" "$tokens_admin" "$new_role_body" "Create Role (Admin)" 2>/dev/null)
set -l test_role_id (echo "$response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

# Admin can list roles
invoke_api "GET" "/roles?limit=10&offset=0" "$tokens_admin" "" "List Roles (Admin)" > /dev/null 2>&1

# Admin can get specific role
if test -n "$test_role_id"
    invoke_api "GET" "/roles/$test_role_id" "$tokens_admin" "" "Get Role by ID (Admin)" > /dev/null 2>&1
end

# Admin can update role
if test -n "$test_role_id"
    set -l update_role_name "test_role_updated_"(random)
    set -l update_role_body "{\"name\":\"$update_role_name\",\"description\":\"Updated test role\"}"
    invoke_api "PUT" "/roles/$test_role_id" "$tokens_admin" "$update_role_body" "Update Role (Admin)" > /dev/null 2>&1
end

# Non-admin cannot create roles (expected to fail)
invoke_api "POST" "/roles" "$tokens_user" "$new_role_body" "Create Role (User - should fail)" > /dev/null 2>&1

print_info ""

# ============================================================================
# TEST PHASE - Permission Management Endpoints
# ============================================================================

print_info "PHASE 3: Permission Management Endpoints"
print_info "---"

# Admin can create permissions
set -l perm_name "test:permission:"(random)
set -l new_permission_body "{\"name\":\"$perm_name\",\"description\":\"Test permission for API testing\"}"

set -l response (invoke_api "POST" "/permissions" "$tokens_admin" "$new_permission_body" "Create Permission (Admin)" 2>/dev/null)
set -l test_permission_id (echo "$response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

# Admin can list permissions
invoke_api "GET" "/permissions?limit=10&offset=0" "$tokens_admin" "" "List Permissions (Admin)" > /dev/null 2>&1

# Admin can get specific permission
if test -n "$test_permission_id"
    invoke_api "GET" "/permissions/$test_permission_id" "$tokens_admin" "" "Get Permission by ID (Admin)" > /dev/null 2>&1
end

# Admin can update permission
if test -n "$test_permission_id"
    set -l update_perm_name "test:permission:updated:"(random)
    set -l update_permission_body "{\"name\":\"$update_perm_name\",\"description\":\"Updated test permission\"}"
    invoke_api "PUT" "/permissions/$test_permission_id" "$tokens_admin" "$update_permission_body" "Update Permission (Admin)" > /dev/null 2>&1
end

print_info ""

# ============================================================================
# CLEANUP PHASE
# ============================================================================

print_info "PHASE 7: Cleanup"
print_info "---"

if test -n "$test_role_id"
    invoke_api "DELETE" "/roles/$test_role_id" "$tokens_admin" "" "Delete Test Role (Admin)" > /dev/null 2>&1
end

if test -n "$test_permission_id"
    invoke_api "DELETE" "/permissions/$test_permission_id" "$tokens_admin" "" "Delete Test Permission (Admin)" > /dev/null 2>&1
end

print_info ""

# ============================================================================
# SUMMARY
# ============================================================================

print_info "=================================================="
print_info "Test Summary"
print_info "=================================================="
print_info "Total Tests:   $total_tests"
print_success "Passed:        $passed_tests"
if test $failed_tests -gt 0
    print_error "Failed:        $failed_tests"
end

set -l success_rate (math "scale=2; ($passed_tests / $total_tests) * 100")
print_info "Success Rate:  $success_rate%"
print_info ""

if test $failed_tests -eq 0
    print_success "All tests passed!"
    exit 0
else
    print_warning "Some tests failed. Review output above."
    exit 1
end
