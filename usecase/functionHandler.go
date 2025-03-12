package usecase

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FunctionUSeCase() Function {
	return &functionUseCase{}
}

// Delete implements Function.
func (f *functionUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	panic("unimplemented")
}

// FindAll implements Function.
func (f *functionUseCase) FindAll(c *fiber.Ctx, r dto.FunctionPagination) ([]dto.PortalUserDto, int64, *fiber.Error) {
	panic("unimplemented")
}

// FindById implements Function.
func (f *functionUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.PortalDto, *fiber.Error) {
	panic("unimplemented")
}

// Save implements Function.
func (f *functionUseCase) Save(c *fiber.Ctx, data *dto.FunctionDto) *fiber.Error {
	panic("unimplemented")
}

// Update implements Function.
func (f *functionUseCase) Update(c *fiber.Ctx, data *dto.FunctionDto) *fiber.Error {
	panic("unimplemented")
}
