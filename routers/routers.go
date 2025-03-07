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
	api.Get("/auth/logout", middleware.SetMiddlewareAuthNoAcl())
	api.Get("/auth/me", middleware.SetMiddlewareAuthNoAcl())
	api.Get("/auth/acl", middleware.SetMiddlewareAuthNoAcl())
	api.Get("/portal", middleware.SetMiddlewareAuthNoAcl())

	api.Post("/portal", middleware.SetMiddlewareAUTH(), controllers.SavePortal)
	api.Get("/portal/:id", middleware.SetMiddlewareAUTH(), controllers.GetPortalById)
	api.Put("/portal", middleware.SetMiddlewareAUTH(), controllers.UpdatePortal)
	api.Delete("/portal/:id", middleware.SetMiddlewareAUTH(), controllers.DeletePortalById)
	
	api.Get("/acl/function", middleware.SetMiddlewareAUTH())
	api.Get("/acl/function/:id", middleware.SetMiddlewareAUTH())
	api.Post("/acl/function", middleware.SetMiddlewareAUTH())
	api.Put("/acl/function", middleware.SetMiddlewareAUTH())
	api.Delete("/acl/function/:id", middleware.SetMiddlewareAUTH())

}