package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lusoris/venio/internal/models"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func testPassword() string {
	return fmt.Sprintf("pw-%d", time.Now().UnixNano())
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// TestRegister_Success tests successful user registration
func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	req := &models.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  testPassword(),
		FirstName: "Test",
		LastName:  "User",
	}

	mockRepo.On("Exists", mock.Anything, req.Email).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(int64(1), nil)

	user, err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Username, user.Username)
	assert.NotEqual(t, req.Password, user.Password) // Password should be hashed
	assert.True(t, user.IsActive)
	mockRepo.AssertExpectations(t)
}

// TestRegister_InvalidEmail tests registration with invalid email
func TestRegister_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	req := &models.CreateUserRequest{
		Email:    "invalid-email",
		Username: "testuser",
		Password: testPassword(),
	}

	user, err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "validation error")
}

// TestRegister_DuplicateEmail tests registration with existing email
func TestRegister_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	req := &models.CreateUserRequest{
		Email:     "existing@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  testPassword(),
	}

	mockRepo.On("Exists", mock.Anything, req.Email).Return(true, nil)

	user, err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already registered")
	mockRepo.AssertExpectations(t)
}

// TestRegister_RepositoryError tests registration when repository fails
func TestRegister_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	req := &models.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  testPassword(),
	}

	mockRepo.On("Exists", mock.Anything, req.Email).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(int64(0), errors.New("database error"))

	user, err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to create user")
	mockRepo.AssertExpectations(t)
}

// TestGetUser_Success tests successful user retrieval
func TestGetUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	expectedUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		IsActive: true,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

	user, err := service.GetUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

// TestGetUser_NotFound tests user retrieval when user doesn't exist
func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("user not found"))

	user, err := service.GetUser(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUpdateUser_Success tests successful user update
func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	existingUser := &models.User{
		ID:       1,
		Email:    "old@example.com",
		Username: "olduser",
		IsActive: true,
	}

	newEmail := "new@example.com"
	newUsername := "newuser"
	updateReq := &models.UpdateUserRequest{
		Email:    &newEmail,
		Username: &newUsername,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	user, err := service.UpdateUser(context.Background(), 1, updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newEmail, user.Email)
	assert.Equal(t, newUsername, user.Username)
	mockRepo.AssertExpectations(t)
}

// TestUpdateUser_NotFound tests update when user doesn't exist
func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	newEmail := "new@example.com"
	updateReq := &models.UpdateUserRequest{
		Email: &newEmail,
	}

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("user not found"))

	user, err := service.UpdateUser(context.Background(), 999, updateReq)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser_Success tests successful user deletion
func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteUser(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser_NotFound tests deletion when user doesn't exist
func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	mockRepo.On("Delete", mock.Anything, int64(999)).Return(errors.New("user not found"))

	err := service.DeleteUser(context.Background(), 999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestListUsers_Success tests successful user listing
func TestListUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	expectedUsers := []*models.User{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(expectedUsers, nil)

	users, err := service.ListUsers(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

// TestListUsers_Empty tests user listing with no results
func TestListUsers_Empty(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	mockRepo.On("List", mock.Anything, 10, 0).Return([]*models.User{}, nil)

	users, err := service.ListUsers(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0)
	mockRepo.AssertExpectations(t)
}

// TestGetUserByEmail_Success tests successful user retrieval by email
func TestGetUserByEmail_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	expectedUser := &models.User{
		ID:    1,
		Email: "test@example.com",
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expectedUser, nil)

	user, err := service.GetUserByEmail(context.Background(), "test@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

// TestGetUserByUsername_Success tests successful user retrieval by username
func TestGetUserByUsername_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewDefaultUserService(mockRepo)

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
	}

	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(expectedUser, nil)

	user, err := service.GetUserByUsername(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Username, user.Username)
	mockRepo.AssertExpectations(t)
}
