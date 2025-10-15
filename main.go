package main

import (
	"github.com/nurefendi/auth-service-using-golang/config"
	"github.com/nurefendi/auth-service-using-golang/config/database"
	"github.com/nurefendi/auth-service-using-golang/routers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Print("Routed to port "+port, " ", fiber.IsChild())

	// start server in goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("[FATAL] Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
}
