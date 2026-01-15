package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/lusoris/venio/internal/config"
	"github.com/lusoris/venio/internal/models"
)

// MockUserService is a mock implementation of UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// MockUserRoleService is a mock implementation of UserRoleService for testing
type MockUserRoleService struct {
	mock.Mock
}

func (m *MockUserRoleService) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRoleService) HasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	args := m.Called(ctx, userID, roleName)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRoleService) HasPermission(ctx context.Context, userID int64, permissionName string) (bool, error) {
	args := m.Called(ctx, userID, permissionName)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRoleService) AssignRole(ctx context.Context, userID, roleID int64) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockUserRoleService) RemoveRole(ctx context.Context, userID, roleID int64) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func TestLogin_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Password: string(hashedPassword),
		IsActive: true,
	}

	mockUserService.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockUserRoleService.On("GetUserRoles", mock.Anything, int64(1)).Return([]string{"user"}, nil)

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	accessToken, refreshToken, err := authService.Login(context.Background(), "test@example.com", password)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockUserService.AssertCalled(t, "GetUserByEmail", mock.Anything, "test@example.com")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Password: string(hashedPassword),
		IsActive: true,
	}

	mockUserService.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockUserRoleService.On("GetUserRoles", mock.Anything, int64(1)).Return([]string{"user"}, nil)

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	_, _, err := authService.Login(context.Background(), "test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_InactiveUser(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Password: string(hashedPassword),
		IsActive: false,
	}

	mockUserService.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockUserRoleService.On("GetUserRoles", mock.Anything, int64(1)).Return([]string{"user"}, nil)

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	_, _, err := authService.Login(context.Background(), "test@example.com", password)

	assert.Error(t, err)
	assert.Equal(t, "user account is inactive", err.Error())
}

func TestLogin_UserNotFound(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	mockUserService.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").
		Return(nil, assert.AnError)

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	_, _, err := authService.Login(context.Background(), "nonexistent@example.com", "password")

	assert.Error(t, err)
}

func TestValidateToken_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Password: string(hashedPassword),
		IsActive: true,
	}

	mockUserService.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockUserRoleService.On("GetUserRoles", mock.Anything, int64(1)).Return([]string{"user"}, nil)

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	accessToken, _, _ := authService.Login(context.Background(), "test@example.com", password)

	claims, err := authService.ValidateToken(accessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	_, err := authService.ValidateToken("invalid.token.string")

	assert.Error(t, err)
}

func TestTokenExpiration(t *testing.T) {
	mockUserService := new(MockUserService)
	mockUserRoleService := new(MockUserRoleService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-must-be-at-least-32-characters-long-ok",
			ExpirationTime: 24 * time.Hour,
		},
	}

	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Password: string(hashedPassword),
		IsActive: true,
	}

	mockUserService.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockUserRoleService.On("GetUserRoles", mock.Anything, int64(1)).Return([]string{"user"}, nil)

	authService := NewDefaultAuthService(mockUserService, mockUserRoleService, cfg)
	accessToken, _, _ := authService.Login(context.Background(), "test@example.com", password)

	claims, err := authService.ValidateToken(accessToken)
	assert.NoError(t, err)

	// Check that token is set to expire in approximately 24 hours
	expiresIn := time.Until(time.Unix(claims.ExpiresAt.Unix(), 0))
	assert.Greater(t, expiresIn, time.Hour*23)
	assert.Less(t, expiresIn, time.Hour*25)
}
