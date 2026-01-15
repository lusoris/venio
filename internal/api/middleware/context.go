package middleware

import "github.com/gin-gonic/gin"

// Context keys for type-safe value storage
type contextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	// EmailKey is the context key for email
	EmailKey contextKey = "email"
	// UsernameKey is the context key for username
	UsernameKey contextKey = "username"
	// RolesKey is the context key for user roles
	RolesKey contextKey = "roles"
	// RequestIDKey is the context key for request ID
	RequestIDKey contextKey = "request_id"
)

// SetUserID stores user ID in context
func SetUserID(c *gin.Context, userID int64) {
	c.Set(string(UserIDKey), userID)
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(string(UserIDKey))
	if !exists {
		return 0, false
	}
	userID, ok := value.(int64)
	return userID, ok
}

// SetEmail stores email in context
func SetEmail(c *gin.Context, email string) {
	c.Set(string(EmailKey), email)
}

// GetEmail retrieves email from context
func GetEmail(c *gin.Context) (string, bool) {
	value, exists := c.Get(string(EmailKey))
	if !exists {
		return "", false
	}
	email, ok := value.(string)
	return email, ok
}

// SetUsername stores username in context
func SetUsername(c *gin.Context, username string) {
	c.Set(string(UsernameKey), username)
}

// GetUsername retrieves username from context
func GetUsername(c *gin.Context) (string, bool) {
	value, exists := c.Get(string(UsernameKey))
	if !exists {
		return "", false
	}
	username, ok := value.(string)
	return username, ok
}

// SetRoles stores user roles in context
func SetRoles(c *gin.Context, roles []string) {
	c.Set(string(RolesKey), roles)
}

// GetRoles retrieves user roles from context
func GetRoles(c *gin.Context) ([]string, bool) {
	value, exists := c.Get(string(RolesKey))
	if !exists {
		return nil, false
	}
	roles, ok := value.([]string)
	return roles, ok
}

// SetRequestID stores request ID in context
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(string(RequestIDKey), requestID)
}

// GetRequestID retrieves request ID from context
func GetRequestID(c *gin.Context) (string, bool) {
	value, exists := c.Get(string(RequestIDKey))
	if !exists {
		return "", false
	}
	requestID, ok := value.(string)
	return requestID, ok
}
