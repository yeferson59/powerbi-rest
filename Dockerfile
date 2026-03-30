# Dockerfile for powerbi-rest Go application

# Build stage
FROM golang:1.26-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

ENV PORT=8080

# Expose port 8080
EXPOSE ${PORT}

# Command to run the application
CMD ["./main"]
