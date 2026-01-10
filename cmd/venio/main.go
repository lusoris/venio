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
		fmt.Printf("Venio v%s\n", version)
		os.Exit(0)
	}

	log.Println("🚀 Starting Venio Server...")
	log.Printf("Version: %s", version)

	// TODO: Initialize configuration
	// TODO: Initialize database
	// TODO: Initialize API server
	// TODO: Start server

	log.Println("✅ Venio Server running on http://localhost:3690")

	// Keep running
	select {}
}
