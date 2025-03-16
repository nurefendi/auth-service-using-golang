package usecase

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GroupUSeCase() Group {
	return &groupUseCase{}
}

// Delete implements Group.
func (g *groupUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	panic("unimplemented")
}

// FindAll implements Group.
func (g *groupUseCase) FindAll(c *fiber.Ctx, r dto.GroupPagination) ([]dto.GroupDto, int64, *fiber.Error) {
	panic("unimplemented")
}

// FindById implements Group.
func (g *groupUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.GroupDto, *fiber.Error) {
	panic("unimplemented")
}

// Save implements Group.
func (g *groupUseCase) Save(c *fiber.Ctx, data *dto.GroupDto) *fiber.Error {
	panic("unimplemented")
}

// Update implements Group.
func (g *groupUseCase) Update(c *fiber.Ctx, data *dto.GroupDto) *fiber.Error {
	panic("unimplemented")
}
