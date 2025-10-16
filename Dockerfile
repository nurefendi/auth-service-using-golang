# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# IMPORTANT! Disable CGO supaya portable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/app .
COPY .production.env .production.env
COPY .local.env .local.env

# Expose ports for both REST API and gRPC
EXPOSE 9000 9001

# Jalankan binary
CMD ["./app"]
