---
alwaysApply: true
always_on: true
trigger: always_on
applyTo: "internal/config/**/*.go,configs/**/*,*.env.example"
description: Configuration Management Best Practices
---

# Configuration Management Best Practices

## Core Principle

**Configuration is for behavior, not secrets. Environment variables are authoritative. Fail fast on invalid config.**

## Configuration Structure

### ✅ CORRECT: Structured Config

```go
// internal/config/config.go
package config

import (
    "errors"
    "fmt"
    "time"

    "github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
    App      AppConfig
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    CORS     CORSConfig
    RateLimit RateLimitConfig
}

type AppConfig struct {
    Name        string
    Environment string  // "development", "staging", "production"
    LogLevel    string  // "debug", "info", "warn", "error"
}

type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type DatabaseConfig struct {
    Host            string
    Port            int
    User            string
    Password        string
    Name            string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

type RedisConfig struct {
    Host         string
    Port         int
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
}

type JWTConfig struct {
    Secret            string
    ExpirationTime    time.Duration
    RefreshExpiryDays int
}

type CORSConfig struct {
    AllowedOrigins []string
}

type RateLimitConfig struct {
    AuthRequests int
    AuthWindow   time.Duration
    APIRequests  int
    APIWindow    time.Duration
}
```

## Loading Configuration

### ✅ CORRECT: Viper with Validation

```go
func Load() (*Config, error) {
    // Set environment prefix
    viper.SetEnvPrefix("VENIO")
    viper.AutomaticEnv()

    // Set defaults
    setDefaults()

    // Load from .env file (optional, env vars override)
    viper.SetConfigName(".env")
    viper.SetConfigType("env")
    viper.AddConfigPath(".")
    _ = viper.ReadInConfig()  // Ignore error, env vars are primary

    // Build config
    cfg := &Config{
        App: AppConfig{
            Name:        viper.GetString("APP_NAME"),
            Environment: viper.GetString("APP_ENV"),
            LogLevel:    viper.GetString("LOG_LEVEL"),
        },
        Server: ServerConfig{
            Host:         viper.GetString("SERVER_HOST"),
            Port:         viper.GetInt("SERVER_PORT"),
            ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
            WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
            IdleTimeout:  viper.GetDuration("SERVER_IDLE_TIMEOUT"),
        },
        Database: DatabaseConfig{
            Host:            viper.GetString("DATABASE_HOST"),
            Port:            viper.GetInt("DATABASE_PORT"),
            User:            viper.GetString("DATABASE_USER"),
            Password:        viper.GetString("DATABASE_PASSWORD"),
            Name:            viper.GetString("DATABASE_NAME"),
            SSLMode:         viper.GetString("DATABASE_SSL_MODE"),
            MaxOpenConns:    viper.GetInt("DATABASE_MAX_OPEN_CONNS"),
            MaxIdleConns:    viper.GetInt("DATABASE_MAX_IDLE_CONNS"),
            ConnMaxLifetime: viper.GetDuration("DATABASE_CONN_MAX_LIFETIME"),
        },
        Redis: RedisConfig{
            Host:         viper.GetString("REDIS_HOST"),
            Port:         viper.GetInt("REDIS_PORT"),
            Password:     viper.GetString("REDIS_PASSWORD"),
            DB:           viper.GetInt("REDIS_DB"),
            PoolSize:     viper.GetInt("REDIS_POOL_SIZE"),
            MinIdleConns: viper.GetInt("REDIS_MIN_IDLE_CONNS"),
        },
        JWT: JWTConfig{
            Secret:            viper.GetString("JWT_SECRET"),
            ExpirationTime:    viper.GetDuration("JWT_EXPIRATION"),
            RefreshExpiryDays: viper.GetInt("JWT_REFRESH_EXPIRY_DAYS"),
        },
        CORS: CORSConfig{
            AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
        },
        RateLimit: RateLimitConfig{
            AuthRequests: viper.GetInt("RATE_LIMIT_AUTH_REQUESTS"),
            AuthWindow:   viper.GetDuration("RATE_LIMIT_AUTH_WINDOW"),
            APIRequests:  viper.GetInt("RATE_LIMIT_API_REQUESTS"),
            APIWindow:    viper.GetDuration("RATE_LIMIT_API_WINDOW"),
        },
    }

    // Validate configuration
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid configuration: %w", err)
    }

    return cfg, nil
}

func setDefaults() {
    // App defaults
    viper.SetDefault("APP_NAME", "venio")
    viper.SetDefault("APP_ENV", "development")
    viper.SetDefault("LOG_LEVEL", "info")

    // Server defaults
    viper.SetDefault("SERVER_HOST", "0.0.0.0")
    viper.SetDefault("SERVER_PORT", 3690)
    viper.SetDefault("SERVER_READ_TIMEOUT", 10*time.Second)
    viper.SetDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)
    viper.SetDefault("SERVER_IDLE_TIMEOUT", 120*time.Second)

    // Database defaults
    viper.SetDefault("DATABASE_HOST", "localhost")
    viper.SetDefault("DATABASE_PORT", 5432)
    viper.SetDefault("DATABASE_USER", "venio")
    viper.SetDefault("DATABASE_NAME", "venio")
    viper.SetDefault("DATABASE_SSL_MODE", "disable")
    viper.SetDefault("DATABASE_MAX_OPEN_CONNS", 25)
    viper.SetDefault("DATABASE_MAX_IDLE_CONNS", 5)
    viper.SetDefault("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute)

    // Redis defaults
    viper.SetDefault("REDIS_HOST", "localhost")
    viper.SetDefault("REDIS_PORT", 6379)
    viper.SetDefault("REDIS_DB", 0)
    viper.SetDefault("REDIS_POOL_SIZE", 10)
    viper.SetDefault("REDIS_MIN_IDLE_CONNS", 5)

    // JWT defaults
    viper.SetDefault("JWT_EXPIRATION", 24*time.Hour)
    viper.SetDefault("JWT_REFRESH_EXPIRY_DAYS", 7)

    // CORS defaults (development)
    viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"})

    // Rate limit defaults
    viper.SetDefault("RATE_LIMIT_AUTH_REQUESTS", 5)
    viper.SetDefault("RATE_LIMIT_AUTH_WINDOW", time.Minute)
    viper.SetDefault("RATE_LIMIT_API_REQUESTS", 100)
    viper.SetDefault("RATE_LIMIT_API_WINDOW", time.Minute)
}
```

