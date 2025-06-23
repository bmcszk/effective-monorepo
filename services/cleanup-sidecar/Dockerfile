# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the cleanup sidecar binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cleanup-sidecar ./services/cleanup-sidecar

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/cleanup-sidecar .

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

USER appuser

# Command to run
CMD ["./cleanup-sidecar"]