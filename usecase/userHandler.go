package usecase

import (
	"auth-service/dto"
	userRepository "auth-service/repository/database/authuser"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func UserUSeCase() User {
	return &userUseCase{}
}

// Delete implements User.
func (u *userUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	_, errf := userRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}
	return userRepository.Delete(c, id)
}

// FindAll implements User.
func (u *userUseCase) FindAll(c *fiber.Ctx, r dto.UserPagination) ([]dto.UserDto, int64, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, total, errf := userRepository.FindAll(c, r)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil,0, errf
	}
	var result []dto.UserDto
	for _, v := range data {
		result = append(result, dto.UserDto{
			ID: &v.ID,
			Gender: v.Gender,
			FullName: v.FullName,
			Email: v.Email,
			Username: v.Username,
			Telephone: v.Telephone,
			Picture: v.Picture,
		})
	}
	return result, total, nil
}

// FindById implements User.
func (u *userUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.UserDto, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, errf := userRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil, errf
	}
	var groupIds []uuid.UUID
	for _, v := range data.Goups {
		groupIds = append(groupIds, v.ID)
	}
	return &dto.UserDto{
		ID: &data.ID,
		Gender: data.Gender,
		FullName: data.FullName,
		Email: data.Email,
		Username: data.Username,
		Telephone: data.Telephone,
		Picture: data.Picture,
		GroupIDs: groupIds,

	}, &fiber.Error{}
}

// Save implements User.
func (u *userUseCase) Save(c *fiber.Ctx, data *dto.UserDto) *fiber.Error {
	panic("unimplemented")
}

// Update implements User.
func (u *userUseCase) Update(c *fiber.Ctx, data *dto.UserDto) *fiber.Error {
	panic("unimplemented")
}