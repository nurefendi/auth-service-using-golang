# Auth Service with Golang, Fiber, GORM, and MariaDB

## Description
This authentication service is built using Golang with the Fiber framework, utilizing GORM as the ORM, and MariaDB as the primary database.

## Features
- User registration
- User login with JWT
- Authentication middleware
- Refresh token
- User management

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

