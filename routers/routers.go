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
	api.Get("/auth/logout", middleware.SetMiddlewareAUTH(false), controllers.AuthLogout)
	api.Get("/auth/refresh-token", middleware.SetMiddlewareAUTH(false), controllers.AuthRefreshTokens)
	api.Get("/auth/me", middleware.SetMiddlewareAUTH(false), controllers.AuthMe)
	api.Get("/auth/chek-access", middleware.SetMiddlewareAUTH(false), controllers.CheckAccess)

	api.Get("/portal", middleware.SetMiddlewareAUTH(true), controllers.GetPortal)
	api.Post("/portal", middleware.SetMiddlewareAUTH(true), controllers.SavePortal)
	api.Get("/portal/:id", middleware.SetMiddlewareAUTH(true), controllers.GetPortalById)
	api.Put("/portal", middleware.SetMiddlewareAUTH(true), controllers.UpdatePortal)
	api.Delete("/portal/:id", middleware.SetMiddlewareAUTH(true), controllers.DeletePortalById)

	api.Get("/function", middleware.SetMiddlewareAUTH(false), controllers.GetFunction)
	api.Get("/function/:id", middleware.SetMiddlewareAUTH(true), controllers.GetFunctionById)
	api.Post("/function", middleware.SetMiddlewareAUTH(true), controllers.SaveFunction)
	api.Put("/function", middleware.SetMiddlewareAUTH(true), controllers.UpdateFunction)
	api.Delete("/function/:id", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)

	api.Get("/user/:id", middleware.SetMiddlewareAUTH(true), controllers.GetUserById)
	api.Get("/user", middleware.SetMiddlewareAUTH(true), controllers.GetUser)
	api.Post("/user", middleware.SetMiddlewareAUTH(true), controllers.SaveUser)
	api.Put("/user", middleware.SetMiddlewareAUTH(true), controllers.UpdateUser)
	api.Delete("/user/:id", middleware.SetMiddlewareAUTH(true), controllers.DeleteUserById)

	api.Get("/group/:id", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
	api.Get("/group", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
	api.Post("/group", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
	api.Put("/group", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
	api.Delete("/group/:id", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
	
	api.Get("/acl", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
	api.Post("/acl", middleware.SetMiddlewareAUTH(true), controllers.DeleteFunctionById)
}