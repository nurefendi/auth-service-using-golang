package routers

import (
	"auth-service/controllers"
	"auth-service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func HandleRouter(app *fiber.App) {
	log.Info("Initialize routers")
	
	api := app.Group("/v1", func(c *fiber.Ctx) error { 
        c.Set("Version", "v1")
        return c.Next()
    })
	api.Post("/auth/register", middleware.SetMiddlewareJSON(), controllers.AuthRegister)
	api.Post("/auth/login", middleware.SetMiddlewareJSON())
	
	// need whitelist this path in midleware
	api.Get("/auth/logout", middleware.SetMiddlewareAUTH())
	api.Get("/auth/me", middleware.SetMiddlewareAUTH())
	api.Get("/auth/acl", middleware.SetMiddlewareAUTH())
	api.Get("/portal", middleware.SetMiddlewareAUTH())

	api.Post("/portal", middleware.SetMiddlewareAUTH(), controllers.SavePortal)
	api.Get("/portal/:id", middleware.SetMiddlewareAUTH(), controllers.SavePortal)
	api.Put("/portal", middleware.SetMiddlewareAUTH(), controllers.SavePortal)
	api.Delete("/portal/:id", middleware.SetMiddlewareAUTH(), controllers.SavePortal)
	
	api.Get("/acl/function", middleware.SetMiddlewareAUTH())
	api.Get("/acl/function/:id", middleware.SetMiddlewareAUTH())
	api.Post("/acl/function", middleware.SetMiddlewareAUTH())
	api.Put("/acl/function", middleware.SetMiddlewareAUTH())
	api.Delete("/acl/function/:id", middleware.SetMiddlewareAUTH())

}