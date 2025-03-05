package usecase

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	authUseCase   struct{}
	portalUseCase struct{}

	Auth interface {
		Register(c *fiber.Ctx) *fiber.Error
		CheckEmailExist(c *fiber.Ctx, email *string) *fiber.Error
		CheckUserNameExist(c *fiber.Ctx, userName *string) *fiber.Error
		Login(c *fiber.Ctx) *fiber.Error
	}
	Portal interface {
		Save(c *fiber.Ctx, data *dto.PortalSaveRequest) *fiber.Error
		Update(c *fiber.Ctx, data *dto.PortalSaveRequest) *fiber.Error
		Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error
	}
)
