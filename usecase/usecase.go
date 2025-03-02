package usecase

import "github.com/gofiber/fiber/v2"

type (
	authUseCase struct{}

	Auth interface {
		Register(c *fiber.Ctx) *fiber.Ctx
		Login(c *fiber.Ctx) *fiber.Ctx
	}
)
