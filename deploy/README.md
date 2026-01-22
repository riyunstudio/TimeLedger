# TimeLedger Production Deployment Guide

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- VPS with at least:
  - 2 CPU cores
  - 4GB RAM
  - 20GB storage
- Domain name (optional, for SSL)
- SSH access to VPS

## Quick Start

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd TimeLedger
   ```

2. **Configure environment variables**
   ```bash
   cp .env.production.example .env.production
   # Edit .env.production with your configuration
   ```

3. **Generate SSL certificates (optional)**
   ```bash
   # Using Let's Encrypt with Certbot
   certbot certonly --standalone -d your-domain.com
   cp /etc/letsencrypt/live/your-domain.com/fullchain.pem deploy/ssl/
   cp /etc/letsencrypt/live/your-domain.com/privkey.pem deploy/ssl/
   ```

4. **Deploy**
   ```bash
   # Linux/Mac
   chmod +x deploy/deploy.sh
   ./deploy/deploy.sh

   # Windows
   deploy\deploy.bat
   ```

## Services

The deployment includes the following services:

| Service | Port | Description |
|:---|:---:|:---|
| App | 8080 | Go backend API |
| Nginx | 80, 443 | Reverse proxy & static files |
| MySQL | 3306 | Database |
| Redis | 6379 | Cache & sessions |
| RabbitMQ | 5672, 15672 | Message queue (15672 for management UI) |

## Configuration

### Environment Variables

Key variables in `.env.production`:

| Variable | Description | Default |
|:---|:---|:---|
| `APP_ENV` | Environment mode | `production` |
| `APP_DEBUG` | Debug mode | `false` |
| `MYSQL_PASSWORD` | Database password | *Change this!* |
| `REDIS_PASSWORD` | Redis password | *Change this!* |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | Required |
| `TELEGRAM_CHAT_ID` | Telegram chat ID | Required |
| `JWT_SECRET` | JWT signing secret | *Change this!* |

### SSL/TLS

1. Place certificates in `deploy/ssl/`:
   - `fullchain.pem` - Full certificate chain
   - `privkey.pem` - Private key

2. Uncomment HTTPS configuration in `deploy/nginx/conf.d/timeledger.conf`

3. Restart nginx:
   ```bash
   docker-compose restart nginx
   ```

## Management

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f mysql
docker-compose logs -f nginx
```

### Restart services
```bash
# All services
docker-compose restart

# Specific service
docker-compose restart app
```

### Stop services
```bash
docker-compose down
```

### Update application
```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose up -d --build app
```

### Database access
```bash
# Connect to MySQL
docker-compose exec mysql mysql -u timeledger -p timeledger

# Backup database
docker-compose exec mysql mysqldump -u timeledger -p timeledger > backup.sql

# Restore database
docker-compose exec -T mysql mysql -u timeledger -p timeledger < backup.sql
```

### Redis access
```bash
# Connect to Redis
docker-compose exec redis redis-cli -a REDIS_PASSWORD

# Monitor Redis
docker-compose exec redis redis-cli -a REDIS_PASSWORD monitor
```

## Health Checks

- **App Health**: `http://your-server:8080/health`
- **Nginx Status**: Check container logs
- **MySQL Status**: Check container logs or connect via CLI
- **Redis Status**: Check container logs or run `PING` command
- **RabbitMQ Management**: `http://your-server:15672` (guest/guest)

## Troubleshooting

### Containers won't start
```bash
# Check logs
docker-compose logs

# Check disk space
df -h

# Check memory usage
free -h
```

### Database connection errors
```bash
# Check if MySQL is running
docker-compose ps mysql

# Check MySQL logs
docker-compose logs mysql

# Restart MySQL
docker-compose restart mysql
```

### Redis connection errors
```bash
# Check if Redis is running
docker-compose ps redis

# Check Redis logs
docker-compose logs redis

# Test Redis connection
docker-compose exec redis redis-cli -a REDIS_PASSWORD ping
```

### Permission errors
```bash
# Fix file permissions
chmod -R 755 deploy/
chown -R $(whoami):$(whoami) deploy/
```

## Security Best Practices

1. **Change all default passwords** in `.env.production`
2. **Use strong passwords** (minimum 32 characters)
3. **Enable SSL/TLS** in production
4. **Keep system updated**:
   ```bash
   sudo apt update && sudo apt upgrade
   ```
5. **Configure firewall**:
   ```bash
   sudo ufw allow 22/tcp    # SSH
   sudo ufw allow 80/tcp    # HTTP
   sudo ufw allow 443/tcp   # HTTPS
   sudo ufw enable
   ```
6. **Regular backups**:
   ```bash
   # Add to crontab
   0 2 * * * /path/to/backup-script.sh
   ```

## Backup Strategy

### Automated Backup Script
```bash
#!/bin/bash
# backup.sh
BACKUP_DIR="/path/to/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Backup MySQL
docker-compose exec -T mysql mysqldump -u timeledger -p${MYSQL_PASSWORD} timeledger > ${BACKUP_DIR}/mysql_${DATE}.sql

# Backup Redis
docker-compose exec redis redis-cli -a ${REDIS_PASSWORD} --rdb ${BACKUP_DIR}/redis_${DATE}.rdb

# Keep last 7 days
find ${BACKUP_DIR} -name "*.sql" -mtime +7 -delete
find ${BACKUP_DIR} -name "*.rdb" -mtime +7 -delete
```

## Monitoring

### Log Monitoring
```bash
# Monitor all logs
docker-compose logs -f --tail=100
```

### Resource Monitoring
```bash
# Container stats
docker stats

# Disk usage
docker system df

# Clean up unused resources
docker system prune -a
```

## Performance Tuning

### MySQL
- Tune `innodb_buffer_pool_size` in MySQL config
- Enable query caching if needed
- Use read replicas for scaling

### Redis
- Configure `maxmemory` in Redis config
- Enable persistence (AOF or RDB)
- Use Redis clustering for large datasets

### Nginx
- Enable gzip compression (already enabled)
- Configure caching for static assets
- Use HTTP/2 (already enabled)

## Scaling

### Horizontal Scaling
For multiple app instances:
```yaml
# In docker-compose.yml
app:
  deploy:
    replicas: 3
```

### Load Balancing
Nginx already provides load balancing. Update the upstream block:
```nginx
upstream timeledger_backend {
    least_conn;
    server app1:8080;
    server app2:8080;
    server app3:8080;
}
```

## Support

For issues or questions:
- Check logs: `docker-compose logs`
- Review configuration in `.env.production`
- Consult API documentation at `/swagger/index.html`
- Contact development team