## Configuration Validation

### ✅ CORRECT: Comprehensive Validation

```go
func (c *Config) Validate() error {
    var errs []error

    // Validate App
    if c.App.Name == "" {
        errs = append(errs, errors.New("APP_NAME is required"))
    }
    if !isValidEnvironment(c.App.Environment) {
        errs = append(errs, fmt.Errorf("invalid APP_ENV: %s (must be development, staging, or production)", c.App.Environment))
    }
    if !isValidLogLevel(c.App.LogLevel) {
        errs = append(errs, fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.App.LogLevel))
    }

    // Validate Server
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        errs = append(errs, fmt.Errorf("invalid SERVER_PORT: %d (must be 1-65535)", c.Server.Port))
    }

    // Validate Database
    if c.Database.Host == "" {
        errs = append(errs, errors.New("DATABASE_HOST is required"))
    }
    if c.Database.Port < 1 || c.Database.Port > 65535 {
        errs = append(errs, fmt.Errorf("invalid DATABASE_PORT: %d", c.Database.Port))
    }
    if c.Database.User == "" {
        errs = append(errs, errors.New("DATABASE_USER is required"))
    }
    if c.Database.Password == "" {
        errs = append(errs, errors.New("DATABASE_PASSWORD is required"))
    }
    if c.Database.Name == "" {
        errs = append(errs, errors.New("DATABASE_NAME is required"))
    }

    // Validate JWT
    if c.JWT.Secret == "" {
        errs = append(errs, errors.New("JWT_SECRET is required"))
    }
    if len(c.JWT.Secret) < 32 {
        errs = append(errs, errors.New("JWT_SECRET must be at least 32 characters"))
    }
    if c.JWT.ExpirationTime < time.Minute {
        errs = append(errs, errors.New("JWT_EXPIRATION must be at least 1 minute"))
    }

    // Validate CORS
    if len(c.CORS.AllowedOrigins) == 0 {
        errs = append(errs, errors.New("CORS_ALLOWED_ORIGINS is required"))
    }
    for _, origin := range c.CORS.AllowedOrigins {
        if !isValidOrigin(origin) {
            errs = append(errs, fmt.Errorf("invalid CORS origin: %s", origin))
        }
    }

    if len(errs) > 0 {
        return errors.Join(errs...)
    }

    return nil
}

func isValidEnvironment(env string) bool {
    return env == "development" || env == "staging" || env == "production"
}

func isValidLogLevel(level string) bool {
    return level == "debug" || level == "info" || level == "warn" || level == "error"
}

func isValidOrigin(origin string) bool {
    return strings.HasPrefix(origin, "http://") || strings.HasPrefix(origin, "https://")
}
```

## Database Connection String

