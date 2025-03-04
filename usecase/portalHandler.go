package usecase

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PortalUseCase() Portal {
	return &portalUseCase{}
}

// Delete implements Portal.
func (p *portalUseCase) Delete(c *fiber.Ctx, id *uuid.UUID) *fiber.Error {
	panic("unimplemented")
}

// Save implements Portal.
func (p *portalUseCase) Save(c *fiber.Ctx, data *dto.PortalSaveRequest) *fiber.Error {
	panic("unimplemented")
}

// Update implements Portal.
func (p *portalUseCase) Update(c *fiber.Ctx, data *dto.PortalSaveRequest) *fiber.Error {
	panic("unimplemented")
}
