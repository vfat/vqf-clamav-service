# ==============================================================================
# STAGE 1: Build Go Binary
# ==============================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build tools if CGO or sqlite requires it
RUN apk add --no-cache git gcc musl-dev

# Copy module files first for dependency caching
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source code
COPY . .

# Build statically-linked binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w -extldflags '-static'" -o clamav-service ./cmd/server

# ==============================================================================
# STAGE 2: Final All-in-One Runtime Container
# ==============================================================================
FROM alpine:3.20

# Install ClamAV daemon, Freshclam, and essential utilities
RUN apk add --no-cache \
    clamav \
    clamav-daemon \
    ca-certificates \
    curl \
    tzdata

# Create directory structure and set ownership
RUN mkdir -p /app /data /data/quarantine /var/run/clamav /var/log/clamav /var/lib/clamav && \
    chown -R clamav:clamav /var/run/clamav /var/log/clamav /var/lib/clamav /data

# Copy custom ClamAV configs
COPY configs/clamd.conf /etc/clamav/clamd.conf
COPY configs/freshclam.conf /etc/clamav/freshclam.conf

# Copy compiled binary from builder stage
COPY --from=builder /build/clamav-service /app/clamav-service

# Expose HTTP API & Web UI port
EXPOSE 8080

# Persistent volume for SQLite DB & Quarantine Vault
VOLUME ["/data"]

WORKDIR /app

# Run Go Native Process Supervisor as PID 1
ENTRYPOINT ["/app/clamav-service"]
CMD ["run"]