### ✅ CORRECT: DSN Builder

```go
// DSN builds PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host,
        c.Port,
        c.User,
        c.Password,
        c.Name,
        c.SSLMode,
    )
}

// RedisAddr returns Redis address
func (c *RedisConfig) Addr() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
```

## Environment File Structure

### .env.example (Template)

```bash
# ===================================================================
# VENIO APPLICATION CONFIGURATION
# ===================================================================
# Copy this file to .env and replace placeholder values with actual values.
# DO NOT commit .env to version control!

# -------------------------------------------------------------------
# Application
# -------------------------------------------------------------------
VENIO_APP_NAME=venio
VENIO_APP_ENV=development  # development, staging, production
VENIO_LOG_LEVEL=info       # debug, info, warn, error

# -------------------------------------------------------------------
# Server
# -------------------------------------------------------------------
VENIO_SERVER_HOST=0.0.0.0
VENIO_SERVER_PORT=3690
VENIO_SERVER_READ_TIMEOUT=10s
VENIO_SERVER_WRITE_TIMEOUT=10s
VENIO_SERVER_IDLE_TIMEOUT=120s

# -------------------------------------------------------------------
# Database (PostgreSQL)
# -------------------------------------------------------------------
VENIO_DATABASE_HOST=localhost
VENIO_DATABASE_PORT=5432
VENIO_DATABASE_USER=venio
VENIO_DATABASE_PASSWORD=change-me-in-production
VENIO_DATABASE_NAME=venio
VENIO_DATABASE_SSL_MODE=disable  # disable, require, verify-ca, verify-full
VENIO_DATABASE_MAX_OPEN_CONNS=25
VENIO_DATABASE_MAX_IDLE_CONNS=5
VENIO_DATABASE_CONN_MAX_LIFETIME=5m

# -------------------------------------------------------------------
# Redis
# -------------------------------------------------------------------
VENIO_REDIS_HOST=localhost
VENIO_REDIS_PORT=6379
VENIO_REDIS_PASSWORD=change-me-in-production
VENIO_REDIS_DB=0
VENIO_REDIS_POOL_SIZE=10
VENIO_REDIS_MIN_IDLE_CONNS=5

# -------------------------------------------------------------------
# JWT Authentication
# -------------------------------------------------------------------
# IMPORTANT: Generate a strong random secret (min 32 characters)
# Example: openssl rand -base64 32
VENIO_JWT_SECRET=generate-random-32-char-string-here
VENIO_JWT_EXPIRATION=24h
VENIO_JWT_REFRESH_EXPIRY_DAYS=7

# -------------------------------------------------------------------
# CORS
# -------------------------------------------------------------------
# Comma-separated list of allowed origins
VENIO_CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001

# -------------------------------------------------------------------
# Rate Limiting
# -------------------------------------------------------------------
VENIO_RATE_LIMIT_AUTH_REQUESTS=5
VENIO_RATE_LIMIT_AUTH_WINDOW=1m
VENIO_RATE_LIMIT_API_REQUESTS=100
VENIO_RATE_LIMIT_API_WINDOW=1m
```

### .env (Actual Values - NEVER COMMIT)

```bash
# Copy from .env.example and replace with real values
VENIO_JWT_SECRET=7Jk9mN2pQ5rS8tV1wX4yZ6aB3cD0eF7gH9iK2lM5nO8pQ
VENIO_DATABASE_PASSWORD=K8mP2qT5vY9zA3cF6iL0oR4sU7xB1dE
VENIO_REDIS_PASSWORD=nQ9sV2yB5eH8kN1pT4wZ7aC0fI3lO6rU
```

## Multi-Environment Strategy

### Development

