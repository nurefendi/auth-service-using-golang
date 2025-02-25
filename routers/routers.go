package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HandleRouter(app *fiber.App) {
	log.Info("Initialize routers")

	app.Post("/login")
	app.Get("/logout")
	app.Get("/me")
	app.Post("/acl")

}