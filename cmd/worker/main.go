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
	
	// TODO: Initialize configuration
	// TODO: Initialize database connection
	// TODO: Initialize job queue (asynq)
	// TODO: Register job handlers
	// TODO: Start worker
	
	log.Println("✅ Venio Worker running")
	
	// Keep running
	select {}
}
