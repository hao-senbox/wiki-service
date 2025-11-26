# Wiki Service - Docker Build & Deploy

## 📋 Tổng quan

Hướng dẫn build và deploy Wiki Service sử dụng Docker và Docker Compose.

## 🚀 Quick Start

### 1. Build và chạy full stack
```bash
chmod +x build.sh
./build.sh --compose
```

### 2. Kiểm tra services
- **Wiki Service**: http://localhost:8080
- **MongoDB**: localhost:27017
- **Redis**: localhost:6379
- **Consul UI**: http://localhost:8500

## 🏗️ Build Options

### Full Stack (Khuyến nghị)
```bash
./build.sh --compose
```
Khởi động tất cả services: Wiki Service + MongoDB + Redis + Consul

### Standalone
```bash
./build.sh --standalone
```
Chỉ build và chạy Wiki Service (cần MongoDB external)

### Chỉ build image
```bash
docker build -t wiki-service:latest .
```

## 📁 File Structure

```
.
├── Dockerfile              # Multi-stage build với Go 1.25.3
├── docker-compose.yml      # Full stack configuration
├── build.sh               # Build script với options
├── .dockerignore          # Exclude files khỏi build context
├── env.example            # Environment variables template
└── DOCKER_README.md       # This file
```

## 🔧 Configuration

### Environment Variables
Copy `env.example` thành `.env` và chỉnh sửa:

```bash
cp env.example .env
```

### Key Settings
```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
MONGO_HOST=mongodb
MONGO_PORT=27017
MONGO_DB_NAME=services_management

# Cache & Discovery
REDIS_HOST=redis
CONSUL_HOST=consul
```

## 🐳 Dockerfile Features

### Build Stage
- **Go 1.25.3**: Match với go.mod
- **Alpine Linux**: Small base image
- **Multi-stage**: Tối ưu image size
- **Optimized build**: Static linking, stripped binary

### Runtime Stage
- **Non-root user**: Security best practice
- **Minimal Alpine**: Chỉ cần thiết dependencies
- **Health check**: Built-in monitoring
- **Proper permissions**: Secure file ownership

## 🚢 Docker Compose Services

### wiki-service
- **Port**: 8080
- **Health check**: HTTP endpoint
- **Dependencies**: MongoDB healthy
- **Environment**: Full configuration

### mongodb
- **Version**: 7.0
- **Authentication**: Root user configured
- **Health check**: MongoDB ping
- **Volume**: Persistent data

### redis
- **Version**: 7.2-alpine
- **Append-only**: Data persistence
- **Health check**: Redis ping

### consul
- **Version**: 1.15
- **UI enabled**: Web interface
- **Bootstrap mode**: Single node setup

## 🔍 Troubleshooting

### Build Issues
```bash
# Check Go version
go version

# Clean build cache
docker system prune -f

# Rebuild without cache
docker build --no-cache -t wiki-service:latest .
```

### Runtime Issues
```bash
# Check container logs
docker-compose logs wiki-service

# Check container status
docker-compose ps

# Restart services
docker-compose restart
```

### Database Connection
```bash
# Test MongoDB connection
docker exec -it wiki-service-mongodb mongosh -u root -p rootpassword services_management

# Check Redis
docker exec -it wiki-service-redis redis-cli ping
```

## 📊 Monitoring

### Health Checks
- **Wiki Service**: `http://localhost:8080/health`
- **MongoDB**: Container health check
- **Redis**: Container health check
- **Consul**: `http://localhost:8500`

### Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f wiki-service

# Last 100 lines
docker-compose logs --tail=100 wiki-service
```

## 🔄 Updates & Deployment

### Update Application
```bash
# Rebuild and restart
docker-compose build --no-cache wiki-service
docker-compose up -d wiki-service
```

### Scale Services
```bash
# Scale wiki-service to 3 instances
docker-compose up -d --scale wiki-service=3
```

### Backup Data
```bash
# Backup MongoDB data
docker run --rm -v wiki-service_mongo_data:/data -v $(pwd):/backup alpine tar czf /backup/mongo_backup.tar.gz -C /data .
```

## 🧪 Testing

Sau khi deploy, test với Postman collection:
- Import `postman_requests_data.json`
- Set token trong environment
- Run từng request theo thứ tự

## 🔒 Security Notes

- **Non-root user**: Application chạy với user `appuser`
- **Minimal attack surface**: Chỉ cần thiết packages
- **No secrets in image**: Config qua environment variables
- **Health checks**: Automatic monitoring và restart

## 📈 Performance

- **Go binary**: Optimized với `-ldflags='-w -s'`
- **Alpine Linux**: Small runtime image (~10MB)
- **Multi-stage build**: Build artifacts không vào final image
- **Static linking**: No external dependencies

## 🆘 Support

Nếu gặp vấn đề:
1. Check logs: `docker-compose logs`
2. Verify configuration: `docker-compose config`
3. Test network: `docker network ls`
4. Restart stack: `docker-compose down && docker-compose up -d`
