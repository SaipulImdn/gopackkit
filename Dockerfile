# Multi-stage Dockerfile for GoPackKit

# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache ca-certificates git tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# Test stage
FROM builder AS tester

# Run tests
RUN go test -v ./...

# Run linting
RUN go vet ./...

# Security scan
RUN go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
RUN gosec ./...

# Final stage - minimal image
FROM scratch AS final

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/main /main

# Set timezone
ENV TZ=UTC

# Expose port (if needed for examples)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/main"]

# Development stage for testing
FROM golang:1.21-alpine AS development

# Install development tools and Go tools
RUN apk add --no-cache bash curl git make && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest && \
    go install golang.org/x/vuln/cmd/govulncheck@latest

# Set working directory
WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Default command for development
CMD ["bash"]
