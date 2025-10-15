package usecase

import (
	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/repository/dao"
	userRepository "github.com/nurefendi/auth-service-using-golang/repository/database/authuser"
	"github.com/nurefendi/auth-service-using-golang/tools/helper"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"
	"time"

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
		return nil, 0, errf
	}
	var result []dto.UserDto
	for _, v := range data {
		result = append(result, dto.UserDto{
			ID:        &v.ID,
			Gender:    v.Gender,
			FullName:  v.FullName,
			Email:     v.Email,
			Username:  v.Username,
			Telephone: v.Telephone,
			Picture:   v.Picture,
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
	for _, v := range data.Groups {
		groupIds = append(groupIds, v.ID)
	}
	return &dto.UserDto{
		ID:        &data.ID,
		Gender:    data.Gender,
		FullName:  data.FullName,
		Email:     data.Email,
		Username:  data.Username,
		Telephone: data.Telephone,
		Picture:   data.Picture,
		GroupIDs:  groupIds,
	}, &fiber.Error{}
}

// Save implements User.
func (u *userUseCase) Save(c *fiber.Ctx, data *dto.UserDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	log.Info(currentAccess.RequestID, " add user ", data.Email)

	password, err := helper.HashBcript(data.Password)
	if err != nil {
		log.Error(currentAccess.RequestID, " Invalid hash password")
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Invalid hash password",
		}
	}
	existUser, fibererr := userRepository.FindByEmail(c, &data.Email)
	if fibererr != nil || existUser.ID != uuid.Nil {
		return &fiber.Error{
			Code:    fibererr.Code,
			Message: fibererr.Error(),
		}
	}

	var groups []dao.AuthUserGroup
	for _, v := range data.GroupIDs {
		groups = append(groups, dao.AuthUserGroup{
			GroupID: v,
			AuditorDAO: dao.AuditorDAO{
				CreatedBy: currentAccess.UserAccess.Email,
				ID:        uuid.New(),
			},
		})
	}
	dataUser := dao.AuthUser{
		FullName:   data.FullName,
		Email:      data.Email,
		Password:   password,
		Gender:     data.Gender,
		HasDeleted: false,
		Username:   data.Username,
		Picture:    data.Picture,
		Groups:     groups,
		Telephone:  data.Telephone,
		AuditorDAO: dao.AuditorDAO{
			ID:        uuid.New(),
			CreatedBy: currentAccess.UserAccess.Email,
		},
	}
	c.Locals(locals.Entity, dataUser)
	fibererr = userRepository.Save(c)
	if fibererr != nil {
		return fibererr
	}
	return nil
}

// Update implements User.
func (u *userUseCase) Update(c *fiber.Ctx, r *dto.UserDto) *fiber.Error {
	currentAccess := locals.GetLocals[dto.UserLocals](c, locals.UserLocalKey)
	if r.ID == &uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}
	user, errr := userRepository.FindById(c, *r.ID)
	if errr != nil {
		return errr
	}
	curtime := time.Now()
	user.FullName = r.FullName
	user.ModifiedAt = &curtime
	user.ModifiedBy = &currentAccess.UserAccess.Email
	user.Email = r.Email
	user.Gender = r.Gender
	user.Picture = r.Picture
	user.Telephone = r.Telephone
	c.Locals(locals.Entity, user)
	return userRepository.Save(c)

}
