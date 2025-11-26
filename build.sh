#!/bin/bash

# Wiki Service - Build and Deploy Script

set -e

echo "🏗️  Building Wiki Service..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

# Build the application
print_status "Building Docker image..."
docker build -t wiki-service:latest .

if [ $? -eq 0 ]; then
    print_success "Docker image built successfully!"
else
    print_error "Failed to build Docker image"
    exit 1
fi

# Option to run with docker-compose
if [ "$1" = "--compose" ] || [ "$1" = "-c" ]; then
    print_status "Starting services with docker-compose..."
    docker-compose up -d

    print_success "Services started!"
    print_status "Wiki Service: http://localhost:8080"
    print_status "MongoDB: localhost:27017"
    print_status "Redis: localhost:6379"
    print_status "Consul UI: http://localhost:8500"
fi

# Option to run standalone
if [ "$1" = "--standalone" ] || [ "$1" = "-s" ]; then
    print_warning "Make sure MongoDB is running on localhost:27017"

    print_status "Running standalone container..."
    docker run -d \
        --name wiki-service-standalone \
        -p 8080:8080 \
        -e MONGO_HOST=host.docker.internal \
        -e MONGO_PORT=27017 \
        -e MONGO_DB_NAME=services_management \
        wiki-service:latest

    print_success "Container started!"
    print_status "Wiki Service: http://localhost:8080"
fi

# Show usage if no arguments
if [ $# -eq 0 ]; then
    echo ""
    print_status "Usage:"
    echo "  ./build.sh --compose    # Build and run with docker-compose (full stack)"
    echo "  ./build.sh --standalone # Build and run standalone (needs external MongoDB)"
    echo ""
    print_success "Docker image 'wiki-service:latest' built successfully!"
    echo "Run './build.sh --compose' to start the full application stack."
fi
