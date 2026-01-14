// Package api contains API routing and handler setup
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/lusoris/venio/internal/api/handlers"
	"github.com/lusoris/venio/internal/api/middleware"
	"github.com/lusoris/venio/internal/config"
	"github.com/lusoris/venio/internal/database"
	"github.com/lusoris/venio/internal/repositories"
	"github.com/lusoris/venio/internal/services"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(cfg *config.Config, db *database.DB) *gin.Engine {
	router := gin.Default()

	// Initialize repositories
	userRepo := repositories.NewPostgresUserRepository(db.Pool())

	// Initialize services
	userService := services.NewDefaultUserService(userRepo)
	authService := services.NewDefaultAuthService(userService, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":     cfg.App.Name,
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected user routes
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(authService))
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return router
}
