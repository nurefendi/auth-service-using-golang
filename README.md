# Auth Service with Golang, Fiber, GORM, MariaDB, and gRPC

## Description
This authentication service is built using Golang with the Fiber framework, utilizing GORM as the ORM, and MariaDB as the primary database. The service provides both REST API and gRPC endpoints for maximum flexibility.

## Features
- User registration
- User login with JWT
- Authentication middleware
- Refresh token
- User management
- **gRPC API** - Full gRPC implementation alongside REST API
- Dual server support (REST + gRPC)
- JWT authentication for both REST and gRPC

## Prerequisites
Ensure you have installed:
- [Golang](https://go.dev/dl/)
- [MariaDB](https://mariadb.org/download/)
- [Docker](https://www.docker.com/) (Optional)

## Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/nurefendi/auth-service.git
    cd auth-service
    ```
2. Create a `.local.env` or `.production.env` file and add the database configuration:
    ```sh
    DB_USER=root
    DB_PASSWORD=yourpassword
    DB_HOST=localhost
    DB_PORT=3306
    DB_NAME=auth_db
    JWT_SECRET=your_secret_key
    ```
3. Run MariaDB database:
    ```sh
    docker run --name mariadb -e MYSQL_ROOT_PASSWORD=yourpassword -e MYSQL_DATABASE=auth_db -p 3306:3306 -d mariadb
    ```
4. Install dependencies:
    ```sh
    go mod tidy
    ```
5. Run database migrations:
    ```sh
    go run main.go migrate
    ```
6. Start the application:
    ```sh
    go run main.go
    ```

## API Endpoints
<!-- ### 1. Registration
- **Endpoint**: `POST /register`
- **Request Body**:
    ```json
    {
      "name": "John Doe",
      "email": "johndoe@example.com",
      "password": "password123"
    }
    ```
- **Response**:
    ```json
    {
      "message": "User registered successfully"
    }
    ```

### 2. Login
- **Endpoint**: `POST /login`
- **Request Body**:
    ```json
    {
      "email": "johndoe@example.com",
      "password": "password123"
    }
    ```
- **Response**:
    ```json
    {
      "token": "your_jwt_token",
      "refresh_token": "your_refresh_token"
    }
    ```

### 3. Authentication Middleware
Add JWT middleware to protect endpoints:
```go
app.Use(jwtware.New(jwtware.Config{
    SigningKey: []byte(os.Getenv("JWT_SECRET")),
}))
``` -->

<!-- ### 4. User Profile
- **Endpoint**: `GET /profile`
- **Header**:
    ```sh
    Authorization: Bearer your_jwt_token
    ```
- **Response**:
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "email": "johndoe@example.com"
    }
    ``` -->

## gRPC Support

This service also provides a full gRPC API alongside the REST API. Both servers run simultaneously:

- **REST API**: Port 9000 (default) or `PORT` environment variable
- **gRPC Server**: Port 9001 (default) or `GRPC_PORT` environment variable

### gRPC Services

The gRPC API provides the following services:

- **Register**: User registration
- **Login**: User authentication
- **RefreshToken**: Refresh access token
- **Logout**: User logout
- **ChangePassword**: Change user password
- **CheckAccess**: Check user access permissions
- **GetUserProfile**: Get user profile information
- **GetUserFunctions**: Get user function permissions

### Testing gRPC

1. **Run the gRPC client example**:
   ```sh
   go run examples/grpc_client.go
   ```

2. **Run the complete demo**:
   ```sh
   ./examples/run_demo.sh
   ```

3. **Using grpcurl** (if installed):
   ```sh
   # List available services
   grpcurl -plaintext localhost:9001 list
   
   # Test registration
   grpcurl -plaintext -d '{
     "full_name": "John Doe",
     "email": "john@example.com",
     "password": "password123",
     "gender": 1
   }' localhost:9001 auth.AuthService/Register
   ```

### gRPC Authentication

gRPC calls support JWT authentication via metadata:
```go
md := metadata.Pairs("authorization", "Bearer " + accessToken)
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

For more details about the gRPC implementation, see [grpc/README.md](grpc/README.md).

<!-- ## Project Structure
```
/auth-service
│── main.go
│── config/
│   ├── database.go
│── models/
│   ├── user.go
│── routes/
│   ├── auth.go
│── controllers/
│   ├── authController.go
│── middleware/
│   ├── jwtMiddleware.go
│── .env
│── go.mod
│── go.sum
``` -->

## License
This project is licensed under the [MIT License](LICENSE).

