package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lusoris/venio/internal/models"
)

// MockPermissionRepository is a mock implementation of PermissionRepository
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetByID(ctx context.Context, id int64) (*models.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Create(ctx context.Context, req *models.CreatePermissionRequest) (*models.Permission, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Update(ctx context.Context, id int64, req *models.UpdatePermissionRequest) (*models.Permission, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) List(ctx context.Context, limit, offset int) ([]models.Permission, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Permission, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) AssignToRole(ctx context.Context, roleID, permissionID int64) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockPermissionRepository) RemoveFromRole(ctx context.Context, roleID, permissionID int64) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

// TestPermissionService_GetByID_Success tests successful permission retrieval by ID
func TestPermissionService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	expectedPerm := &models.Permission{
		ID:          1,
		Name:        "users.read",
		Description: "Read users",
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedPerm, nil)

	perm, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, perm)
	assert.Equal(t, expectedPerm.ID, perm.ID)
	mockRepo.AssertExpectations(t)
}

// TestPermissionService_GetByID_NotFound tests retrieval when permission doesn't exist
func TestPermissionService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	perm, err := service.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, perm)
	mockRepo.AssertExpectations(t)
}

// TestPermissionService_GetByName_Success tests successful permission retrieval by name
func TestPermissionService_GetByName_Success(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	expectedPerm := &models.Permission{
		ID:   1,
		Name: "users.read",
	}

	mockRepo.On("GetByName", mock.Anything, "users.read").Return(expectedPerm, nil)

	perm, err := service.GetByName(context.Background(), "users.read")

	assert.NoError(t, err)
	assert.NotNil(t, perm)
	assert.Equal(t, expectedPerm.Name, perm.Name)
	mockRepo.AssertExpectations(t)
}

// TestPermissionService_Create_Success tests successful permission creation
func TestPermissionService_Create_Success(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	req := models.CreatePermissionRequest{
		Name:        "posts.write",
		Description: "Write posts",
	}

	mockRepo.On("GetByName", mock.Anything, req.Name).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, &req).Return(&models.Permission{
		ID:          1,
		Name:        req.Name,
		Description: req.Description,
	}, nil)

	perm, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, perm)
	assert.Equal(t, req.Name, perm.Name)
	mockRepo.AssertExpectations(t)
}

// TestPermissionService_Create_DuplicateName tests creation with existing name
func TestPermissionService_Create_DuplicateName(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	req := models.CreatePermissionRequest{
		Name:        "users.read",
		Description: "Read users",
	}

	existingPerm := &models.Permission{
		ID:   1,
		Name: "users.read",
	}

	mockRepo.On("GetByName", mock.Anything, req.Name).Return(existingPerm, nil)

	perm, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, perm)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

// TestPermissionService_List_Success tests successful permission listing
func TestPermissionService_List_Success(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	expectedPerms := []models.Permission{
		{ID: 1, Name: "users.read"},
		{ID: 2, Name: "users.write"},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(expectedPerms, int64(2), nil)

	perms, total, err := service.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, perms)
	assert.Len(t, perms, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// TestPermissionService_Delete_Success tests successful permission deletion
func TestPermissionService_Delete_Success(t *testing.T) {
	mockRepo := new(MockPermissionRepository)
	service := NewPermissionService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&models.Permission{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := service.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
