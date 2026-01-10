# Internal

This directory contains private application code that cannot be imported by external packages.

## Structure

- **api/** - HTTP handlers, routing, middleware
- **services/** - Business logic layer
- **database/** - Database queries and models (sqlc generated)
- **models/** - Shared data structures
- **config/** - Configuration loading and management
- **providers/** - External API clients (Overseerr, Arrs, etc.)
- **proxy/** - Metadata proxy implementation

## Guidelines

- Keep packages focused and loosely coupled
- Use dependency injection
- Write tests for all business logic
- Document exported functions and types
