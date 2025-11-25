# Setup Guide

This guide will help you set up and run the Auth Service on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21 or higher**
  ```bash
  go version
  ```

- **MySQL 8.0 or higher**
  ```bash
  mysql --version
  ```

- **Make** (optional, but recommended)
  ```bash
  make --version
  ```

- **Docker and Docker Compose** (optional, for containerized setup)
  ```bash
  docker --version
  docker-compose --version
  ```

## Setup Methods

You can set up the service in three ways:
1. Local setup with existing MySQL
2. Docker Compose (easiest)
3. Docker with external MySQL

---

## Method 1: Local Setup with Existing MySQL

### Step 1: Install Dependencies

```bash
# Navigate to project directory
cd wiki-service

# Download Go dependencies
go mod download
```

### Step 2: Configure Environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your settings
nano .env  # or use your preferred editor
```

Update the following variables in `.env`:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=auth_service
JWT_SECRET=change_this_to_a_secure_random_string
```

### Step 3: Create Database

```bash
# Login to MySQL
mysql -u root -p

# Create database
CREATE DATABASE auth_service;

# Exit MySQL
exit;
```

### Step 4: Run Migrations

```bash
# Run the migration script
mysql -u root -p auth_service < migrations/001_create_users_table.sql

# Or using Make
make migrate-up DB_HOST=localhost DB_PORT=3306 DB_USER=root DB_PASSWORD=your_password DB_NAME=auth_service
```

### Step 5: Run the Application

```bash
# Run directly with Go
go run cmd/server/main.go

# Or using Make
make run

# Or build and run
make build
./bin/wiki-service
```

The server will start on `http://localhost:8080`

---

## Method 2: Docker Compose (Recommended for Development)

This method sets up both MySQL and the application in containers.

### Step 1: Update Configuration

Edit `docker-compose.yml` if needed to change:
- MySQL credentials
- JWT secret
- Port mappings

### Step 2: Start Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

The application will be available at `http://localhost:8080`

MySQL will be available at `localhost:3306`

### Step 3: Stop Services

```bash
# Stop services
docker-compose down

# Stop and remove volumes (data will be lost)
docker-compose down -v
```

---

## Method 3: Docker with External MySQL

### Step 1: Build Docker Image

```bash
docker build -t wiki-service:latest .
```

### Step 2: Run Container

```bash
docker run -d \
  --name wiki-service \
  -p 8080:8080 \
  -e DB_HOST=your_mysql_host \
  -e DB_PORT=3306 \
  -e DB_USER=your_db_user \
  -e DB_PASSWORD=your_db_password \
  -e DB_NAME=auth_service \
  -e JWT_SECRET=your_jwt_secret \
  -e JWT_EXPIRATION_HOURS=24 \
  wiki-service:latest
```

Or use an env file:

```bash
docker run -d \
  --name wiki-service \
  -p 8080:8080 \
  --env-file .env \
  wiki-service:latest
```

---

## Verification

### 1. Check Health Endpoint

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "wiki-service"
}
```

### 2. Register a Test User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the token from the response.

### 4. Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Troubleshooting

### Database Connection Issues

**Problem:** `Failed to connect to database`

**Solutions:**
1. Check MySQL is running:
   ```bash
   # macOS
   brew services list | grep mysql
   
   # Linux
   systemctl status mysql
   
   # Docker
   docker ps | grep mysql
   ```

2. Verify credentials in `.env`

3. Test connection manually:
   ```bash
   mysql -h localhost -u root -p
   ```

4. Check firewall settings

### Port Already in Use

**Problem:** `bind: address already in use`

**Solutions:**
1. Check what's using port 8080:
   ```bash
   # macOS/Linux
   lsof -i :8080
   
   # Kill the process
   kill -9 PID
   ```

2. Change port in `.env`:
   ```env
   SERVER_PORT=8081
   ```

### Migration Errors

**Problem:** `Table already exists`

**Solution:**
- This is usually fine if running migrations again
- To reset: `make migrate-down` then `make migrate-up`

### Module Not Found

**Problem:** `cannot find module`

**Solution:**
```bash
go mod download
go mod tidy
```

### Docker Build Issues

**Problem:** Build fails or is very slow

**Solutions:**
1. Clean Docker cache:
   ```bash
   docker system prune -a
   ```

2. Rebuild without cache:
   ```bash
   docker-compose build --no-cache
   ```

---

## Development Tips

### Hot Reload (Optional)

Install Air for hot reloading:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Create .air.toml config
air init

# Run with hot reload
air
```

### VS Code Setup

Install recommended extensions:
- Go extension
- REST Client (for api-examples.http)
- Docker extension

`.vscode/settings.json`:
```json
{
  "go.formatTool": "gofmt",
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### Database Management Tools

- **MySQL Workbench**: GUI for MySQL
- **DBeaver**: Universal database tool
- **phpMyAdmin**: Web-based management

### API Testing Tools

- **Postman**: Full-featured API client
- **Insomnia**: Alternative to Postman
- **httpie**: Command-line HTTP client
- **REST Client (VS Code)**: Test APIs from editor

---

## Next Steps

1. Read the [README.md](README.md) for API documentation
2. Check [ARCHITECTURE.md](ARCHITECTURE.md) to understand the codebase
3. Review `api-examples.http` for API examples
4. Start implementing your features!

## Support

If you encounter any issues:
1. Check the troubleshooting section above
2. Review application logs
3. Check MySQL logs
4. Verify environment variables
5. Open an issue on the repository

## Security Notes

⚠️ **Important for Production:**

1. **Change JWT Secret**: Use a strong, random secret
   ```bash
   # Generate a random secret
   openssl rand -base64 32
   ```

2. **Use Strong Passwords**: For database and user accounts

3. **Enable HTTPS**: Use a reverse proxy (nginx, Caddy)

4. **Environment Variables**: Never commit `.env` file

5. **Database Security**: 
   - Don't use root user
   - Create dedicated database user
   - Use strong password
   - Restrict network access

6. **Update Dependencies**: Regularly update packages
   ```bash
   go get -u ./...
   go mod tidy
   ```

