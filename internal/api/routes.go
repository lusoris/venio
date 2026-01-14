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
	roleRepo := repositories.NewRoleRepository(db.Pool())
	permissionRepo := repositories.NewPermissionRepository(db.Pool())
	userRoleRepo := repositories.NewUserRoleRepository(db.Pool())

	// Initialize services
	userService := services.NewDefaultUserService(userRepo)
	authService := services.NewDefaultAuthService(userService, cfg)
	roleService := services.NewRoleService(roleRepo)
	permissionService := services.NewPermissionService(permissionRepo)
	userRoleService := services.NewUserRoleService(userRoleRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	userRoleHandler := handlers.NewUserRoleHandler(userRoleService)
	adminHandler := handlers.NewAdminHandler(userService, roleService, permissionService, userRoleService)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(authService)
	rbacMiddleware := middleware.NewRBACMiddleware(userRoleService)

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
		users.Use(authMiddleware)
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			// User-role management routes
			users.GET("/:userId/roles", userRoleHandler.GetUserRoles)
			users.POST("/:userId/roles", rbacMiddleware.RequireRole("admin"), userRoleHandler.AssignRoleToUser)
			users.DELETE("/:userId/roles/:roleId", rbacMiddleware.RequireRole("admin"), userRoleHandler.RemoveRoleFromUser)
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
