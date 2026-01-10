# Database Migrations

SQL migration files for database schema changes.

## Naming Convention

```
YYYYMMDDHHMMSS_description.up.sql
YYYYMMDDHHMMSS_description.down.sql
```

Example:
```
20260110000001_create_users_table.up.sql
20260110000001_create_users_table.down.sql
```

## Creating Migrations

```bash
# TODO: Add migration creation command
```

## Running Migrations

```bash
make migrate-up    # Apply pending migrations
make migrate-down  # Rollback last migration
```

## Best Practices

- Always provide both up and down migrations
- Test migrations on dev database first
- Keep migrations atomic and focused
- Never modify committed migrations
