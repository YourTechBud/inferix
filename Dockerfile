# Stage 1: Build the Go binary
FROM golang:1.23.2 AS builder

# Set environment variables to enforce better build security practices
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Create and set working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . ./

# Build the Go application
RUN go build -a -ldflags="-s -w" -o /app/inferix

# Stage 2: Create a minimal container to run the Go binary
FROM debian:12.8-slim AS runtime

# Create a non-root user and group
RUN groupadd --gid 1001 appuser && \
    useradd --uid 1001 --gid appuser --shell /bin/bash --create-home appuser

# Set working directory
WORKDIR /home/appuser

# Copy the compiled Go binary from the builder image and set ownership
COPY --from=builder /app/inferix ./inferix
RUN chown appuser:appuser ./inferix

# Set permissions on the binary to allow execution by the new user
RUN chmod +x ./inferix

# Switch to the non-root user
USER appuser

# Expose application port
EXPOSE 4386

# Command to run the binary
ENTRYPOINT ["./inferix"]
CMD ["--config-path", "/config/inferix.yaml"]