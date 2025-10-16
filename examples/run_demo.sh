#!/bin/bash

echo "Starting Auth Service gRPC Demo..."

# Build the main application
echo "Building auth service..."
go build -o auth-service

# Build the gRPC client
echo "Building gRPC client..."
go build -o grpc-client examples/grpc_client.go

# Start the auth service in background
echo "Starting auth service on port 9000 (REST) and 9001 (gRPC)..."
./auth-service &
AUTH_PID=$!

# Wait a bit for the server to start
sleep 3

# Run the gRPC client
echo "Running gRPC client tests..."
./grpc-client

# Cleanup
echo "Stopping auth service..."
kill $AUTH_PID

echo "Demo completed!"