# Configuration Reference

## Environment Variables

### Server

```bash
PORT=3690                    # API server port
ENV=production              # Environment: development|production|test
LOG_LEVEL=info              # Log level: debug|info|warn|error
LOG_FORMAT=json             # Log format: json|text
```

### Database (PostgreSQL)

```bash
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=venio
POSTGRES_PASSWORD=changeme
POSTGRES_DB=venio
DATABASE_URL=postgres://venio:changeme@postgres:5432/venio?sslmode=disable
```

### Redis

```bash
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=changeme
REDIS_DB=0
```

### Typesense

```bash
TYPESENSE_HOST=typesense
TYPESENSE_PORT=8108
TYPESENSE_API_KEY=changeme
```

### Security

```bash
JWT_SECRET=minimum_32_characters_required
API_KEY=changeme
```

### OIDC (Optional)

```bash
OIDC_ENABLED=false
OIDC_ISSUER=https://auth.example.com
OIDC_CLIENT_ID=venio
OIDC_CLIENT_SECRET=changeme
```

### External Services (Optional)

```bash
# Overseerr
OVERSEERR_URL=http://overseerr:5055
OVERSEERR_API_KEY=

# Arrs
SONARR_URL=http://sonarr:8989
SONARR_API_KEY=

RADARR_URL=http://radarr:7878
RADARR_API_KEY=

LIDARR_URL=http://lidarr:8686
LIDARR_API_KEY=

WHISPARR_URL=http://whisparr:6969
WHISPARR_API_KEY=
```

## Configuration Files

### config.yml (Future)

```yaml
# Example structure (not yet implemented)
server:
  port: 3690
  timeout: 30s

modules:
  overseerr:
    enabled: true
  lidarr:
    enabled: true
  whisparr:
    enabled: false

features:
  voting:
    enabled: true
    threshold: 3
```

## Best Practices

### Security

1. **Use strong secrets** - Generate random 32+ character strings
2. **Never commit .env** - It's gitignored for a reason
3. **Rotate keys regularly** - Especially API keys
4. **Use OIDC in production** - Centralized auth is safer

### Performance

1. **Configure connection pools** - PostgreSQL/Redis
2. **Enable caching** - Redis for metadata
3. **Set appropriate timeouts** - Prevent hanging requests

### Deployment

1. **Use environment-specific configs** - dev/staging/prod
2. **Keep secrets in secrets manager** - Not in env files
3. **Monitor resource usage** - Adjust as needed

---

## Document Revision

| Version | Date | Author | Changes | Source Version |
|---------|------|--------|---------|----------------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial configuration reference | - |

## Referenced Documentation

- **PostgreSQL Configuration:** [PostgreSQL 18 Configuration](https://www.postgresql.org/docs/18/runtime-config.html) (Released: 2025-11-14)
- **Redis Configuration:** [Redis 8.4 Configuration](https://redis.io/docs/latest/operate/oss_and_stack/management/config/) (Released: 2025-12-15)
- **Environment Variables:** [12 Factor App - Config](https://12factor.net/config)
- **JWT Best Practices:** [JWT Handbook](https://auth0.com/resources/ebooks/jwt-handbook)
- **Database Connection Pooling:** [pgx Pool Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)

*Full configuration options coming soon*
