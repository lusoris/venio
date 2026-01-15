package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lusoris/venio/internal/models"
)

// MockAuthServiceForHandler is a mock implementation of AuthService for testing handlers
type MockAuthServiceForHandler struct {
	mock.Mock
}

func (m *MockAuthServiceForHandler) Login(ctx context.Context, email, password string) (string, string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthServiceForHandler) ValidateToken(token string) (*models.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenClaims), args.Error(1)
}

func (m *MockAuthServiceForHandler) RefreshToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

// MockUserServiceForHandler is a mock implementation of UserService for testing handlers
type MockUserServiceForHandler struct {
	mock.Mock
}

func (m *MockUserServiceForHandler) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) GetUser(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/register", handler.Register)

	requestBody := models.CreateUserRequest{
		Email:     "newuser@example.com",
		Username:  "newuser",
		FirstName: "New",
		LastName:  "User",
		Password:  "password123",
	}

	body, _ := json.Marshal(requestBody)

	newUser := &models.User{
		ID:       1,
		Email:    "newuser@example.com",
		Username: "newuser",
		IsActive: true,
	}

	mockUserService.On("Register", mock.Anything, mock.MatchedBy(func(req *models.CreateUserRequest) bool {
		return req.Email == "newuser@example.com"
	})).Return(newUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "newuser@example.com", response.Email)
}

func TestAuthHandler_Register_InvalidEmail(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/register", handler.Register)

	requestBody := map[string]interface{}{
		"email":      "invalid-email",
		"username":   "testuser",
		"first_name": "Test",
		"last_name":  "User",
		"password":   "password123",
	}

	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/login", handler.Login)

	requestBody := models.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(requestBody)

	// Mock both calls
	mockAuthService.On("Login", mock.Anything, "user@example.com", "password123").
		Return("access.token.here", "refresh.token.here", nil)

	mockUserService.On("GetUserByEmail", mock.Anything, "user@example.com").Return(&models.User{
		ID:       1,
		Email:    "user@example.com",
		Username: "testuser",
		IsActive: true,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["access_token"])
	assert.NotNil(t, response["refresh_token"])
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/login", handler.Login)

	requestBody := models.LoginRequest{
		Email:    "user@example.com",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(requestBody)

	mockAuthService.On("Login", mock.Anything, "user@example.com", "wrongpassword").
		Return("", "", assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/refresh", handler.RefreshToken)

	requestBody := RefreshTokenRequest{
		RefreshToken: "refresh.token.here",
	}

	body, _ := json.Marshal(requestBody)

	mockAuthService.On("RefreshToken", mock.Anything, "refresh.token.here").
		Return("new.access.token", nil)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["access_token"])
}

func TestAuthHandler_Login_MissingEmail(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/login", handler.Login)

	requestBody := map[string]string{
		"password": "password123",
	}

	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
