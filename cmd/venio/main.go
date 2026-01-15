// Copyright (C) 2026 Venio Contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0
//
// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lusoris/venio/internal/api"
	"github.com/lusoris/venio/internal/config"
	"github.com/lusoris/venio/internal/database"
	"github.com/lusoris/venio/internal/logger"
	redisClient "github.com/lusoris/venio/internal/redis"
)

const version = "0.1.0-dev"

func main() {
	// Parse CLI flags
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Venio v%s\n", version)
		os.Exit(0)
	}

	log.Println("🚀 Starting Venio Server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg.LogConfig()

	// Initialize structured logger
	appLogger := logger.New(&cfg.App)
	appLogger.Info("Starting Venio Server", "version", version, "env", cfg.App.Env)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize database
	db, err := database.Connect(ctx, &cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v", closeErr)
		}
	}()

	// Initialize Redis
	redis, err := redisClient.Connect(ctx, &cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}()

	// Setup router with all routes
	router := api.SetupRouter(cfg, db, redis, appLogger)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("✅ Venio Server running on http://localhost:%d", cfg.Server.Port)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("⏹️  Shutting down server...")
}
