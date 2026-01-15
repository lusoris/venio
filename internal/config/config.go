// Package config handles application configuration
package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// AppConfig holds application-level settings
type AppConfig struct {
	Name    string
	Version string // Format: YYYY.MM.PATCH (e.g., 2026.01.0) - CalVer with patch number
	Env     string // development, staging, production
	Debug   bool
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	MaxConns int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Address returns the Redis connection address
func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string
	ExpirationTime    time.Duration
	RefreshExpiryDays int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
}

// DSN returns the PostgreSQL Data Source Name
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

// Load reads and parses configuration from environment and .env file
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Set default values
	setDefaults()

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		// .env not found is OK, use environment variables
	}

	// Bind environment variables
	viper.AutomaticEnv()

	// Validate required fields
	if err := validateRequired(); err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Name:    "Venio",
			Version: "2026.01.0",
			Env:     viper.GetString("APP_ENV"),
			Debug:   viper.GetBool("DEBUG"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetInt("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Database: viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
			MaxConns: viper.GetInt("POSTGRES_MAX_CONNS"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:            viper.GetString("JWT_SECRET"),
			ExpirationTime:    viper.GetDuration("JWT_EXPIRATION"),
			RefreshExpiryDays: viper.GetInt("JWT_REFRESH_EXPIRY_DAYS"),
		},
		Server: ServerConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetInt("SERVER_PORT"),
		},
	}

	return cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DEBUG", true)

	// Database defaults
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USER", "venio")
	viper.SetDefault("POSTGRES_DB", "venio")
	viper.SetDefault("POSTGRES_SSLMODE", "disable")
	viper.SetDefault("POSTGRES_MAX_CONNS", 25)

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)

	// JWT defaults
	viper.SetDefault("JWT_EXPIRATION", 24*time.Hour)
	viper.SetDefault("JWT_REFRESH_EXPIRY_DAYS", 7)

	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 3690)
}

// validateRequired checks that all required configuration is present
func validateRequired() error {
	required := []string{
		"POSTGRES_PASSWORD",
		"JWT_SECRET",
	}

	for _, key := range required {
		if viper.GetString(key) == "" {
			return fmt.Errorf("required config variable not set: %s", key)
		}
	}

	// JWT_SECRET must be at least 32 chars
	if len(viper.GetString("JWT_SECRET")) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	return nil
}

// LogConfig logs the current configuration (excluding sensitive data)
func (c *Config) LogConfig() {
	log.Printf("=== Configuration ===")
	log.Printf("App: %s v%s (env: %s)", c.App.Name, c.App.Version, c.App.Env)
	log.Printf("Database: %s@%s:%d", c.Database.User, c.Database.Host, c.Database.Port)
	log.Printf("Redis: %s:%d", c.Redis.Host, c.Redis.Port)
	log.Printf("Server: %s:%d", c.Server.Host, c.Server.Port)
	log.Printf("JWT Expiration: %v", c.JWT.ExpirationTime)
}
