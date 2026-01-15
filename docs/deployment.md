# Production Deployment Guide

This guide covers deploying Venio to production with full observability stack.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Infrastructure Setup](#infrastructure-setup)
- [Configuration](#configuration)
- [Deployment](#deployment)
- [Observability](#observability)
- [Security](#security)
- [Maintenance](#maintenance)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- Docker 24.0+ and Docker Compose 2.20+
- Linux server (Ubuntu 22.04 LTS recommended)
- 4GB RAM minimum (8GB recommended)
- 20GB disk space minimum
- Domain name with DNS configured
- SSL certificate (Let's Encrypt recommended)

### Required Services

- PostgreSQL 18.1+ (managed service or self-hosted)
- Redis 8.4+ (managed service or self-hosted)
- SMTP server (for alerts)

## Infrastructure Setup

### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create deployment user
sudo useradd -m -s /bin/bash venio
sudo usermod -aG docker venio

# Switch to deployment user
sudo su - venio
```

### 2. Firewall Configuration

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
```

### 3. Directory Structure

```bash
mkdir -p ~/venio/{config,data,logs,backups}
cd ~/venio
```

## Configuration

### 1. Environment Variables

Create `.env` file:

```bash
# Application
APP_ENV=production
APP_NAME=Venio
APP_VERSION=1.0.0
SERVER_PORT=3690

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=venio
DB_PASSWORD=<strong-password-here>
DB_NAME=venio
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_LIFETIME=15m

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=<strong-password-here>
REDIS_DB=0

# JWT
JWT_SECRET=<generate-with-openssl-rand-base64-32>
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com

# Rate Limiting
RATE_LIMIT_AUTH_REQUESTS=5
RATE_LIMIT_AUTH_WINDOW=1m
RATE_LIMIT_GENERAL_REQUESTS=100
RATE_LIMIT_GENERAL_WINDOW=1m

# Observability
GRAFANA_USER=admin
GRAFANA_PASSWORD=<strong-password-here>
```

### 2. Generate Secrets

```bash
# JWT secret (keep this secure!)
openssl rand -base64 32

# Database password
openssl rand -base64 32

# Redis password
openssl rand -base64 24
```

### 3. SSL Configuration

#### Using Let's Encrypt

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Generate certificate
sudo certbot certonly --standalone -d api.yourdomain.com

# Certificates will be at:
# /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem
# /etc/letsencrypt/live/api.yourdomain.com/privkey.pem
```

## Deployment

### 1. Clone Repository

```bash
git clone https://github.com/lusoris/venio.git
cd venio
```

### 2. Build Production Image

```bash
docker build -t venio:latest -f Dockerfile .
```

### 3. Start Services

```bash
# Start database and Redis first
docker compose up -d postgres redis

# Wait for database to be ready
sleep 10

# Run migrations
docker compose run --rm venio ./bin/venio migrate up

# Start all services
docker compose up -d
```

### 4. Verify Deployment

```bash
# Check service health
curl http://localhost:3690/health/ready

# Check logs
docker compose logs -f venio

# Check metrics
curl http://localhost:3690/metrics
```

## Observability

### Accessing Observability Tools

| Service | URL | Purpose |
|---------|-----|---------|
| Grafana | http://your-ip:3001 | Dashboards and visualization |
| Prometheus | http://your-ip:9090 | Metrics storage and querying |
| Alertmanager | http://your-ip:9093 | Alert management |

### Configuring Alerts

Edit `deployments/alertmanager/config.yml` to add notification channels:

#### Email Notifications

```yaml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alerts@yourdomain.com'
  smtp_auth_username: 'alerts@yourdomain.com'
  smtp_auth_password: 'your-app-password'

receivers:
  - name: 'critical'
    email_configs:
      - to: 'ops-team@yourdomain.com'
        headers:
          Subject: '[CRITICAL] Venio Alert'
```

#### Slack Notifications

```yaml
receivers:
  - name: 'critical'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
        channel: '#venio-alerts'
        title: '[CRITICAL] {{ .GroupLabels.alertname }}'
```

#### PagerDuty

```yaml
receivers:
  - name: 'critical'
    pagerduty_configs:
      - service_key: 'YOUR_PAGERDUTY_SERVICE_KEY'
```

### Grafana Setup

1. Access Grafana at `http://your-ip:3001`
2. Login with credentials from `.env`
3. Dashboards are auto-loaded from `deployments/grafana/dashboards/`
4. Customize dashboards as needed

## Security

### 1. Reverse Proxy (Nginx)

Create `/etc/nginx/sites-available/venio`:

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    location / {
        proxy_pass http://localhost:3690;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/venio /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 2. Database Security

```sql
-- Restrict database user permissions
REVOKE ALL ON DATABASE venio FROM PUBLIC;
GRANT CONNECT ON DATABASE venio TO venio;
GRANT ALL ON ALL TABLES IN SCHEMA public TO venio;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO venio;

-- Enable SSL connections only
ALTER SYSTEM SET ssl = on;
```

### 3. Redis Security

Edit redis configuration:

```conf
# Require password
requirepass your-strong-password

# Disable dangerous commands
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command CONFIG ""
```

### 4. Application Security

- Use strong JWT secret (32+ characters)
- Enable HTTPS only in production
- Set secure CORS origins
- Implement rate limiting (already configured)
- Regular security audits with Snyk

## Maintenance

### Backups

#### Database Backup

```bash
#!/bin/bash
# save as backup-db.sh

BACKUP_DIR=~/venio/backups
DATE=$(date +%Y%m%d_%H%M%S)

docker compose exec -T postgres pg_dump -U venio venio > $BACKUP_DIR/venio_$DATE.sql
gzip $BACKUP_DIR/venio_$DATE.sql

# Keep last 7 days
find $BACKUP_DIR -name "venio_*.sql.gz" -mtime +7 -delete
```

Schedule with cron:

```bash
crontab -e
# Add: 0 2 * * * /home/venio/backup-db.sh
```

#### Redis Backup

Redis automatically creates snapshots in `/data` volume.

### Updates

```bash
# Pull latest code
git pull origin main

# Rebuild image
docker build -t venio:latest -f Dockerfile .

# Run migrations
docker compose run --rm venio ./bin/venio migrate up

# Rolling update
docker compose up -d --no-deps --build venio

# Verify
curl https://api.yourdomain.com/health/ready
```

### Monitoring

Check metrics regularly:

```bash
# Application health
curl https://api.yourdomain.com/health/ready

# Prometheus targets
curl http://localhost:9090/api/v1/targets

# View Grafana dashboards
open http://localhost:3001
```

## Troubleshooting

### API Not Responding

```bash
# Check service status
docker compose ps

# View logs
docker compose logs -f venio

# Check health
curl http://localhost:3690/health/ready
```

### Database Connection Issues

```bash
# Test database connectivity
docker compose exec postgres psql -U venio -d venio -c "SELECT 1;"

# Check connection pool metrics
curl http://localhost:3690/metrics | grep venio_db_connections
```

### Redis Connection Issues

```bash
# Test Redis connectivity
docker compose exec redis redis-cli -a your-password ping

# Check Redis metrics
curl http://localhost:3690/metrics | grep venio_redis
```

### High Memory Usage

```bash
# Check container stats
docker stats

# Restart service if needed
docker compose restart venio
```

### Certificate Renewal

```bash
# Test renewal
sudo certbot renew --dry-run

# Renew certificates (automatic via cron, but can be manual)
sudo certbot renew
sudo systemctl reload nginx
```

## Scaling

### Horizontal Scaling

For high traffic, run multiple API instances:

```yaml
# docker-compose.yml
services:
  venio:
    deploy:
      replicas: 3
```

Add load balancer (Nginx):

```nginx
upstream venio_backend {
    least_conn;
    server venio-1:3690;
    server venio-2:3690;
    server venio-3:3690;
}
```

### Vertical Scaling

Increase container resources:

```yaml
services:
  venio:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

## Kubernetes Deployment

For Kubernetes, see [docs/kubernetes.md](./kubernetes.md).

## Support

- Documentation: https://docs.venio.dev
- Issues: https://github.com/lusoris/venio/issues
- Email: support@venio.dev
