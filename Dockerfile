# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/etl ./cmd/etl

# Final stage
FROM alpine:3.17

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/etl .

# Create necessary directories
RUN mkdir -p data/raw data/processed logs

# Expose metrics port
EXPOSE 2112

# Run the application
CMD ["./etl"] 