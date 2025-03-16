package usecase

import (
	"auth-service/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UserUSeCase() User {
	return &userUseCase{}
}

// Delete implements User.
func (u *userUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	panic("unimplemented")
}

// FindAll implements User.
func (u *userUseCase) FindAll(c *fiber.Ctx, r dto.UserPagination) ([]dto.UserDto, int64, *fiber.Error) {
	panic("unimplemented")
}

// FindById implements User.
func (u *userUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.UserDto, *fiber.Error) {
	panic("unimplemented")
}

// Save implements User.
func (u *userUseCase) Save(c *fiber.Ctx, data *dto.UserDto) *fiber.Error {
	panic("unimplemented")
}

// Update implements User.
func (u *userUseCase) Update(c *fiber.Ctx, data *dto.UserDto) *fiber.Error {
	panic("unimplemented")
}