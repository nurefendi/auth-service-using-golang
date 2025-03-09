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
	api.Post("/auth/login", middleware.SetMiddlewareJSON(), controllers.AuthLogin)
	
	// need whitelist this path in midleware
	api.Get("/auth/logout", middleware.SetMiddlewareAuthNoAcl(), controllers.AuthLogout)
	api.Get("/auth/refresh-token", middleware.SetMiddlewareAuthNoAcl(), controllers.AuthRefreshTokens)
	api.Get("/auth/me", middleware.SetMiddlewareAuthNoAcl(), controllers.AuthMe)
	api.Get("/auth/chek-access", middleware.SetMiddlewareAuthNoAcl(), controllers.CheckAccess)
	api.Get("/portal", middleware.SetMiddlewareAuthNoAcl())
	api.Post("/portal", middleware.SetMiddlewareAUTH(), controllers.SavePortal)
	api.Get("/portal/:id", middleware.SetMiddlewareAUTH(), controllers.GetPortalById)
	api.Put("/portal", middleware.SetMiddlewareAUTH(), controllers.UpdatePortal)
	api.Delete("/portal/:id", middleware.SetMiddlewareAUTH(), controllers.DeletePortalById)
	api.Get("/acl/function", middleware.SetMiddlewareAuthNoAcl())
	api.Get("/acl/function/:id", middleware.SetMiddlewareAUTH())
	api.Post("/acl/function", middleware.SetMiddlewareAUTH())
	api.Put("/acl/function", middleware.SetMiddlewareAUTH())
	api.Delete("/acl/function/:id", middleware.SetMiddlewareAUTH())

}