# Build stage
FROM golang:1.23.4-alpine AS builder

# Set working directory for the build
WORKDIR /build

# Install git, swag tool, and other necessary tools
RUN apk add --no-cache git postgresql-client curl bash

# Install swag tool using Go get
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod and go.sum to leverage Docker layer caching
COPY go.mod go.sum ./

# Generate go.sum if it does not exist
RUN go mod tidy

# Copy the entire source code
COPY . .

# Generate Swagger documentation using swag
RUN swag init -g src/cmd/api/main.go # Ensure the path to main.go is correct

# List the contents of the docs directory for verification
RUN ls -al /build/docs
# Build the application binary in /build directory
RUN go clean -modcache
RUN go build -o /build/main ./src/cmd/api

# Run stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install runtime dependencies, including PostgreSQL client, bash, and ca-certificates
RUN apk add --no-cache ca-certificates postgresql-client bash

# Copy the built binary from the builder stage to /app
COPY --from=builder /build/main /app/main

# Copy the entrypoint script
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

# Ensure proper file permissions
RUN chmod +x /app/docker-entrypoint.sh \
    && chmod 755 /app/docker-entrypoint.sh

# Copy the generated Swagger docs to /app/docs
# Ensure that swagger.json and swagger.yaml are available for the Swagger UI
COPY --from=builder /build/docs /app/docs
# Expose the application port
EXPOSE 8085

# Set the entrypoint to run the entrypoint script
ENTRYPOINT ["/bin/sh", "/app/docker-entrypoint.sh"]
