package usecase

import (
	"auth-service/dto"
	"auth-service/repository/dao"
	groupRepository "auth-service/repository/database/authgroup"
	"auth-service/tools/helper"
	"auth-service/tools/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func GroupUSeCase() Group {
	return &groupUseCase{}
}

// Delete implements Group.
func (g *groupUseCase) Delete(c *fiber.Ctx, id uuid.UUID) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	_, errf := groupRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}
	return groupRepository.Delete(c, id)
}

// FindAll implements Group.
func (g *groupUseCase) FindAll(c *fiber.Ctx, r dto.GroupPagination) ([]dto.GroupDto, int64, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, total, errf := groupRepository.FindAll(c, r)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil,0, errf
	}
	var result []dto.GroupDto
	for _, v := range data {
		result = append(result, dto.GroupDto{
			ID: &v.ID,
			Name: v.Name,
			Description: v.Description,
		})
	}
	return result, total, nil
}

// FindById implements Group.
func (g *groupUseCase) FindById(c *fiber.Ctx, id uuid.UUID) (*dto.GroupDto, *fiber.Error) {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	data, errf := groupRepository.FindById(c, id)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return nil, errf
	}
	var result dto.GroupDto
	err := helper.Map(data, &result)
	if err != nil {
		log.Error(currentAccess.RequestID, err.Error())
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Erorr bind data")
	}
	return &result, &fiber.Error{}
}

// Save implements Group.
func (g *groupUseCase) Save(c *fiber.Ctx, data *dto.GroupDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	saveError := groupRepository.Save(c, &dao.AuthGroup{
		Name: data.Name,
		Description: data.Description,
		AuditorDAO: dao.AuditorDAO{
			CreatedBy: currentAccess.UserAccess.UserName,
		},
	})
	if saveError != nil {
		log.Error(currentAccess.RequestID, saveError.Error())
		return fiber.NewError(fiber.StatusInternalServerError, " Can't save data to db")
	}
	return nil
}

// Update implements Group.
func (g *groupUseCase) Update(c *fiber.Ctx, data *dto.GroupDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if data.ID == nil {
		log.Error(currentAccess.RequestID, " ID must not nul")
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}
	_, errf := groupRepository.FindById(c, *data.ID)
	if errf != nil {
		log.Error(currentAccess.RequestID, errf.Message)
		return errf
	}

	saveError := groupRepository.Save(c, &dao.AuthGroup{
		Name: data.Name,
		Description: data.Description,
		AuditorDAO: dao.AuditorDAO{
			ModifiedBy: &currentAccess.UserAccess.Email,
		},
	})
	if saveError != nil {
		log.Error(currentAccess.RequestID, "Failed to save function:", saveError.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save data")
	}

	return nil
}
