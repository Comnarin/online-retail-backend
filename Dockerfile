# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/app/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates and tzdata
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations if any
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 4001

# Command to run
CMD ["./main"]
