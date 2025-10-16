# gRPC Implementation for Auth Service

This directory contains the gRPC implementation for the authentication service, providing both REST API and gRPC endpoints.

## Structure

```
grpc/
├── pb/                 # Generated protobuf files
│   ├── auth.pb.go     # Protobuf message definitions
│   └── auth_grpc.pb.go # gRPC service definitions
├── server.go          # gRPC server implementation
proto/
└── auth.proto         # Protocol buffer definition
examples/
├── grpc_client.go     # gRPC client example
└── run_demo.sh        # Demo script
```

## Proto Definition

The `proto/auth.proto` file defines the following services:

### AuthService

- **Register**: Register a new user
- **Login**: User authentication
- **RefreshToken**: Refresh access token
- **Logout**: User logout
- **ChangePassword**: Change user password
- **CheckAccess**: Check user access permissions
- **GetUserProfile**: Get user profile information
- **GetUserFunctions**: Get user function permissions

## Server Configuration

The application now runs both REST API and gRPC servers:

- **REST API**: Port 9000 (default) or `PORT` environment variable
- **gRPC Server**: Port 9001 (default) or `GRPC_PORT` environment variable

## Running the Application

1. **Start the servers**:
   ```bash
   go run main.go
   ```

2. **Test with gRPC client**:
   ```bash
   go run examples/grpc_client.go
   ```

3. **Run complete demo**:
   ```bash
   ./examples/run_demo.sh
   ```

## gRPC Client Example

The client example demonstrates:

1. User registration
2. User login
3. Getting user profile (with JWT authentication)
4. Checking access permissions
5. Getting user functions
6. Token refresh
7. User logout
8. Password change

## Authentication

gRPC calls support JWT authentication via metadata:

```go
md := metadata.Pairs("authorization", "Bearer " + accessToken)
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

## Error Handling

The gRPC server properly handles errors and returns appropriate gRPC status codes:

- `codes.InvalidArgument`: For validation errors
- `codes.Unauthenticated`: For authentication errors
- `codes.PermissionDenied`: For authorization errors
- `codes.Internal`: For internal server errors

## Development Notes

1. The gRPC server reuses the existing business logic from the REST API usecase layer
2. Fiber contexts are created and adapted for gRPC requests to maintain compatibility
3. JWT token parsing and validation work the same way as REST API
4. All existing middleware and validation logic is preserved

## Building Proto Files

To regenerate the protobuf files after modifying `auth.proto`:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/auth.proto

# Move generated files to proper location
mv proto/*.pb.go grpc/pb/
```

## Testing

You can test the gRPC service using:

1. The provided Go client example
2. gRPC testing tools like [grpcurl](https://github.com/fullstorydev/grpcurl)
3. GUI tools like [BloomRPC](https://github.com/uw-labs/bloomrpc) or [Postman](https://www.postman.com/)

Example with grpcurl:
```bash
# List services
grpcurl -plaintext localhost:9001 list

# Call register method
grpcurl -plaintext -d '{
  "full_name": "John Doe",
  "email": "john@example.com", 
  "password": "password123",
  "gender": 1
}' localhost:9001 auth.AuthService/Register
```