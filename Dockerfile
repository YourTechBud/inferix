# Stage 1: Build the Go binary
FROM golang:1.23.2 AS builder

# Set environment variables to enforce better build security practices
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
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
RUN go build -ldflags="-s -w" -o /app/inferix

# Stage 2: Create a minimal container to run the Go binary
FROM gcr.io/distroless/static-debian12 AS runtime

# Set non-root user for running the application
USER nonroot:nonroot

# Copy the compiled Go binary from the builder image
COPY --from=builder /app/inferix /inferix

# Expose application port
EXPOSE 4386

# Command to run the binary
CMD ["/inferix", "--config", "/config/config.yaml"]