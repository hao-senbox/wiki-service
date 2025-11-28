# Build stage - không hardcode GOARCH
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build cho kiến trúc hiện tại, không hardcode amd64
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o wiki-service \
    ./cmd/server/main.go

# Verify binary
RUN file wiki-service && chmod +x wiki-service

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app

WORKDIR /app

COPY --from=builder /build/wiki-service .
COPY --from=builder /build/.env* ./ || true

RUN chmod +x wiki-service

USER appuser

EXPOSE 8023

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8023/health || exit 1

CMD ["./wiki-service"]