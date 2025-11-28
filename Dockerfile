# ----------------------------
# Stage 1: Build binary
# ----------------------------
    FROM golang:1.25-alpine AS builder

    WORKDIR /app
    
    # Copy go.mod/go.sum và tải dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy toàn bộ source
    COPY . .
    
    # Build Go binary
    RUN go build -o wiki-service cmd/server/main.go
    
    # ----------------------------
    # Stage 2: Runtime
    # ----------------------------
    FROM alpine:3.21
    
    WORKDIR /root
    
    # Cài libc6-compat để chạy binary Go static
    RUN apk add --no-cache libc6-compat
    
    # Copy binary từ builder
    COPY --from=builder /app/wiki-service .
    
    # Expose port
    EXPOSE 8023
    
    # Chạy trực tiếp binary
    CMD ["./wiki-service"]
    