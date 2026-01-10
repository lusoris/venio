module github.com/lusoris/venio

go 1.23

require (
	// Web Framework
	github.com/gin-gonic/gin v1.10.0
	
	// Database
	github.com/jackc/pgx/v5 v5.7.2
	
	// Redis
	github.com/redis/go-redis/v9 v9.7.0
	
	// Configuration
	github.com/spf13/viper v1.19.0
	
	// Jobs
	github.com/hibiken/asynq v0.24.1
	
	// Testing
	github.com/stretchr/testify v1.10.0
	
	// OpenAPI/Swagger
	github.com/swaggo/swag v1.16.4
	github.com/swaggo/gin-swagger v1.6.0
)
