# Build stage
FROM golang:1.22.1-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/server .
COPY --from=builder /app/migrations ./migrations

# Install migrate for database migrations
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    rm -rf migrate.linux-amd64.tar.gz

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
