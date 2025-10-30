FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the exporter
WORKDIR /build/cmd/mm-exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o mm-exporter .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates dbus

# Create non-root user (note: exporter needs root for D-Bus access in practice)
RUN addgroup -S mmexporter && adduser -S mmexporter -G mmexporter

# Copy binary from builder
COPY --from=builder /build/cmd/mm-exporter/mm-exporter /usr/local/bin/mm-exporter

# Make it executable
RUN chmod +x /usr/local/bin/mm-exporter

# Expose metrics port
EXPOSE 9539

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9539/health || exit 1

# Run as root (required for D-Bus system bus access)
USER root

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/mm-exporter"]

# Default arguments
CMD ["-listen-address=:9539", "-signal-rate=5s"]