```bash
# .env.development
VENIO_APP_ENV=development
VENIO_LOG_LEVEL=debug
VENIO_DATABASE_SSL_MODE=disable
VENIO_CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

### Staging

```bash
# .env.staging
VENIO_APP_ENV=staging
VENIO_LOG_LEVEL=info
VENIO_DATABASE_SSL_MODE=require
VENIO_CORS_ALLOWED_ORIGINS=https://staging.venio.io
```

### Production

```bash
# .env.production
VENIO_APP_ENV=production
VENIO_LOG_LEVEL=warn
VENIO_DATABASE_SSL_MODE=verify-full
VENIO_CORS_ALLOWED_ORIGINS=https://app.venio.io
```

## Configuration in main.go

### ✅ CORRECT: Fail Fast on Startup

```go
// cmd/venio/main.go
func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize logger
    logger := initLogger(cfg.App.LogLevel)
    logger.Info("Starting application",
        "name", cfg.App.Name,
        "environment", cfg.App.Environment,
        "port", cfg.Server.Port,
    )

    // Validate environment-specific requirements
    if cfg.App.Environment == "production" {
        if cfg.Database.SSLMode == "disable" {
            log.Fatal("Production requires DATABASE_SSL_MODE=require or higher")
        }
        if !strings.HasPrefix(cfg.CORS.AllowedOrigins[0], "https://") {
            log.Fatal("Production requires HTTPS origins only")
        }
    }

    // Initialize database
    db, err := database.Connect(cfg.Database.DSN())
    if err != nil {
        logger.Error("Failed to connect to database", "error", err)
        log.Fatal(err)
    }
    defer db.Close()

    // ... rest of initialization
}
```

## Configuration Testing

### ✅ CORRECT: Test Configuration Loading

```go
// internal/config/config_test.go
func TestLoad_ValidConfig(t *testing.T) {
    // Set required environment variables
    os.Setenv("VENIO_JWT_SECRET", "this-is-a-test-secret-that-is-at-least-32-characters-long")
    os.Setenv("VENIO_DATABASE_PASSWORD", "test-password")
    defer os.Unsetenv("VENIO_JWT_SECRET")
    defer os.Unsetenv("VENIO_DATABASE_PASSWORD")

    cfg, err := Load()
    assert.NoError(t, err)
    assert.NotNil(t, cfg)
    assert.Equal(t, "venio", cfg.App.Name)
}

func TestValidate_MissingJWTSecret(t *testing.T) {
    cfg := &Config{
        App: AppConfig{Name: "venio", Environment: "development"},
        JWT: JWTConfig{Secret: ""}, // Missing
    }

    err := cfg.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "JWT_SECRET is required")
}

func TestValidate_ShortJWTSecret(t *testing.T) {
    cfg := &Config{
        App: AppConfig{Name: "venio", Environment: "development"},
        JWT: JWTConfig{Secret: "too-short"}, // Less than 32 chars
    }

    err := cfg.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "JWT_SECRET must be at least 32 characters")
}
```

## Common Mistakes

### ❌ DON'T Do This

```go
// 1. Hardcoded secrets
const jwtSecret = "my-secret"  // ❌ NEVER!

// 2. No validation
func Load() *Config {
    return &Config{
        JWT: JWTConfig{Secret: os.Getenv("JWT_SECRET")},  // ❌ No validation!
    }
}

// 3. Silently use defaults for secrets
func Load() *Config {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "default-secret"  // ❌ Insecure default!
    }
    return &Config{JWT: JWTConfig{Secret: secret}}
}

// 4. Mixed configuration sources without priority
// Read from file, env, flags - which takes precedence? Unclear!

// 5. No environment-specific validation
// Allow DATABASE_SSL_MODE=disable in production

// 6. Panic on missing optional config
port := viper.GetInt("SERVER_PORT")  // ❌ Returns 0 if missing, no error!
if port == 0 {
    panic("PORT not set")  // ❌ Panic for optional config
}
```

### ✅ DO This

```go
// 1. Load from environment
secret := os.Getenv("JWT_SECRET")

// 2. Validate thoroughly
if err := cfg.Validate(); err != nil {
    return nil, err
}

// 3. Fail fast on missing secrets
if secret == "" {
    return nil, errors.New("JWT_SECRET is required")
}

// 4. Clear precedence: env vars > .env file > defaults
viper.AutomaticEnv()  // Env vars override everything
viper.ReadInConfig()  // .env file
setDefaults()         // Defaults only if not set

// 5. Environment-specific validation
if cfg.App.Environment == "production" {
    if cfg.Database.SSLMode == "disable" {
        return nil, errors.New("production requires SSL")
    }
}

// 6. Validate optional config if set
if port != 0 && (port < 1 || port > 65535) {
    return nil, fmt.Errorf("invalid port: %d", port)
}
```

## Configuration Checklist

- [ ] **All secrets from environment variables**
- [ ] **No hardcoded secrets or passwords**
- [ ] **Validation on startup (fail fast)**
- [ ] **Sensible defaults for non-secrets**
- [ ] **Clear precedence: env > file > defaults**
- [ ] **.env.example committed (template)**
- [ ] **.env in .gitignore**
- [ ] **Environment-specific validation (SSL in prod)**
- [ ] **Struct-based config (not map[string]interface{})**
- [ ] **DSN/connection string builders**
- [ ] **Tests for validation logic**

---

## Document Version

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-15 | AI Assistant | Initial configuration management guide |

**Remember**: Configuration should be easy to change between environments without code changes. Secrets must never be committed to version control.
