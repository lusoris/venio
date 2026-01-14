package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(3, 1*time.Second)

	// First 3 requests should succeed
	assert.True(t, rl.Allow("192.168.1.1"))
	assert.True(t, rl.Allow("192.168.1.1"))
	assert.True(t, rl.Allow("192.168.1.1"))

	// Fourth request should fail
	assert.False(t, rl.Allow("192.168.1.1"))
}

func TestRateLimiter_Different_IPs(t *testing.T) {
	rl := NewRateLimiter(2, 1*time.Second)

	// Each IP should have its own limit
	assert.True(t, rl.Allow("192.168.1.1"))
	assert.True(t, rl.Allow("192.168.1.1"))
	assert.False(t, rl.Allow("192.168.1.1"))

	// Different IP should still have its limit
	assert.True(t, rl.Allow("192.168.1.2"))
	assert.True(t, rl.Allow("192.168.1.2"))
	assert.False(t, rl.Allow("192.168.1.2"))
}

func TestRateLimiter_Window_Reset(t *testing.T) {
	rl := NewRateLimiter(2, 100*time.Millisecond)

	assert.True(t, rl.Allow("192.168.1.1"))
	assert.True(t, rl.Allow("192.168.1.1"))
	assert.False(t, rl.Allow("192.168.1.1"))

	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)

	// Should be able to make requests again
	assert.True(t, rl.Allow("192.168.1.1"))
}

func TestRateLimiter_Middleware_Success(t *testing.T) {
	rl := NewRateLimiter(2, 1*time.Second)

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First two requests should succeed
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Third request should be rate limited
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestSecurityHeaders_All_Headers_Present(t *testing.T) {
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check all security headers are present
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
	assert.NotEmpty(t, w.Header().Get("Permissions-Policy"))
}

func TestStrictSecurityHeaders_All_Headers_Present(t *testing.T) {
	router := gin.New()
	router.Use(StrictSecurityHeaders(31536000))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check all strict security headers
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Contains(t, w.Header().Get("Content-Security-Policy"), "default-src 'none'")
	assert.Equal(t, "strict-no-referrer", w.Header().Get("Referrer-Policy"))
	assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "max-age=31536000")
}

func TestAuthRateLimiter_Default_Settings(t *testing.T) {
	rl := AuthRateLimiter()

	// Should allow 5 requests per minute
	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow("test-ip"))
	}

	// 6th should fail
	assert.False(t, rl.Allow("test-ip"))
}

func TestGeneralRateLimiter_Default_Settings(t *testing.T) {
	rl := GeneralRateLimiter()

	// Should allow 100 requests per minute
	for i := 0; i < 100; i++ {
		assert.True(t, rl.Allow("test-ip"))
	}

	// 101st should fail
	assert.False(t, rl.Allow("test-ip"))
}
