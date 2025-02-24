package main

import (
	"auth-service/config"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:     true,
	})

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods:"GET,POST,OPTIONS,PUT,DELETE,PATCH",
	}))

	config.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Print("Routed to port " + port, " ", fiber.IsChild())
	log.Fatal(app.Listen(":" + port))
}