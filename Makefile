# Makefile for Auth Service

.PHONY: build run test clean proto grpc-client demo help

# Build the main application
build:
	@echo "Building auth service..."
	@go build -o auth-service

# Run the application (both REST and gRPC)
run:
	@echo "Starting auth service..."
	@echo "REST API will be available on :9000"
	@echo "gRPC server will be available on :9001"
	@go run main.go

# Build gRPC client
grpc-client:
	@echo "Building gRPC client..."
	@go build -o grpc-client examples/grpc_client.go

# Run gRPC client test
test-grpc:
	@echo "Running gRPC client test..."
	@go run examples/grpc_client.go

# Run the complete demo
demo:
	@echo "Running complete demo..."
	@chmod +x examples/run_demo.sh
	@./examples/run_demo.sh

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto
	@mv proto/*.pb.go grpc/pb/
	@echo "Protobuf files generated successfully"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f auth-service grpc-client
	@echo "Cleaned successfully"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Install protoc tools
install-protoc-tools:
	@echo "Installing protoc tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t auth-service .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 9000:9000 -p 9001:9001 auth-service

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the auth service binary"
	@echo "  run                - Run the auth service (REST + gRPC)"
	@echo "  grpc-client        - Build the gRPC client"
	@echo "  test-grpc          - Run gRPC client test"
	@echo "  demo               - Run complete demo"
	@echo "  proto              - Generate protobuf files"
	@echo "  deps               - Install dependencies"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  install-protoc-tools - Install protoc tools"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  help               - Show this help message"