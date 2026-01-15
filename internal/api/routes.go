// Package api contains API routing and handler setup
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lusoris/venio/docs/swagger" // Import generated docs
	"github.com/lusoris/venio/internal/api/handlers"
	"github.com/lusoris/venio/internal/api/middleware"
	"github.com/lusoris/venio/internal/config"
	"github.com/lusoris/venio/internal/database"
	"github.com/lusoris/venio/internal/logger"
	redisClient "github.com/lusoris/venio/internal/redis"
	"github.com/lusoris/venio/internal/repositories"
	"github.com/lusoris/venio/internal/services"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(cfg *config.Config, db *database.DB, redis *redisClient.Client, log *logger.Logger) *gin.Engine {
	router := gin.Default()

	// Apply structured logging middleware
	router.Use(middleware.LoggingMiddleware(log))

	// Apply Prometheus metrics middleware
	router.Use(middleware.PrometheusMiddleware())

	// Apply global security middleware
	router.Use(middleware.SecurityHeaders())

	// Apply CORS middleware
	if cfg.App.Env == "development" {
		router.Use(middleware.CORSDevelopment())
	} else {
		// In production, specify the frontend URL from config
		// For now, allow localhost:3000 for testing
		router.Use(middleware.CORS("http://localhost:3000"))
	}

	// Initialize repositories
	userRepo := repositories.NewPostgresUserRepository(db.Pool())
	roleRepo := repositories.NewRoleRepository(db.Pool())
	permissionRepo := repositories.NewPermissionRepository(db.Pool())
	userRoleRepo := repositories.NewUserRoleRepository(db.Pool())

	// Initialize services
	userService := services.NewDefaultUserService(userRepo)
	userRoleService := services.NewUserRoleService(userRoleRepo)
	authService := services.NewDefaultAuthService(userService, userRoleService, cfg)
	roleService := services.NewRoleService(roleRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	userRoleHandler := handlers.NewUserRoleHandler(userRoleService)
	adminHandler := handlers.NewAdminHandler(userService, roleService, permissionService, userRoleService)
	healthHandler := handlers.NewHealthHandler(db.Pool(), redis.Client, cfg.App.Version, cfg.App.Env)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(authService)
	rbacMiddleware := middleware.NewRBACMiddleware(userRoleService)

	// Initialize Redis-based rate limiters (distributed, production-ready)
	authRateLimiter := middleware.RedisAuthRateLimiter(redis.Client)
	generalRateLimiter := middleware.RedisGeneralRateLimiter(redis.Client)

	// Metrics endpoint (Prometheus)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoints (Kubernetes probes)
	health := router.Group("/health")
	{
		health.GET("/live", healthHandler.Liveness)
		health.GET("/ready", healthHandler.Readiness)
	}

	// API documentation (Swagger UI) - Only in development
	if cfg.App.Env != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	v1.Use(generalRateLimiter.Middleware())
	{
		// Public auth routes (with stricter rate limiting)
		auth := v1.Group("/auth")
		auth.Use(authRateLimiter.Middleware())
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected user routes
		users := v1.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			// User-role management routes
			users.GET("/:id/roles", userRoleHandler.GetUserRoles)
			users.POST("/:id/roles", rbacMiddleware.RequireRole("admin"), userRoleHandler.AssignRoleToUser)
			users.DELETE("/:id/roles/:roleId", rbacMiddleware.RequireRole("admin"), userRoleHandler.RemoveRoleFromUser)
		}

		// Protected role routes
		roles := v1.Group("/roles")
		roles.Use(authMiddleware, rbacMiddleware.RequireRole("admin"))
		{
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.POST("", roleHandler.CreateRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			// Role-permission management routes
			roles.GET("/:id/permissions", roleHandler.GetRolePermissions)
			roles.POST("/:id/permissions", roleHandler.AssignPermissionToRole)
			roles.DELETE("/:id/permissions/:permissionId", roleHandler.RemovePermissionFromRole)
		}

		// Protected permission routes
		permissions := v1.Group("/permissions")
		permissions.Use(authMiddleware, rbacMiddleware.RequireRole("admin"))
		{
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)
			permissions.POST("", permissionHandler.CreatePermission)
			permissions.PUT("/:id", permissionHandler.UpdatePermission)
			permissions.DELETE("/:id", permissionHandler.DeletePermission)
		}

		// Admin-only routes
		admin := v1.Group("/admin")
		admin.Use(authMiddleware, rbacMiddleware.RequireRole("admin"))
		{
			// User management
			admin.GET("/users", adminHandler.ListUsers)
			admin.POST("/users", adminHandler.CreateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			// Role management
			admin.GET("/roles", adminHandler.ListRoles)
			admin.POST("/roles", adminHandler.CreateRole)
			admin.DELETE("/roles/:id", adminHandler.DeleteRole)

			// Permission management
			admin.GET("/permissions", adminHandler.ListPermissions)

			// User-role assignments
			admin.GET("/user-roles", adminHandler.ListUserRoles)
			admin.DELETE("/user-roles/:id", adminHandler.RemoveUserRole)
		}
	}

	return router
}
