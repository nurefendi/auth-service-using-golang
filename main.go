package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"

	"github.com/nurefendi/auth-service-using-golang/config"
	"github.com/nurefendi/auth-service-using-golang/config/database"
	grpcServer "github.com/nurefendi/auth-service-using-golang/grpc"
	"github.com/nurefendi/auth-service-using-golang/grpc/pb"
	"github.com/nurefendi/auth-service-using-golang/routers"
)

func main() {
	config.Init()
	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	if !app.Config().Prefork {
		database.InitGlobalDB()
	}

	// Ensure DB closed on exit
	defer database.CloseDBConnection()

	// Middleware DB (Hanya berguna untuk Prefork)
	app.Use(database.DBMiddleware())
	app.Use(logger.New())

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8001",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,OPTIONS,PUT,DELETE,PATCH",
		AllowCredentials: true,
	}))
	routers.HandleRouter(app)

	// Setup gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9001"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	authServer := grpcServer.NewAuthServer()
	pb.RegisterAuthServiceServer(s, authServer)

	log.Printf("gRPC server listening on port %s", grpcPort)

	// Start gRPC server in goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[FATAL] Failed to start gRPC server: %v", err)
		}
	}()

	// Setup REST API server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Printf("REST API server listening on port %s", port)

	// start REST API server in goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("[FATAL] Failed to start REST API server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Stop gRPC server
	s.GracefulStop()
	log.Println("gRPC server stopped")

	// Stop REST API server
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down REST API server: %v", err)
	}
	log.Println("REST API server stopped")
}
