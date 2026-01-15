package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	version string
	env     string
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(db *pgxpool.Pool, redis *redis.Client, version, env string) *HealthHandler {
	return &HealthHandler{
		db:      db,
		redis:   redis,
		version: version,
		env:     env,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Version   string                 `json:"version"`
	Services  map[string]ServiceInfo `json:"services"`
}

// ServiceInfo represents the status of a service
type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Liveness godoc
// @Summary Liveness probe
// @Description Check if the application is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness godoc
// @Summary Readiness probe
// @Description Check if the application is ready to serve traffic
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health/ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Check dependencies
	dbStatus := h.checkDatabase(ctx)
	redisStatus := h.checkRedis(ctx)
	allHealthy := dbStatus.Status == "healthy" && redisStatus.Status == "healthy"

	// In production, return minimal information
	if h.env == "production" {
		if allHealthy {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		}
		return
	}

	// In development, return detailed information
	response := HealthResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   h.version,
		Services:  make(map[string]ServiceInfo),
	}

	response.Services["database"] = dbStatus
	response.Services["redis"] = redisStatus

	if allHealthy {
		response.Status = "healthy"
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

func (h *HealthHandler) checkDatabase(ctx context.Context) ServiceInfo {
	if err := h.db.Ping(ctx); err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: "Database connection failed",
		}
	}

	return ServiceInfo{
		Status: "healthy",
	}
}

func (h *HealthHandler) checkRedis(ctx context.Context) ServiceInfo {
	if err := h.redis.Ping(ctx).Err(); err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: "Redis connection failed",
		}
	}

	return ServiceInfo{
		Status: "healthy",
	}
}
