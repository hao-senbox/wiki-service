.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building wiki-service..."
	@go build -o bin/wiki-service cmd/server/main.go

run: ## Run the application
	@echo "Running wiki-service..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) < migrations/001_create_users_table.sql
	@echo "Migrations completed successfully!"

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) -e "DROP TABLE IF EXISTS users;"

logs: ## View application logs
	@tail -f logs/app.log

audit-logs: ## View audit logs
	@tail -f logs/audit.log

clean-logs: ## Clean old log files
	@echo "Cleaning logs..."
	@rm -rf logs/*.log

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t wiki-service:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env wiki-service:latest

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

dev: ## Run in development mode with hot reload (requires air)
	@echo "Running in development mode..."
	@air

.DEFAULT_GOAL := help

