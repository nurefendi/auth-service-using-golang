package usecase

import "github.com/gofiber/fiber/v2"

type (
	authUseCase struct{}

	Auth interface {
		Register(c *fiber.Ctx) *fiber.Error
		CheckEmailExist(c *fiber.Ctx, email *string) *fiber.Error
		Login(c *fiber.Ctx) *fiber.Error
	}
)
