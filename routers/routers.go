package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HandleRouter(app *fiber.App) {
	log.Info("Initialize routers")
	api := app.Group("/v1")
	api.Post("/login")
	api.Get("/logout")
	api.Get("/me")
	api.Post("/acl")

}