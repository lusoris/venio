# Deployment Guide

## Docker Compose (Recommended for most users)

### Quick Start

```bash
# Download docker-compose.yml
curl -sL https://raw.githubusercontent.com/lusoris/venio/main/docker-compose.yml -o docker-compose.yml

# Create .env
cat > .env << EOF
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_redis_password
TYPESENSE_API_KEY=your_typesense_key
JWT_SECRET=your_jwt_secret_minimum_32_chars
API_KEY=your_api_key
EOF

# Start
docker compose up -d

# Check logs
docker compose logs -f venio
```

Access Venio at `http://localhost:3690`

### Production Deployment

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  venio:
    image: ghcr.io/lusoris/venio:latest
    restart: unless-stopped
    ports:
      - "3690:3690"
    env_file:
      - .env
    depends_on:
      - postgres
      - redis
    networks:
      - venio
    
  # ... other services

volumes:
  postgres-data:
  redis-data:

networks:
  venio:
```

### Reverse Proxy (Nginx)

```nginx
server {
    listen 80;
    server_name venio.example.com;
    
    location / {
        proxy_pass http://localhost:3690;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

With SSL (Certbot):
```bash
sudo certbot --nginx -d venio.example.com
```

## Kubernetes

### Helm Chart (Coming Soon)

```bash
helm repo add venio https://charts.venio.io
helm install venio venio/venio
```

### Manual Deployment

```yaml
# venio-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: venio
spec:
  replicas: 2
  selector:
    matchLabels:
      app: venio
  template:
    metadata:
      labels:
        app: venio
    spec:
      containers:
      - name: venio
        image: ghcr.io/lusoris/venio:latest
        ports:
        - containerPort: 3690
        env:
        - name: POSTGRES_HOST
          value: postgres-service
        envFrom:
        - secretRef:
            name: venio-secrets
---
apiVersion: v1
kind: Service
metadata:
  name: venio-service
spec:
  selector:
    app: venio
  ports:
  - port: 80
    targetPort: 3690
```

Apply:
```bash
kubectl apply -f venio-deployment.yaml
```

## Security Hardening

### 1. Use Non-Root User

Already done in Dockerfile - runs as user `venio` (UID 1000)

### 2. Limit Resources

```yaml
# docker-compose.yml
services:
  venio:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### 3. Network Isolation

```yaml
networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true  # No internet access
```

### 4. Enable TLS

Always use HTTPS in production via reverse proxy.

### 5. Secrets Management

Use Docker secrets or Kubernetes secrets, not .env files:

```bash
echo "my_secret" | docker secret create postgres_password -
```

## Backup Strategy

### Database Backups

```bash
# PostgreSQL
docker exec venio-postgres pg_dump -U venio venio | gzip > backup-$(date +%Y%m%d).sql.gz

# Automated with cron
0 2 * * * /path/to/backup-script.sh
```

### Configuration Backups

```bash
# Backup .env and configs
tar -czf venio-config-$(date +%Y%m%d).tar.gz .env configs/
```

## Monitoring

### Health Checks

```bash
# Check if Venio is healthy
curl http://localhost:3690/health

# Expected response
{"status":"ok","version":"v1.0.0"}
```

### Logs

```bash
# Docker Compose
docker compose logs -f venio

# Kubernetes
kubectl logs -f deployment/venio
```

### Metrics (Future)

Prometheus endpoint: `http://localhost:3690/metrics`

## Upgrading

### Docker Compose

```bash
# Pull latest image
docker compose pull

# Restart with new image
docker compose up -d

# Check logs
docker compose logs -f venio
```

### Kubernetes

```bash
kubectl set image deployment/venio venio=ghcr.io/lusoris/venio:v1.1.0
kubectl rollout status deployment/venio
```

## Rollback

### Docker Compose

```bash
# Use specific version
docker compose down
# Edit docker-compose.yml to use old version
docker compose up -d
```

### Kubernetes

```bash
kubectl rollout undo deployment/venio
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose logs venio

# Common issues:
# - Database not ready
# - Missing environment variables
# - Port already in use
```

### Database Connection Issues

```bash
# Test database connection
docker exec venio-postgres psql -U venio -c "SELECT 1"

# Check network
docker compose exec venio ping postgres
```

### High Memory Usage

```bash
# Check stats
docker stats venio

# Adjust limits in docker-compose.yml
```

## Performance Tuning

### PostgreSQL

```bash
# Increase connection pool
POSTGRES_MAX_CONNECTIONS=200
```

### Redis

```bash
# Increase memory limit
REDIS_MAXMEMORY=512mb
```

### Application

```bash
# Adjust worker concurrency
WORKER_CONCURRENCY=10
```

---

## Document Revision

| Version | Date | Author | Changes | Source Version |
|---------|------|--------|---------|----------------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial deployment guide | - |

## Referenced Documentation

- **Docker Compose:** [Docker Compose v2 Documentation](https://docs.docker.com/compose/) (v2.29.0)
- **Docker Best Practices:** [Docker Production Best Practices](https://docs.docker.com/engine/install/)
- **PostgreSQL Docker:** [PostgreSQL Docker Image](https://hub.docker.com/_/postgres) (18.1-alpine)
- **Redis Docker:** [Redis Docker Image](https://hub.docker.com/_/redis) (8.4-alpine)
- **Kubernetes Deployment:** [Kubernetes Documentation](https://kubernetes.io/docs/home/) (v1.31)
- **Security Hardening:** [Docker Security Best Practices](https://docs.docker.com/engine/security/)

*More deployment guides coming soon*
