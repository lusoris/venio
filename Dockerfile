# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/venio cmd/venio/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/worker cmd/worker/main.go

# Runtime Stage
FROM alpine:3.23

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/venio /app/venio
COPY --from=builder /app/worker /app/worker

# Copy config files
COPY configs /app/configs

# Create non-root user
RUN addgroup -g 1000 venio && \
    adduser -D -u 1000 -G venio venio && \
    chown -R venio:venio /app

USER venio

EXPOSE 3690

# Default to venio server (can be overridden to run worker)
ENTRYPOINT ["/app/venio"]
