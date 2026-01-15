# Admin Operations Guide

Complete guide for day-to-day administration of Venio.

## Table of Contents

- [User Management](#user-management)
- [Role & Permission Management](#role--permission-management)
- [Monitoring & Alerts](#monitoring--alerts)
- [Database Operations](#database-operations)
- [Logs & Troubleshooting](#logs--troubleshooting)
- [Performance Tuning](#performance-tuning)
- [Backup & Recovery](#backup--recovery)
- [Security Operations](#security-operations)

## User Management

### Creating Users

**Via API:**
```bash
curl -X POST https://api.venio.dev/api/v1/admin/users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "newuser",
    "first_name": "John",
    "last_name": "Doe",
    "password": "SecurePassword123!"
  }'
```

**Via Database (Emergency):**
```sql
-- Generate bcrypt hash first
-- Cost 12: echo -n "password" | htpasswd -nBC 12 "" | cut -d: -f2

INSERT INTO users (email, username, first_name, last_name, password, is_active)
VALUES ('user@example.com', 'newuser', 'John', 'Doe', '$2a$12$...', true);
```

### Disabling Users

**Temporary Suspension:**
```bash
curl -X PUT https://api.venio.dev/api/v1/admin/users/123 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'
```

**Permanent Deletion:**
```bash
curl -X DELETE https://api.venio.dev/api/v1/admin/users/123 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Listing Users

```bash
# All users
curl https://api.venio.dev/api/v1/admin/users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Filter by status
curl "https://api.venio.dev/api/v1/admin/users?is_active=true" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Password Reset

**Force password reset:**
```sql
UPDATE users 
SET password = '$2a$12$NEW_BCRYPT_HASH_HERE'
WHERE email = 'user@example.com';
```

## Role & Permission Management

### Assigning Roles to Users

```bash
curl -X POST https://api.venio.dev/api/v1/users/123/roles \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 2
  }'
```

### Removing Roles from Users

```bash
curl -X DELETE https://api.venio.dev/api/v1/users/123/roles/2 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Creating Custom Roles

```bash
curl -X POST https://api.venio.dev/api/v1/roles \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "content_moderator",
    "description": "Can moderate user-generated content"
  }'
```

### Assigning Permissions to Roles

```bash
curl -X POST https://api.venio.dev/api/v1/roles/5/permissions \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_id": 10
  }'
```

### Viewing Role Hierarchy

```bash
# List all roles with their permissions
curl https://api.venio.dev/api/v1/admin/roles \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Get specific role details
curl https://api.venio.dev/api/v1/roles/2/permissions \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## Monitoring & Alerts

### Accessing Grafana

1. Navigate to `https://grafana.venio.dev`
2. Login with admin credentials
3. Navigate to Dashboards → Venio Overview

### Key Metrics to Monitor

**Application Health:**
- Request rate (should be steady, spikes need investigation)
- Error rate (should be <1%)
- P95 latency (should be <500ms)
- Active users

**Database:**
- Connection pool utilization (<80%)
- Query duration
- Failed queries
- Connection errors

**Redis:**
- Command rate
- Command duration
- Rate limit denials
- Memory usage

### Managing Alerts

**Silencing Alerts (Maintenance Window):**

1. Go to Alertmanager: `https://alerts.venio.dev`
2. Click "New Silence"
3. Set matchers:
   - `alertname = HighLatency`
   - `severity = warning`
4. Set duration (e.g., 2 hours)
5. Add comment: "Planned maintenance"

**Modifying Alert Thresholds:**

Edit `deployments/prometheus/alerts.yml`:

```yaml
- alert: HighErrorRate
  expr: |
    sum(rate(venio_http_requests_total{status=~"5.."}[5m]))
    /
    sum(rate(venio_http_requests_total[5m])) > 0.05  # Change threshold here
  for: 5m  # Change evaluation period
```

Apply changes:
```bash
docker compose restart prometheus
```

### Viewing Metrics Directly

**Prometheus Query:**
```bash
curl 'https://prometheus.venio.dev/api/v1/query?query=venio_http_requests_total'
```

**From Application:**
```bash
curl https://api.venio.dev/metrics
```

## Database Operations

### Checking Database Status

```bash
# Via health endpoint
curl https://api.venio.dev/health/ready

# Direct check
docker compose exec postgres pg_isready -U venio
```

### Running Migrations

```bash
# Check current version
docker compose exec venio ./venio migrate version

# Migrate up
docker compose exec venio ./venio migrate up

# Rollback (if needed)
docker compose exec venio ./venio migrate down
```

### Viewing Active Connections

```sql
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query_start,
    state_change
FROM pg_stat_activity
WHERE datname = 'venio'
ORDER BY query_start DESC;
```

### Terminating Stuck Queries

```sql
-- Find long-running queries
SELECT 
    pid,
    now() - query_start as duration,
    query
FROM pg_stat_activity
WHERE state = 'active'
  AND now() - query_start > interval '5 minutes';

-- Kill specific query
SELECT pg_terminate_backend(pid);
```

### Database Maintenance

**Vacuum (cleanup):**
```bash
docker compose exec postgres psql -U venio -d venio -c "VACUUM ANALYZE;"
```

**Reindex (performance):**
```bash
docker compose exec postgres psql -U venio -d venio -c "REINDEX DATABASE venio;"
```

## Logs & Troubleshooting

### Viewing Logs

**Application logs:**
```bash
# Real-time
docker compose logs -f venio

# Last 100 lines
docker compose logs --tail=100 venio

# Specific time range
docker compose logs --since="2026-01-15T10:00:00" --until="2026-01-15T11:00:00" venio
```

**Database logs:**
```bash
docker compose logs postgres
```

**Redis logs:**
```bash
docker compose logs redis
```

### Log Analysis

**Find errors:**
```bash
docker compose logs venio | grep "level=ERROR"
```

**Find slow queries:**
```bash
docker compose logs venio | grep "duration" | awk '$NF > 1000'
```

**Authentication failures:**
```bash
docker compose logs venio | grep "Authentication failed"
```

### Common Issues

#### High Memory Usage

**Check container stats:**
```bash
docker stats venio-api
```

**If memory is high:**
```bash
# Restart application
docker compose restart venio

# Check for memory leaks in metrics
curl https://api.venio.dev/metrics | grep process_resident_memory_bytes
```

#### Database Connection Errors

**Symptoms:** "Connection pool exhausted" errors

**Fix:**
1. Check active connections (see Database Operations)
2. Increase pool size in `.env`:
   ```
   DB_MAX_CONNECTIONS=50
   ```
3. Restart application

#### Redis Connection Errors

**Check Redis health:**
```bash
docker compose exec redis redis-cli -a YOUR_PASSWORD ping
```

**Restart Redis:**
```bash
docker compose restart redis
```

## Performance Tuning

### Database Optimization

**Check slow queries:**
```sql
-- Enable pg_stat_statements (if not enabled)
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Find slow queries
SELECT 
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Add indexes:**
```sql
-- Example: Index on email lookups
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- Monitor index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

### Redis Optimization

**Check memory usage:**
```bash
docker compose exec redis redis-cli -a YOUR_PASSWORD INFO memory
```

**Configure eviction policy:**
```bash
# In redis.conf or via command
docker compose exec redis redis-cli -a YOUR_PASSWORD CONFIG SET maxmemory 1gb
docker compose exec redis redis-cli -a YOUR_PASSWORD CONFIG SET maxmemory-policy allkeys-lru
```

### Application Optimization

**Check connection pool stats:**
```bash
curl https://api.venio.dev/metrics | grep venio_db_connections
```

**Tune pool settings in `.env`:**
```bash
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_LIFETIME=15m
```

## Backup & Recovery

### Automated Backups

**Setup cron job:**
```bash
crontab -e

# Add:
0 2 * * * /home/venio/backup-db.sh
```

**Backup script:**
```bash
#!/bin/bash
BACKUP_DIR=/home/venio/backups
DATE=$(date +%Y%m%d_%H%M%S)

# Database backup
docker compose exec -T postgres pg_dump -U venio venio > $BACKUP_DIR/venio_$DATE.sql
gzip $BACKUP_DIR/venio_$DATE.sql

# Redis backup (RDB snapshot)
docker compose exec redis redis-cli -a YOUR_PASSWORD BGSAVE
cp /var/lib/docker/volumes/venio_redis-data/_data/dump.rdb $BACKUP_DIR/redis_$DATE.rdb

# Keep last 7 days
find $BACKUP_DIR -name "venio_*.sql.gz" -mtime +7 -delete
find $BACKUP_DIR -name "redis_*.rdb" -mtime +7 -delete
```

### Manual Backup

```bash
# Database
docker compose exec postgres pg_dump -U venio venio > backup_$(date +%Y%m%d).sql

# Redis
docker compose exec redis redis-cli -a YOUR_PASSWORD SAVE
```

### Restore from Backup

**Database:**
```bash
# Stop application first
docker compose stop venio

# Restore
cat backup_20260115.sql | docker compose exec -T postgres psql -U venio -d venio

# Start application
docker compose start venio
```

**Redis:**
```bash
# Stop Redis
docker compose stop redis

# Replace dump.rdb
cp backup_redis.rdb /var/lib/docker/volumes/venio_redis-data/_data/dump.rdb

# Start Redis
docker compose start redis
```

## Security Operations

### Reviewing Authentication Logs

```bash
# Failed login attempts
docker compose logs venio | grep "Authentication failed"

# Successful logins
docker compose logs venio | grep "Login successful"

# Rate limit denials
curl https://api.venio.dev/metrics | grep venio_rate_limit_hits_total
```

### Blocking IP Addresses

**Using firewall:**
```bash
# Block specific IP
sudo ufw deny from 192.168.1.100

# List blocked IPs
sudo ufw status numbered
```

**Application-level (future feature):**
```bash
# Add to blocklist (when implemented)
curl -X POST https://api.venio.dev/api/v1/admin/blocklist \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{"ip": "192.168.1.100", "reason": "Brute force attempt"}'
```

### Rotating Secrets

**JWT Secret:**
```bash
# Generate new secret
openssl rand -base64 32

# Update .env
JWT_SECRET=new_secret_here

# Restart application (invalidates all tokens)
docker compose restart venio
```

**Database Password:**
```sql
-- In PostgreSQL
ALTER USER venio WITH PASSWORD 'new_strong_password';

-- Update .env and restart
```

### Security Audit

**Check for outdated dependencies:**
```bash
cd /home/venio/venio
snyk test
```

**Review user permissions:**
```bash
curl https://api.venio.dev/api/v1/admin/users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  | jq '.[] | select(.is_active == true) | {email, roles}'
```

---

## Quick Reference

### Health Checks

| Endpoint | Purpose |
|----------|---------|
| `/health/live` | Is application running? |
| `/health/ready` | Can application serve requests? |
| `/metrics` | Prometheus metrics |

### Common Commands

```bash
# Check service status
docker compose ps

# Restart all services
docker compose restart

# View logs
docker compose logs -f venio

# Check database health
docker compose exec postgres pg_isready

# Check Redis health
docker compose exec redis redis-cli ping

# Backup database
docker compose exec postgres pg_dump -U venio venio > backup.sql

# Run migrations
docker compose exec venio ./venio migrate up
```

### Emergency Contacts

- **On-Call DevOps:** ops@venio.dev
- **Security Issues:** security@venio.dev
- **Database Issues:** dba@venio.dev

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-15  
**Maintained By:** Operations Team
