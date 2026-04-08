# Dockerfile for powerbi-rest Go application

# Build stage
FROM golang:1.26-alpine AS builder

RUN apk update && apk add --no-cache \
  git \
  ca-certificates \
  tzdata \
  upx \
  && update-ca-certificates

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

RUN upx --best --lzma main || echo "UPX compression failed, continuing without compression";

# Runtime stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy user/group information for security
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/dashboard.html .

ENV PORT=8080

# Expose port 8080
EXPOSE ${PORT}

USER 65534:65534

# Command to run the application
CMD ["./main"]
