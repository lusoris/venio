package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lusoris/venio/internal/models"
	"github.com/lusoris/venio/internal/services"
)

const (
	validVerificationToken   = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	expiredVerificationToken = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
)

func testSecret() string {
	return fmt.Sprintf("s-%d", time.Now().UnixNano())
}

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

func (m *MockAuthServiceForHandler) GenerateEmailVerificationToken(ctx context.Context, userID int64) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthServiceForHandler) VerifyEmail(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthServiceForHandler) ResendVerificationEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
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

func (m *MockUserServiceForHandler) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserServiceForHandler) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserServiceForHandler) GetByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/register", handler.Register)

	password := testSecret()
	requestBody := models.CreateUserRequest{
		Email:     "newuser@example.com",
		Username:  "newuser",
		FirstName: "New",
		LastName:  "User",
		Password:  password,
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
		"password":   testSecret(),
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

	password := testSecret()
	requestBody := models.LoginRequest{
		Email:    "user@example.com",
		Password: password,
	}

	body, _ := json.Marshal(requestBody)

	accessToken := fmt.Sprintf("access-%d", time.Now().UnixNano())
	refreshToken := fmt.Sprintf("refresh-%d", time.Now().UnixNano())

	// Mock both calls
	mockAuthService.On("Login", mock.Anything, "user@example.com", password).
		Return(accessToken, refreshToken, nil)

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

	wrongPassword := testSecret()
	requestBody := models.LoginRequest{
		Email:    "user@example.com",
		Password: wrongPassword,
	}

	body, _ := json.Marshal(requestBody)

	mockAuthService.On("Login", mock.Anything, "user@example.com", wrongPassword).
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

	refreshToken := fmt.Sprintf("refresh-%d", time.Now().UnixNano())
	requestBody := RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	body, _ := json.Marshal(requestBody)

	newAccessToken := fmt.Sprintf("access-%d", time.Now().UnixNano())
	mockAuthService.On("RefreshToken", mock.Anything, refreshToken).
		Return(newAccessToken, nil)

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
		"password": testSecret(),
	}

	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_VerifyEmail_Success(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/verify-email", handler.VerifyEmail)

	body, _ := json.Marshal(VerifyEmailRequest{Token: validVerificationToken})

	mockAuthService.On("VerifyEmail", mock.Anything, validVerificationToken).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_VerifyEmail_ExpiredToken(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/verify-email", handler.VerifyEmail)

	body, _ := json.Marshal(VerifyEmailRequest{Token: expiredVerificationToken})

	mockAuthService.On("VerifyEmail", mock.Anything, expiredVerificationToken).Return(services.ErrVerificationTokenExpired)

	req := httptest.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_ResendVerificationEmail_Success(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/resend-verification", handler.ResendVerificationEmail)

	body, _ := json.Marshal(ResendVerificationRequest{Email: "user@example.com"})

	mockAuthService.On("ResendVerificationEmail", mock.Anything, "user@example.com").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_ResendVerificationEmail_UserNotFound(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/resend-verification", handler.ResendVerificationEmail)

	body, _ := json.Marshal(ResendVerificationRequest{Email: "missing@example.com"})

	mockAuthService.On("ResendVerificationEmail", mock.Anything, "missing@example.com").Return(services.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAuthHandler_ResendVerificationEmail_AlreadyVerified(t *testing.T) {
	mockAuthService := new(MockAuthServiceForHandler)
	mockUserService := new(MockUserServiceForHandler)

	handler := NewAuthHandler(mockAuthService, mockUserService)

	router := gin.New()
	router.POST("/resend-verification", handler.ResendVerificationEmail)

	body, _ := json.Marshal(ResendVerificationRequest{Email: "verified@example.com"})

	mockAuthService.On("ResendVerificationEmail", mock.Anything, "verified@example.com").Return(services.ErrEmailAlreadyVerified)

	req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}
