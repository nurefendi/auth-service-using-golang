# Pakai official Golang image
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy semua source code
COPY . .

# Build aplikasi
RUN go build -o app

EXPOSE 9100

# Jalankan binary
CMD ["sh", "-c", "./app"]
