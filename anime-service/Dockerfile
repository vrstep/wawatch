FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o anime-service .

# Create a smaller final image
FROM alpine:latest

WORKDIR /app

# Install PostgreSQL client for potential migrations
RUN apk --no-cache add postgresql-client

# Copy the binary from builder
COPY --from=builder /app/anime-service .
COPY --from=builder /app/db/migrations ./db/migrations

# Expose the service port
EXPOSE 8081

# Add to backend Dockerfile before CMD
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Install migrate tool
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate

# Command to run
CMD ["./anime-service"]