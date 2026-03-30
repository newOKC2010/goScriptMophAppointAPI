# Build stage - MOPH Appointment Scheduler
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Production stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set timezone to Bangkok
ENV TZ=Asia/Bangkok

# Copy binary from builder
COPY --from=builder /app/server .

# .env file will be provided via docker-compose volumes

# Run scheduler (no HTTP port needed)
CMD ["./server"]