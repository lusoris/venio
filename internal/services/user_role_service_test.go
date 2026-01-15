package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lusoris/venio/internal/models"
)

// MockUserRoleRepositoryForService is a mock implementation
type MockUserRoleRepositoryForService struct {
	mock.Mock
}

func (m *MockUserRoleRepositoryForService) GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockUserRoleRepositoryForService) HasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	args := m.Called(ctx, userID, roleName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRoleRepositoryForService) HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	args := m.Called(ctx, userID, permissionName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRoleRepositoryForService) AssignRole(ctx context.Context, userID, roleID int64) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockUserRoleRepositoryForService) RemoveRole(ctx context.Context, userID, roleID int64) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// TestUserRoleService_GetUserRoles_Success tests successful role retrieval
func TestUserRoleService_GetUserRoles_Success(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	expectedRoles := []models.Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "user"},
	}

	mockRepo.On("GetUserRoles", mock.Anything, int64(1)).Return(expectedRoles, nil)

	roles, err := service.GetUserRoles(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 2)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_GetUserRoles_NotFound tests when user has no roles
func TestUserRoleService_GetUserRoles_NotFound(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("GetUserRoles", mock.Anything, int64(999)).Return([]models.Role{}, nil)

	roles, err := service.GetUserRoles(context.Background(), 999)

	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 0)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_HasRole_True tests when user has role
func TestUserRoleService_HasRole_True(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("HasRole", mock.Anything, int64(1), "admin").Return(true, nil)

	has, err := service.HasRole(context.Background(), 1, "admin")

	assert.NoError(t, err)
	assert.True(t, has)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_HasRole_False tests when user doesn't have role
func TestUserRoleService_HasRole_False(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("HasRole", mock.Anything, int64(1), "admin").Return(false, nil)

	has, err := service.HasRole(context.Background(), 1, "admin")

	assert.NoError(t, err)
	assert.False(t, has)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_HasPermission_True tests when user has permission
func TestUserRoleService_HasPermission_True(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("HasPermission", mock.Anything, int64(1), "users.read").Return(true, nil)

	has, err := service.HasPermission(context.Background(), 1, "users.read")

	assert.NoError(t, err)
	assert.True(t, has)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_HasPermission_False tests when user doesn't have permission
func TestUserRoleService_HasPermission_False(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("HasPermission", mock.Anything, int64(1), "users.delete").Return(false, nil)

	has, err := service.HasPermission(context.Background(), 1, "users.delete")

	assert.NoError(t, err)
	assert.False(t, has)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_AssignRole_Success tests successful role assignment
func TestUserRoleService_AssignRole_Success(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("AssignRole", mock.Anything, int64(1), int64(2)).Return(nil)

	err := service.AssignRole(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_AssignRole_Error tests role assignment failure
func TestUserRoleService_AssignRole_Error(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("AssignRole", mock.Anything, int64(1), int64(999)).Return(errors.New("role not found"))

	err := service.AssignRole(context.Background(), 1, 999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_RemoveRole_Success tests successful role removal
func TestUserRoleService_RemoveRole_Success(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("RemoveRole", mock.Anything, int64(1), int64(2)).Return(nil)

	err := service.RemoveRole(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRoleService_RemoveRole_Error tests role removal failure
func TestUserRoleService_RemoveRole_Error(t *testing.T) {
	mockRepo := new(MockUserRoleRepositoryForService)
	service := NewUserRoleService(mockRepo)

	mockRepo.On("RemoveRole", mock.Anything, int64(1), int64(999)).Return(errors.New("role assignment not found"))

	err := service.RemoveRole(context.Background(), 1, 999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
