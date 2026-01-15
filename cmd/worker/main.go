// Copyright (C) 2026 Venio Contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0
//
// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"fmt"
	"log"
	"os"
)

const version = "0.1.0-dev"

func main() {
	// Parse CLI flags
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Venio Worker v%s\n", version)
		os.Exit(0)
	}

	log.Println("🔧 Starting Venio Worker...")
	log.Printf("Version: %s", version)

	// Worker is currently a stub and not implemented
	// Future implementation will use Asynq for background job processing:
	// - Email notifications
	// - Media processing
	// - Database cleanup tasks
	// - Integration with external services (Arr, Overseerr)
	//
	// For now, all operations run synchronously in the main server.
	// See: https://github.com/lusoris/venio/issues/[TBD] for roadmap

	log.Println("⚠️  Worker stub - not yet implemented")
	log.Println("    Background tasks currently run in main server")

	// Exit instead of hanging
	os.Exit(0)
}
