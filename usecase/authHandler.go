package usecase

import "github.com/gofiber/fiber/v2"

func AuthUSeCase() Auth {
	return &authUseCase{}
}

// Login implements Auth.
func (a *authUseCase) Login(c *fiber.Ctx) *fiber.Ctx {
	panic("unimplemented")
}

// Register implements Auth.
func (a *authUseCase) Register(c *fiber.Ctx) *fiber.Ctx {
	panic("unimplemented")
}