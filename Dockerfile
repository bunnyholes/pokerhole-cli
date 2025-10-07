# Multi-stage Dockerfile for PokerHole Client
# Stage 1: Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /workspace/app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy source code
COPY cmd cmd
COPY internal internal

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o poker cmd/poker-client/main.go

# Stage 2: Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS (if needed) and create non-root user
RUN apk --no-cache add ca-certificates && \
    addgroup -S pokerhole && adduser -S pokerhole -G pokerhole

WORKDIR /app

# Copy binary from builder
COPY --from=builder /workspace/app/poker .

# Change ownership
RUN chown pokerhole:pokerhole /app/poker

# Switch to non-root user
USER pokerhole

# Set environment variable for server URL (can be overridden)
ENV POKERHOLE_SERVER=ws://pokerhole-server:8080/ws/game

# Run the client
ENTRYPOINT ["./poker"]
