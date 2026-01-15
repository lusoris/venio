package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lusoris/venio/internal/models"
)

// MockRoleRepositoryForTest is a mock implementation of RoleRepository
type MockRoleRepositoryForTest struct {
	mock.Mock
}

func (m *MockRoleRepositoryForTest) GetByID(ctx context.Context, id int64) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepositoryForTest) GetByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepositoryForTest) Create(ctx context.Context, req *models.CreateRoleRequest) (*models.Role, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepositoryForTest) Update(ctx context.Context, id int64, req *models.UpdateRoleRequest) (*models.Role, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepositoryForTest) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepositoryForTest) List(ctx context.Context, limit, offset int) ([]models.Role, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepositoryForTest) GetPermissions(ctx context.Context, roleID int64) ([]models.Permission, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}

// TestRoleService_GetByID_Success tests successful role retrieval by ID
func TestRoleService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	expectedRole := &models.Role{
		ID:          1,
		Name:        "admin",
		Description: "Administrator role",
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedRole, nil)

	role, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, expectedRole.ID, role.ID)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_GetByID_InvalidID tests retrieval with invalid ID
func TestRoleService_GetByID_InvalidID(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	role, err := service.GetByID(context.Background(), 0)

	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Contains(t, err.Error(), "invalid")
	mockRepo.AssertExpectations(t)
}

// TestRoleService_GetByID_NotFound tests retrieval when role doesn't exist
func TestRoleService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	role, err := service.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, role)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_GetByName_Success tests successful role retrieval by name
func TestRoleService_GetByName_Success(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	expectedRole := &models.Role{
		ID:   1,
		Name: "admin",
	}

	mockRepo.On("GetByName", mock.Anything, "admin").Return(expectedRole, nil)

	role, err := service.GetByName(context.Background(), "admin")

	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, expectedRole.Name, role.Name)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_GetByName_EmptyName tests retrieval with empty name
func TestRoleService_GetByName_EmptyName(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	role, err := service.GetByName(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, role)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_Create_Success tests successful role creation
func TestRoleService_Create_Success(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	req := models.CreateRoleRequest{
		Name:        "moderator",
		Description: "Moderator role",
	}

	mockRepo.On("GetByName", mock.Anything, req.Name).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, &req).Return(&models.Role{
		ID:          1,
		Name:        req.Name,
		Description: req.Description,
	}, nil)

	role, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, req.Name, role.Name)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_Create_DuplicateName tests creation with existing name
func TestRoleService_Create_DuplicateName(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	req := models.CreateRoleRequest{
		Name:        "admin",
		Description: "Admin role",
	}

	existingRole := &models.Role{
		ID:   1,
		Name: "admin",
	}

	mockRepo.On("GetByName", mock.Anything, req.Name).Return(existingRole, nil)

	role, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

// TestRoleService_Delete_Success tests successful role deletion
func TestRoleService_Delete_Success(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&models.Role{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := service.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_Delete_NotFound tests deletion when role doesn't exist
func TestRoleService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, nil)

	err := service.Delete(context.Background(), 999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_List_Success tests successful role listing
func TestRoleService_List_Success(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	expectedRoles := []models.Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "user"},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(expectedRoles, int64(2), nil)

	roles, total, err := service.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// TestRoleService_GetPermissions_Success tests successful permission retrieval for role
func TestRoleService_GetPermissions_Success(t *testing.T) {
	mockRepo := new(MockRoleRepositoryForTest)
	service := NewRoleService(mockRepo)

	expectedPerms := []models.Permission{
		{ID: 1, Name: "users.read"},
		{ID: 2, Name: "users.write"},
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&models.Role{ID: 1}, nil)
	mockRepo.On("GetPermissions", mock.Anything, int64(1)).Return(expectedPerms, nil)

	perms, err := service.GetPermissions(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, perms)
	assert.Len(t, perms, 2)
	mockRepo.AssertExpectations(t)
}
