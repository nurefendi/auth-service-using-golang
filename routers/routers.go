package routers

import (
	"auth-service/controllers"
	"auth-service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HandleRouter(app *fiber.App) {
	log.Info("Initialize routers")
	api := app.Group("/v1")
	api.Post("/register", middleware.SetMiddlewareJSON(), controllers.AuthRegister)
	// api.Post("/login", middleware.SetMiddlewareJSON())
	// api.Get("/logout")
	// api.Get("/me")
	// api.Post("/acl")

}